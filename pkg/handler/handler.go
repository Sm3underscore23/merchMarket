package handler

import (
	"github.com/Sm3underscore23/merchStore/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/info", h.getInfo)
		api.POST("/sendCoin", h.sendCoins)
		api.GET("/buy/:id", h.buyItem)
		api.POST("/auth", h.singUpIn)
	}

	return router
}
