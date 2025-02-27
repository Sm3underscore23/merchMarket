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

func TestHandler_sendCoins(t *testing.T) {
	testInputBody := `{"toUser":"test_username","amount":100}`
	testInputUser := models.SendCoinRequest{
		Receiver: "test_username",
		Coins:    100,
	}

	type mockBehavior func(
		serviceAuth *mock_service.MockAuthorization,
		serviceSC *mock_service.MockSendCoins,
		sendCoinsRequest models.SendCoinRequest,
	)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            models.SendCoinRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(sendCoinsRequest.Receiver, 1, sendCoinsRequest.Coins).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"toUser":"test_username"}`,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(
				`{"errors":"%s"}`,
				customerrors.ErrInvalidInputBody,
			),
		},
		{
			name:      "ToUser not found",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrUserNotFound)
			},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrUserNotFound),
		},
		{
			name:      "Not enoghf coins",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrChangeBalance)
			},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrChangeBalance),
		},
		{
			name:      "Unauthorized - No Token",
			inputBody: testInputBody,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
			}, // Не мокаем сервис, так как middleware сразу вернет 401
			expectedStatusCode:   401,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name:      "Unauthorized - Invalid Token",
			inputBody: testInputBody,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("invalid_token").Return(0, customerrors.ErrParseToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseToken),
		},

		{
			name:      "Get balance error",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrGetBalance)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrGetBalance),
		},
		{
			name:      "Update balance error",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrUpdateBalance)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrUpdateBalance),
		},
		{
			name:      "Db transaction start error",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrTxStart)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStart),
		},
		{
			name:      "Db transaction stop error",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrTxStop)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStop),
		},
		{
			name:      "Add transaction to list error",
			inputBody: testInputBody,
			inputUser: testInputUser,
			mockBehavior: func(
				serviceAuth *mock_service.MockAuthorization,
				serviceSC *mock_service.MockSendCoins,
				sendCoinsRequest models.SendCoinRequest,
			) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceSC.EXPECT().SendCoins(
					sendCoinsRequest.Receiver,
					1, sendCoinsRequest.Coins,
				).Return(customerrors.ErrAddTrxToList)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrAddTrxToList),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			snd := mock_service.NewMockSendCoins(c)

			testCase.mockBehavior(auth, snd, testCase.inputUser)

			services := &service.Service{Authorization: auth, SendCoins: snd}
			hendler := NewHandler(services)

			r := gin.New()

			r.POST("/sendCoin", hendler.userIdentity, hendler.sendCoins, func(ctx *gin.Context) {})

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/sendCoin", bytes.NewBufferString(testCase.inputBody))

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
