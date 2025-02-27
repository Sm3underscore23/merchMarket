package handler

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/Sm3underscore23/merchStore/pkg/service"
	mock_service "github.com/Sm3underscore23/merchStore/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestHandler_singUpIn(t *testing.T) {
	type mockBehavior func(
		s *mock_service.MockAuthorization,
		authRequest models.AuthRequest,
	)

	testTable := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		inputUser            models.AuthRequest
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test_username","password":"test_password"}`,
			inputUser: models.AuthRequest{
				Username: "test_username",
				Password: "test_password",
			},
			mockBehavior: func(
				s *mock_service.MockAuthorization,
				authRequest models.AuthRequest,
			) {
				s.EXPECT().AuthUser(authRequest.Username, authRequest.Password).Return("test_token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"test_token"}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"username":"test_username"}`,
			mockBehavior: func(
				s *mock_service.MockAuthorization,
				authRequest models.AuthRequest,
			) {
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"errors":"%s"}`,
				customerrors.ErrInvalidInputBody,
			),
		},
		{
			name:      "Service failure",
			inputBody: `{"username":"test_username","password":"test_password"}`,
			inputUser: models.AuthRequest{
				Username: "test_username",
				Password: "test_password",
			},
			mockBehavior: func(
				s *mock_service.MockAuthorization,
				authRequest models.AuthRequest,
			) {
				s.EXPECT().AuthUser(
					authRequest.Username,
					authRequest.Password,
				).Return(
					"",
					customerrors.ErrDatabase,
				)
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(
				`{"errors":"%s"}`,
				customerrors.ErrDatabase,
			),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)

			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()

			r.POST("/api/auth", handler.singUpIn)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/api/auth", bytes.NewBufferString(testCase.inputBody))

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
