package main

import (
	"io"
	"io/ioutil"
	"os"
	filepath "path"
	"strings"

	"github.com/buaazp/fasthttprouter"
	logger "github.com/jeanphorn/log4go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

func startHTTP(config Config) {
	tmp := config.VideoTmp
	if len(config.VideoTmp) == 0 {
		tmpDir, errCreateDir := ioutil.TempDir("", "stream_tmp_video")
		if errCreateDir != nil {
			panic(errCreateDir)
		}
		tmp = tmpDir
	}
	if _, err := os.Stat(tmp); os.IsNotExist(err) {
		if err = os.MkdirAll(tmp, os.ModePerm); err != nil {
			panic(err)
		}
	}
	logger.Info("Read video tmp")
	readTmpDir(tmp)
	// Setup FS handler
	spaHandler := func(ctx *fasthttp.RequestCtx) {
		f, err := os.Open(filepath.Join(config.SiteDir, "index.html"))
		defer f.Close()
		if err != nil {
			logger.Error(err)
			ctx.NotFound()
			return
		}
		ctx.Response.Header.SetContentType("text/html")
		io.Copy(ctx, f)
	}
	fs := &fasthttp.FS{
		Root:               config.SiteDir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		Compress:           true,
		AcceptByteRange:    false,
		PathNotFound:       spaHandler,
	}
	fsHandler := fs.NewRequestHandler()
	apiHandler := func(ctx *fasthttp.RequestCtx) {
		params := make(map[string]string)
		fn := func(key []byte, val interface{}) {
			value, _ := val.(string)
			params[string(key)] = value
		}
		fnArgs := func(key []byte, val []byte) {
			params[string(key)] = string(val)
		}
		ctx.VisitUserValues(fn)
		if arg := ctx.QueryArgs(); arg != nil {
			arg.VisitAll(fnArgs)
		}
		if arg := ctx.PostArgs(); arg != nil {
			arg.VisitAll(fnArgs)
		}
		ctByte := ctx.Request.Header.ContentType()
		ct := strings.ToLower(string(ctByte))
		if strings.HasPrefix(ct, "application/json") {
			params["json"] = string(ctx.PostBody())
		}
		if strings.HasPrefix(ct, "application/xml") {
			params["xml"] = string(ctx.PostBody())
		}
		logger.Debug("RequestUri: %s\n", ctx.RequestURI())
		ctx.Request.Header.VisitAll(func(key, value []byte) {
			logger.Debug("Header: %s=%s\n", key, value)
		})
		logger.Debug("Params: %s\n", params)
		requestHandler(tmp, config, params, ctx)
	}

	router := fasthttprouter.New()
	router.NotFound = fsHandler
	router.GET("/stats", expvarhandler.ExpvarHandler)
	router.GET("/api/:query", apiHandler)
	router.POST("/api/:query", apiHandler)
	logger.Info("Http start %s", config.Listen)
	if err := fasthttp.ListenAndServe(config.Listen, router.Handler); err != nil {
		logger.Error("error in ListenAndServe: %s", err)
		panic(err)
	}
}
