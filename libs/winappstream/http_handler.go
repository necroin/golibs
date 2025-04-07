//go:build windows

package winappstream

import (
	"net/http"
)

type HttpImageCaptureHandler struct {
	app *App
}

func NewHttpImageCaptureHandler(app *App) HttpImageCaptureHandler {
	return HttpImageCaptureHandler{
		app: app,
	}
}

func (handler HttpImageCaptureHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Cache-Control", "max-age=2")

	promise, ok := <-handler.app.encodedData
	if !ok {
		responseWriter.Write([]byte("channel closed"))
		return
	}

	data, _ := promise.Get()
	responseWriter.Write(data)
}
