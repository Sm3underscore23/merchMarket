package handler

import (
	"net/http"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) singUpIn(c *gin.Context) {
	var input models.AuthRequest
	if err := c.BindJSON(&input); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, customerrors.ErrInvalidInputBody.Error())
		return
	}

	token, err := h.service.Authorization.AuthUser(input.Username, input.Password)
	if err != nil {
		models.NewErrorResponse(c, customerrors.ErrWithStatus[err], err.Error())
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
	})
}
