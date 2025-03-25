package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/Sm3underscore23/merchStore/internal/service"
	mock_service "github.com/Sm3underscore23/merchStore/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestHandler_info(t *testing.T) {
	commonMockBehavior := func(userInfoStruct models.UserInfoResponse, returnErr error) func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
		return func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
			serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
			serviseInfo.EXPECT().GetUserInfo(1).Return(userInfoStruct, returnErr)
		}
	}

	testUserInfoStruct := models.UserInfoResponse{
		Balance: 1000,
		Inventory: []models.InventoryItem{
			{
				ItemType: "test_item",
				Quantity: 3,
			},
		},
		CoinHistory: models.TransactionHistory{
			Received: []models.IncomingTransaction{
				{
					FromUser: "test_user",
					Amount:   100,
				},
			},
			Sent: []models.OutgoingTransaction{
				{
					ToUser: "test_user",
					Amount: 100,
				},
			},
		},
	}

	testUserInfoString := `{"coins":1000,"inventory":[{"type":"test_item","quantity":3}],"coinHistory":{"received":[{"fromUser":"test_user","amount":100}],"sent":[{"toUser":"test_user","amount":100}]}}`

	testTable := []struct {
		name         string
		mockBehavior func(
			serviceAuth *mock_service.MockAuthorization,
			serviseInfo *mock_service.MockInfo,
		)
		noToken              bool
		invalidToken         bool
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name:                 "OK",
			mockBehavior:         commonMockBehavior(testUserInfoStruct, nil),
			expectedStatusCode:   200,
			expectedResponseBody: testUserInfoString,
		},
		{
			name:                 "Unauthorized - No Token",
			noToken:              true,
			mockBehavior:         func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name:         "Unauthorized - Invalid Token",
			invalidToken: true,
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviseInfo *mock_service.MockInfo) {
				serviceAuth.EXPECT().ParseToken("invalid_token").Return(0, customerrors.ErrParseToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseToken),
		},
		{
			name:                 "Get balance error",
			mockBehavior:         commonMockBehavior(models.UserInfoResponse{}, customerrors.ErrGetBalance),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrGetBalance),
		},
		{
			name:                 "Db transaction start error",
			mockBehavior:         commonMockBehavior(models.UserInfoResponse{}, customerrors.ErrTxStart),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStart),
		},
		{
			name:                 "Db transaction stop error",
			mockBehavior:         commonMockBehavior(models.UserInfoResponse{}, customerrors.ErrTxStop),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStop),
		},
		{
			name:                 "Data base error",
			mockBehavior:         commonMockBehavior(models.UserInfoResponse{}, customerrors.ErrDatabase),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrDatabase),
		},
		{
			name:                 "Parse inventory error",
			mockBehavior:         commonMockBehavior(models.UserInfoResponse{}, customerrors.ErrParseInventory),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseInventory),
		},
		{
			name:                 "Parse transaction history error",
			mockBehavior:         commonMockBehavior(models.UserInfoResponse{}, customerrors.ErrParseTrx),
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

			if testCase.invalidToken {
				req.Header.Set("Authorization", "Bearer invalid_token")

			} else if !testCase.noToken {
				req.Header.Set("Authorization", "Bearer valid_token")
			}

			r.ServeHTTP(w, req)

			fmt.Println(w.Code, testCase.expectedStatusCode)

			fmt.Println(w.Body.String(), testCase.expectedResponseBody)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)

			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}

}
