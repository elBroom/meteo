package router

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/elBroom/meteo/app/handler"
)

func Routing() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/meteo/:token/update/:pin", handler.CreteValueEndpoint)
	router.GET("/meteo/get_data/:pins", handler.GetValuesEndpoint)
	router.GET("/meteo/ws/:channel", handler.WSEndpoint)
	return router
}
