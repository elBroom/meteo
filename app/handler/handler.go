package handler

import (
	"net/http"

	"log"

	"strconv"

	"strings"

	"regexp"

	"time"

	"github.com/elBroom/meteo/app/config"
	"github.com/elBroom/meteo/app/db"
	"github.com/elBroom/meteo/app/model"
	"github.com/elBroom/meteo/app/schema"
	"github.com/elBroom/meteo/app/ws"
	"github.com/fasthttp-contrib/websocket"
	"github.com/getsentry/raven-go"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func init() {
	raven.SetDSN(config.GetApp().SentryDsn)
}

func CreteValueEndpoint(ctx *fasthttp.RequestCtx) {
	log.Print(ctx.RemoteIP(), " ", string(ctx.Path()), " ", string(ctx.UserAgent()))

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

	value_str := strings.Replace(string(data), "value=", "", 1)
	value, _ := strconv.ParseFloat(value_str, 64)

	db.Sql_connect().Create(&model.Indication{CreateDate: time.Now(), Value: value, Pin: disignation.Pin})

	var indication schema.Indication
	indication.Value = value
	indication.CreateDate = time.Now().Unix() * 1000
	indication.Pin = disignation.Pin
	hub := config.GetApp().Hub
	if hub != nil {
		hub.SendMessage(&indication)
	}

	ctx.SetStatusCode(http.StatusOK)
	writeStr(ctx, "OK")
}

func GetValuesEndpoint(ctx *fasthttp.RequestCtx) {
	log.Print(ctx.RemoteIP(), " ", string(ctx.Path()), " ", string(ctx.UserAgent()))

	if ctx.UserValue("pins") == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		writeStr(ctx, "Invalid params")
		return
	}
	pins := strings.Split(ctx.UserValue("pins").(string), ",")

	var indications []model.Indication
	query := db.Sql_connect().Where("Pin IN (?)", pins)

	params := ctx.QueryArgs()
	r, _ := regexp.Compile("^20(1|2)[0-9]-[0,1]?[0-9]-[0-3]?[0-9]$") // Y-m-d
	if params.Has("start_date") {
		start_date := string(params.Peek("start_date"))
		if r.MatchString(start_date) {
			query = query.Where("create_date::date >= ?::date", start_date)
		}
	}
	if params.Has("end_date") {
		end_date := string(params.Peek("end_date"))
		if r.MatchString(end_date) {
			query = query.Where("create_date::date <= ?::date", end_date)
		}
	}
	query.Order("pin, create_date").Find(&indications)

	pin_data := make(map[string][][]interface{})
	for _, item := range indications {
		elem := []interface{}{item.CreateDate.Unix() * 1000, item.Value}
		pin_data[item.Pin] = append(pin_data[item.Pin], elem)
	}

	var disignations []model.Disignation
	db.Sql_connect().Where("Pin IN (?)", pins).Order("pin").Find(&disignations)

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

var upgrader = websocket.New(func(conn *websocket.Conn) {
	ws.ServeWs(config.GetApp().Hub, conn)
})

func WSEndpoint(ctx *fasthttp.RequestCtx) {
	if err := upgrader.Upgrade(ctx); err != nil {
		log.Println(err)
		return
	}
}

func writeStr(ctx *fasthttp.RequestCtx, s string) {
	b := []byte(s)
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetBody(b)
}
