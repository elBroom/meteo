package handler

import (
	"net/http"

	"log"

	"strconv"

	"strings"

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

	var disignation model.Disignation
	ok := !db.Sql_connect().Where("Pin = ?", ctx.UserValue("pin")).First(&disignation).RecordNotFound()
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
	log.Println("GetValuesEndpoint")

	if ctx.UserValue("pins") == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		writeStr(ctx, "Invalid params")
		return
	}
	pins := strings.Split(ctx.UserValue("pins").(string), ",")

	var indications []model.Indication
	db.Sql_connect().Where("Pin IN (?)", pins).Find(&indications).Order("Pin,CreatedAt")

	pin_data := make(map[string][][]interface{})
	for _, item := range indications {
		elem := []interface{}{item.CreatedAt.Unix() * 1000, item.Value}
		pin_data[item.Pin] = append(pin_data[item.Pin], elem)
	}

	var disignations []model.Disignation
	db.Sql_connect().Where("Pin IN (?)", pins).Find(&disignations).Order("Pin")
	var data schema.DisignationList
	for _, item := range disignations {
		_, ok := pin_data[item.Pin]
		if ok {
			var disignation schema.Disignation
			disignation.Pin = item.Pin
			disignation.Color = item.Color
			disignation.Name = item.Name
			disignation.Unit = item.Unit
			disignation.Data = pin_data[item.Pin]

			data = append(data, disignation)
		}
	}

	b, _ := easyjson.Marshal(data)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetBody(b)
}

func writeStr(ctx *fasthttp.RequestCtx, s string) {
	b := []byte(s)
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetBody(b)
}
