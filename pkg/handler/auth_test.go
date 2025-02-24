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
		name                string
		inputBody           string
		mockBehavior        mockBehavior
		inputUser           models.AuthRequest
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"testusername","password":"testpassword"}`,
			inputUser: models.AuthRequest{
				Username: "testusername",
				Password: "testpassword",
			},
			mockBehavior: func(
				s *mock_service.MockAuthorization,
				authRequest models.AuthRequest,
			) {
				s.EXPECT().AuthUser(authRequest.Username, authRequest.Password).Return("testtoken1234", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"testtoken1234"}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"username":"testusername"}`,
			mockBehavior: func(
				s *mock_service.MockAuthorization,
				authRequest models.AuthRequest,
			) {
			},
			expectedStatusCode: 400,
			expectedRequestBody: fmt.Sprintf(
				`{"errors":"%s"}`,
				customerrors.ErrInvalidInputBody.Error(),
			),
		},
		{
			name:      "Service failure",
			inputBody: `{"username":"testusername","password":"testpassword"}`,
			inputUser: models.AuthRequest{
				Username: "testusername",
				Password: "testpassword",
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
			expectedRequestBody: fmt.Sprintf(
				`{"errors":"%s"}`,
				customerrors.ErrDatabase.Error(),
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

			if testCase.expectedStatusCode != w.Code ||
				testCase.expectedRequestBody != w.Body.String() {
				t.Error("stutus codes or request bodies are different")
			}
		})
	}
}
