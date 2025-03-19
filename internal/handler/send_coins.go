package handler

import (
	"net/http"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) sendCoins(c *gin.Context) {
	var input models.SendCoinRequest

	if err := c.BindJSON(&input); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, customerrors.ErrInvalidInputBody.Error())
		return
	}

	senderUserID, err := getIdFromCtx(c)
	if err != nil {
		return
	}

	err = h.service.SendCoins.SendCoins(input.Receiver, senderUserID, input.Coins)
	if err != nil {
		models.NewErrorResponse(c, customerrors.ErrWithStatus[err], err.Error())
		return
	}
}
