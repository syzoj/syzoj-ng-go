package http

import (
	"encoding/json"
	"fmt"
	"mime"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var log = logrus.StandardLogger()

func SendJSON(ctx *fasthttp.RequestCtx, data interface{}) {
	val, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("SendJSON: Failed to marshal data")
		return
	}
	ctx.SetBody(val)
}

func SendError(ctx *fasthttp.RequestCtx, err string) {
	SendJSON(ctx, struct {
		Error string `json:"error"`
	}{Error: err})
}

func SendInternalError(ctx *fasthttp.RequestCtx, err error) {
	log.WithError(err).Error("Internal server error")
	ctx.SetStatusCode(500)
	SendError(ctx, "Internal server error")
}

func SendConflict(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(409)
	SendError(ctx, err.Error())
}

func ReadBodyJSON(ctx *fasthttp.RequestCtx, val interface{}) error {
	mediatype, _, err := mime.ParseMediaType(string(ctx.Request.Header.ContentType()))
	if err != nil {
		return err
	}
	if mediatype != "application/json" {
		return fmt.Errorf("Invalid media type: %s", mediatype)
	}
	return json.Unmarshal(ctx.Request.Body(), val)
}

func BadRequest(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(400)
	SendError(ctx, err.Error())
}

func NotFound(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(404)
	SendError(ctx, err.Error())
}
