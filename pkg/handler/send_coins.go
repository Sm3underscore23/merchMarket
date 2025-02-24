package handler

import (
	"fmt"
	"net/http"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) sendCoins(c *gin.Context) {
	var input models.SendCoinRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(input.Receiver, input.Coins)

	senderUserID, err := getIdFromCtx(c)
	if err != nil {
		return
	}

	err = h.service.SendCoins.SendCoins(input.Receiver, senderUserID, input.Coins)
	if err != nil {
		statusCode, message := customerrors.ClassifyError(err)
		newErrorResponse(c, statusCode, message)
		return
	}
}
