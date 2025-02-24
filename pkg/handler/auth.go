package handler

import (
	"errors"
	"net/http"

	merchstore "github.com/Sm3underscore23/merchStore"
	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) singUpIn(c *gin.Context) {
	var input merchstore.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.GetUser(input.Username, input.Password)

	if errors.Is(err, customerrors.ErrWrongPasswod) {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.service.Authorization.GenerateToken(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
