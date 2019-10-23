package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	logger "github.com/jeanphorn/log4go"
	uuid "github.com/satori/go.uuid"

	"github.com/valyala/fasthttp"
)

type videoTmp struct {
	FtLastRead     time.Time
	FtCreate       time.Time
	FkID           string
	FkParent       string
	FlFinishEncode bool
	FnCurrentSize  int64
	FvPath         string
}

var videosTmp map[string]*videoTmp = make(map[string]*videoTmp)
var videosTmpUUID map[string]*videoTmp = make(map[string]*videoTmp)

func requestHandler(tmpDir string, config Config, params map[string]string, ctx *fasthttp.RequestCtx) {
	query, ok := params["query"]
	if !ok {
		ctx.Error("Required params query", 500)
		return
	}
	if query == "fileTree" {
		d, _ := videoDirRead(config.VideoDir, config.AcceptsFile)
		ctx.Response.Header.SetContentType("application/json")
		ctx.SetBodyString(string(d))
		return
	}

	if query == "getVideoUuid" {
		jsonText, ok := params["json"]
		if !ok {
			ctx.NotFound()
			return
		}
		var data dirPath
		if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
			ctx.NotFound()
			return
		}
		video, ok := videosTmp[data.FkID]
		if !ok {
			uuid, errUUID := uuid.NewV4()
			if errUUID != nil {
				logger.Error(errUUID)
				ctx.Error("Failed", 500)
				return
			}
			uuidByte, errMarchal := uuid.MarshalText()
			if errMarchal != nil {
				logger.Error(errMarchal)
				ctx.Error("Failed", 500)
				return
			}
			uuidText := string(uuidByte)
			newVideo := videoTmp{FkID: uuidText, FtCreate: time.Now(), FkParent: data.FkID, FlFinishEncode: false}
			videosTmp[data.FkID] = &newVideo
			videosTmpUUID[uuidText] = &newVideo
			video = &newVideo
			go encodeVideoFF(tmpDir, data, &newVideo, uuidText)
			for {
				if len(video.FvPath) > 0 {
					if _, err := os.Stat(video.FvPath); err == nil {
						break
					}
					if _, ok := videosTmpUUID[uuidText]; !ok {
						ctx.NotFound()
						return
					}
				}
				time.Sleep(100)
			}
		}
		ctx.Response.Header.SetContentType("application/json")
		fmt.Fprintf(ctx, "{\"fv_uuid\":\"%s\"}", video.FkID)
		return
	}

	if query == "getVideo" {
		uuidText, ok := params["uuid"]
		if !ok {
			ctx.NotFound()
			return
		}
		video, ok := videosTmpUUID[uuidText]
		if !ok {
			ctx.NotFound()
			return
		}
		readFileVideo(video, ctx)
		return
	}

	ctx.NotFound()
}
