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

func TestHandler_buyItem(t *testing.T) {
	type mockBehavior func(
		serviceAuth *mock_service.MockAuthorization,
		serviceBuy *mock_service.MockBuy,
	)

	testTable := []struct {
		name                 string
		productSlug          string
		expectedStatusCode   int
		mockBehavior         mockBehavior
		expectedResponseBody string
	}{
		{
			name:        "OK",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:        "Invalid product slug",
			productSlug: "invalid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "invalid_product").Return(customerrors.ErrProductNotFound)
			},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrProductNotFound),
		},
		{
			name:        "Not enoghf coins",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(customerrors.ErrChangeBalance)
			},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrChangeBalance),
		},
		{
			name:                 "Unauthorized - No Token",
			productSlug:          "valid_product",
			mockBehavior:         func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name:        "Unauthorized - Invalid Token",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("invalid_token").Return(0, customerrors.ErrParseToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseToken),
		},
		{
			name:        "Data base error",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(customerrors.ErrDatabase)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrDatabase),
		},
		{
			name:        "Db transaction start error",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(customerrors.ErrTxStart)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStart),
		},
		{
			name:        "Db transaction stop error",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(customerrors.ErrTxStop)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStop),
		},
		{
			name:        "Get balance error",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(customerrors.ErrGetBalance)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrGetBalance),
		},
		{
			name:        "Update balance error",
			productSlug: "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
				serviceBuy.EXPECT().Buy(1, "valid_product").Return(customerrors.ErrUpdateBalance)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrUpdateBalance),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			buy := mock_service.NewMockBuy(c)

			testCase.mockBehavior(auth, buy)

			services := &service.Service{Authorization: auth, Buy: buy}
			hendler := NewHandler(services)

			r := gin.New()

			r.GET("/buy/:slug", hendler.userIdentity, hendler.buyItem)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/buy/%s", testCase.productSlug), nil)

			if testCase.name == "Unauthorized - Invalid Token" {
				req.Header.Set("Authorization", "Bearer invalid_token")

			} else if testCase.name != "Unauthorized - No Token" {
				req.Header.Set("Authorization", "Bearer valid_token")
			}

			r.ServeHTTP(w, req)

			fmt.Println(testCase.expectedStatusCode, w.Code)
			fmt.Println(testCase.expectedResponseBody, w.Body.String())

			if testCase.expectedStatusCode != w.Code {
				t.Error("stutus codes are different")
			}

			if testCase.expectedResponseBody != w.Body.String() {
				t.Error("response bodies are different")
			}
		})
	}
}
