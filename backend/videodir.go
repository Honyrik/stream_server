package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	filepath "path"
	"sort"
	"strconv"
	"strings"
	"time"

	logger "github.com/jeanphorn/log4go"
	"github.com/valyala/fasthttp"
)

type dirPath struct {
	FkID        string    `json:"fk_id"`
	FvName      string    `json:"fv_name"`
	FlLeaf      bool      `json:"fl_leaf"`
	FctChildren []dirPath `json:"fct_children"`
}

func parse(dir string, accepts []string) (paths []dirPath, err error) {
	result := []dirPath{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(a, b int) bool {
		af := files[a]
		bf := files[b]
		if (af.IsDir() && bf.IsDir()) || (!af.IsDir() && !bf.IsDir()) {
			arr := []string{af.Name(), bf.Name()}
			sort.Strings(arr)
			return arr[0] == af.Name()
		}
		if af.IsDir() && !bf.IsDir() {
			return true
		}
		return false
	})
	fn := func(ipaths interface{}, ifile interface{}) (interface{}, error) {
		f, _ := ifile.(os.FileInfo)
		paths, _ := ipaths.([]dirPath)
		fealPath := filepath.Join(dir, f.Name())
		if f.IsDir() {
			children, err := parse(fealPath, accepts)
			if err == nil && len(children) > 0 {
				p := dirPath{
					FkID:        fealPath,
					FvName:      f.Name(),
					FlLeaf:      false,
					FctChildren: children,
				}
				return append(paths, p), nil
			}

		}
		name := strings.ToLower(f.Name())
		filterfn := func(iaccept interface{}) bool {
			accept, _ := iaccept.(string)
			if accept == "all" {
				return true
			}
			return strings.HasSuffix(name, accept)
		}
		if HasElemFn(accepts, filterfn) {
			p := dirPath{
				FkID:   fealPath,
				FvName: f.Name(),
				FlLeaf: true,
			}
			return append(paths, p), nil
		}
		return paths, nil
	}
	res, err := ReduceErr(files, result, fn)
	if err != nil {
		return result, err
	}
	result, _ = res.([]dirPath)

	return result, nil
}

//VideoDirRead(dir string, accept string)
func videoDirRead(dir string, accept string) ([]byte, error) {
	paths, err := parse(dir, strings.Split(strings.ToLower(accept), ","))
	if err != nil {
		return nil, err
	}
	json, errMarshal := json.Marshal(paths)
	if errMarshal != nil {
		return nil, errMarshal
	}
	return json, nil
}

func readTmpDir(tmp string) {
	files, err := ioutil.ReadDir(tmp)
	if err != nil {
		logger.Warn(err)
		return
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".meta") {
			data, err := ioutil.ReadFile(filepath.Join(tmp, file.Name()))

			if err != nil {
				logger.Warn(err)
				continue
			}
			var video videoTmp
			err = json.Unmarshal(data, &video)
			if err != nil {
				logger.Warn(err)
				continue
			}
			videosTmp[video.FkParent] = &video
			videosTmpUUID[video.FkID] = &video
		}
		if strings.HasSuffix(file.Name(), ".webm") {
			metaname := strings.Replace(file.Name(), ".webm", ".meta", -1)
			name := filepath.Join(tmp, metaname)
			if _, err := os.Stat(name); os.IsNotExist(err) {
				err = os.RemoveAll(filepath.Join(tmp, file.Name()))
				if err != nil {
					logger.Warn(err)
				}
			}
		}
	}
}

func encodeVideoFF(tmpDir string, data dirPath, video *videoTmp, uuid string) {
	f, err := os.Open(data.FkID)
	if err != nil {
		delete(videosTmp, data.FkID)
		delete(videosTmpUUID, uuid)
		logger.Error(err)
		return
	}
	state, err := f.Stat()
	if err != nil {
		delete(videosTmp, data.FkID)
		delete(videosTmpUUID, uuid)
		logger.Error(err)
		return
	}
	video.FnCurrentSize = state.Size()
	f.Close()
	video.FvPath = filepath.Join(tmpDir, uuid+".webm")
	cmd := exec.Command("ffmpeg", "-i", data.FkID, "-c:v", "libvpx-vp9", "-c:a", "libvorbis", "-preset", "ultrafast", "-cpu-used", "-5", "-deadline", "realtime", "-b:v", "2m", video.FvPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err)
		logger.Error(string(out))
		delete(videosTmp, data.FkID)
		delete(videosTmpUUID, uuid)
		return
	}
	f, err = os.Open(video.FvPath)
	if err != nil {
		logger.Error(err)
		return
	}
	state, err = f.Stat()
	if err != nil {
		logger.Error(err)
		return
	}
	video.FnCurrentSize = state.Size()
	f.Close()
	video.FlFinishEncode = true
	jsonData, errMarshal := json.Marshal(video)
	if errMarshal != nil {
		logger.Error(errMarshal)
		return
	}
	err = ioutil.WriteFile(filepath.Join(tmpDir, uuid+".meta"), jsonData, 0644)
	if err != nil {
		logger.Error(err)
		return
	}
}

func readFileVideo(video *videoTmp, ctx *fasthttp.RequestCtx) {
	f, err := os.Open(video.FvPath)
	if err != nil {
		ctx.NotFound()
		return
	}
	byteRange := ctx.Request.Header.Peek("Range")
	var videoRange bool = false
	var start int64

	if byteRange != nil && len(byteRange) > 0 {
		videoRange = true
		preRangeStr := string(byteRange)
		preRange := strings.Split(strings.Split(preRangeStr, "=")[1], "-")
		start, _ = strconv.ParseInt(preRange[0], 10, 64)
	}

	if videoRange {
		f.Seek(start, 0)
		var chunkLength int = 100000
		if video.FlFinishEncode {
			chunkLength = int(video.FnCurrentSize - start)
		}
		buf := make([]byte, chunkLength)
		var len int
		var errRead error
		for {
			len, errRead = f.Read(buf)
			if errRead != nil {
				if errRead.Error() == "EOF" {
					if start >= video.FnCurrentSize && video.FlFinishEncode {
						ctx.Response.Header.SetContentType("video/mp4")
						ctx.Response.Header.SetContentLength(0)
						ctx.Response.Header.Set("Accept-Ranges", "bytes")
						ctx.Response.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", video.FnCurrentSize, video.FnCurrentSize, video.FnCurrentSize))
						ctx.Response.Header.SetStatusCode(206)
						ctx.Write(nil)
						return
					}
					time.Sleep(1000)
					continue
				}
				logger.Error(errRead)
				ctx.Error("Error buffer", 404)
				return
			}
			break
		}
		ctx.Response.Header.SetContentType("video/mp4")
		ctx.Response.Header.SetContentLength(len)
		ctx.Response.Header.Set("Accept-Ranges", "bytes")
		ctx.Response.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, int(start)+len-1, video.FnCurrentSize))
		ctx.Response.Header.SetStatusCode(206)
		ctx.Write(buf[:len])
		return
	}
	if video.FlFinishEncode {
		ctx.Response.Header.SetContentType("video/mp4")
		ctx.Response.Header.SetContentLength(int(video.FnCurrentSize))
		io.Copy(ctx, f)
	}
}
