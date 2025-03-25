package handler

import (
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

func TestHandler_info(t *testing.T) {
	testUserInfoStruct := models.UserInfoResponse{
		Balance: 1000,
		Inventory: []models.InventoryItem{
			models.InventoryItem{
				ItemType: "test_item",
				Quantity: 3,
			},
		},
		CoinHistory: models.TransactionHistory{
			Received: []models.IncomingTransaction{
				models.IncomingTransaction{
					FromUser: "test_user",
					Amount:   100,
				},
			},
			Sent: []models.OutgoingTransaction{
				models.OutgoingTransaction{
					ToUser: "test_user",
					Amount: 100,
				},
			},
		},
	}

	testUserInfoString := `{"coins":1000,"inventory":[{"type":"test_item","quantity":3}],"coinHistory":{"received":[{"fromUser":"test_user","amount":100}],"sent":[{"toUser":"test_user","amount":100}]}}`

	type mockBehavior func(
		serviceAuth *mock_service.MockAuthorization,
		serviseInfo *mock_service.MockInfo,
	)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name: "OK",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(testUserInfoStruct, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: testUserInfoString,
		},
		{
			name:                 "Unauthorized - No Token",
			mockBehavior:         func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name: "Unauthorized - Invalid Token",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("invalid_token").Return(0, customerrors.ErrParseToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseToken),
		},
		{
			name: "Get balance error",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(models.UserInfoResponse{}, customerrors.ErrGetBalance)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrGetBalance),
		},
		{
			name: "Db transaction start error",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(models.UserInfoResponse{}, customerrors.ErrTxStart)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStart),
		},
		{
			name: "Db transaction stop error",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(models.UserInfoResponse{}, customerrors.ErrTxStop)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStop),
		},
		{
			name: "Data base error",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(models.UserInfoResponse{}, customerrors.ErrDatabase)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrDatabase),
		},
		{
			name: "Parse inventory error",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(models.UserInfoResponse{}, customerrors.ErrParseInventory)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseInventory),
		},
		{
			name: "Parse transaction history error",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviseInfo.EXPECT().GetUserInfo(1).Return(models.UserInfoResponse{}, customerrors.ErrParseTrx)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseTrx),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			info := mock_service.NewMockInfo(c)

			services := &service.Service{Authorization: auth, Info: info}
			hendler := NewHandler(services)

			testCase.mockBehavior(auth, info)

			r := gin.New()

			r.GET("/info", hendler.userIdentity, hendler.getInfo)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/info", nil)

			if testCase.name == "Unauthorized - Invalid Token" {
				req.Header.Set("Authorization", "Bearer invalid_token")

			} else if testCase.name != "Unauthorized - No Token" {
				req.Header.Set("Authorization", "Bearer valid_token")
			}

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
