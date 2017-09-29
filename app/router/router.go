package router

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/elBroom/meteo/app/handler"
)

func Routing() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/meteo/:token/upload", handler.CreteValueEndpoint)
	router.GET("/meteo/get_data/:pins", handler.GetValuesEndpoint)
	return router
}
