package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/pkg/service"
	mock_service "github.com/Sm3underscore23/merchStore/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name               string
		headerName         string
		headerValue        string
		token              string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: "1",
		},
		{
			name:               "No header",
			headerName:         "",
			headerValue:        "Bearer token",
			token:              "token",
			mockBehavior:       func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name:               "Invalid bearer",
			headerName:         "Authorization",
			headerValue:        "Beer token",
			mockBehavior:       func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"errors":"invalid auth header"}`,
		},
		{
			name:               "Invalid token",
			headerName:         "Authorization",
			headerValue:        "Bearer ",
			mockBehavior:       func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"errors":"empty token"}`,
		},
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, customerrors.ErrParseToken)
			},
			expectedStatusCode: 401,
			expectedResponseBody: fmt.Sprintf(
				`{"errors":"%s"}`,
				customerrors.ErrParseToken,
			),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)

			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get(ctxUserId)
				ctx.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			if testCase.expectedStatusCode != w.Code {
				t.Error("stutus codes are different")
			}

			if testCase.expectedResponseBody != w.Body.String() {
				t.Error("response bodies are different")
			}

		})
	}
}
