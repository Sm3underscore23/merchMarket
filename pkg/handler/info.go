package handler

import (
	"net/http"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getInfo(c *gin.Context) {
	userId, err := getIdFromCtx(c)
	if err != nil {
		return
	}

	var userInfo models.UserInfoResponse

	userInfo, err = h.service.Info.GetUserInfo(userId)
	if err != nil {
		statusCode, message := customerrors.ClassifyError(err)
		newErrorResponse(c, statusCode, message)
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
