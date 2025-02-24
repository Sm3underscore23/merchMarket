package handler

import (
	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) buyItem(c *gin.Context) {
	productType := c.Param("slug")

	userId, err := getIdFromCtx(c)
	if err != nil {
		return
	}

	err = h.service.Buy.Buy(userId, productType)

	if err != nil {
		statusCode, message := customerrors.ClassifyError(err)
		newErrorResponse(c, statusCode, message)
		return
	}
}
