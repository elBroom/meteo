package handler

import (
	"net/http"

	"log"

	"strconv"

	"github.com/elBroom/meteo/app/config"
	"github.com/elBroom/meteo/app/db"
	"github.com/elBroom/meteo/app/model"
	"github.com/elBroom/meteo/app/schema"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func CreteValueEndpoint(ctx *fasthttp.RequestCtx) {
	log.Println("CreteValueEndpoint")

	data := ctx.PostBody()
	if len(data) == 0 || ctx.UserValue("token") == "" || ctx.UserValue("pin") == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		writeStr(ctx, "Invalid params")
		return
	}

	if ctx.UserValue("token") != config.GetApp().Token {
		ctx.SetStatusCode(http.StatusForbidden)
		writeStr(ctx, "Invalid token")
		return
	}

	pin := ctx.UserValue("pin")
	var disignation model.Disignation
	ok := !db.Sql_connect().Where("Pin = ?", pin).First(&disignation).RecordNotFound()
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		writeStr(ctx, "Invalid pin")
		return
	}

	var indication schema.Indication
	_ = easyjson.Unmarshal(data, &indication)
	db.Sql_connect().Create(&model.Indication{Value: indication.Value, Pin: disignation.Pin})

	ctx.SetStatusCode(http.StatusOK)
	writeStr(ctx, "OK")
}

func GetValuesEndpoint(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
}

func writeStr(ctx *fasthttp.RequestCtx, s string) {
	b := []byte(s)
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetBody(b)
}
