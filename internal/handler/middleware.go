package handler

import (
	"net/http"
	"strings"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	ctxUserId           = "UserId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		models.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		models.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[0] != "Bearer" {
		models.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[1] == "" {
		models.NewErrorResponse(c, http.StatusUnauthorized, "empty token")
		return
	}

	id, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(ctxUserId, id)
}

func getIdFromCtx(c *gin.Context) (int, error) {
	id, ok := c.Get(ctxUserId)
	if !ok {
		models.NewErrorResponse(c, http.StatusInternalServerError, customerrors.ErrUserNotFound.Error())
		return 0, customerrors.ErrUserNotFound
	}

	idInt, ok := id.(int)
	if !ok || idInt == 0 {
		models.NewErrorResponse(c, http.StatusInternalServerError, customerrors.ErrUserIdNotInt.Error())
		return 0, customerrors.ErrUserIdNotInt
	}

	return idInt, nil
}
