package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/service"
	mock_service "github.com/Sm3underscore23/merchStore/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestHandler_buyItem(t *testing.T) {

	commonMockBehavior := func(returnErr error) func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
		return func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
			serviceAuth.EXPECT().ParseToken("valid_token").Return(1, nil)
			serviceBuy.EXPECT().Buy(1, "valid_product").Return(returnErr)
		}
	}

	testTable := []struct {
		name               string
		productSlug        string
		invalidToken       bool
		noToken            bool
		expectedStatusCode int
		mockBehavior       func(
			serviceAuth *mock_service.MockAuthorization,
			serviceBuy *mock_service.MockBuy,
		)
		expectedResponseBody string
	}{
		{
			name:               "OK",
			productSlug:        "valid_product",
			mockBehavior:       commonMockBehavior(nil),
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
			name:                 "Not enoghf coins",
			productSlug:          "valid_product",
			mockBehavior:         commonMockBehavior(customerrors.ErrChangeBalance),
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrChangeBalance),
		},
		{
			name:                 "Unauthorized - No Token",
			noToken:              true,
			productSlug:          "valid_product",
			mockBehavior:         func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name:         "Unauthorized - Invalid Token",
			invalidToken: true,
			productSlug:  "valid_product",
			mockBehavior: func(serviceAuth *mock_service.MockAuthorization, serviceBuy *mock_service.MockBuy) {
				serviceAuth.EXPECT().ParseToken("invalid_token").Return(0, customerrors.ErrParseToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrParseToken),
		},
		{
			name:                 "Data base error",
			productSlug:          "valid_product",
			mockBehavior:         commonMockBehavior(customerrors.ErrDatabase),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrDatabase),
		},
		{
			name:                 "Db transaction start error",
			productSlug:          "valid_product",
			mockBehavior:         commonMockBehavior(customerrors.ErrTxStart),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStart),
		},
		{
			name:                 "Db transaction stop error",
			productSlug:          "valid_product",
			mockBehavior:         commonMockBehavior(customerrors.ErrTxStop),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrTxStop),
		},
		{
			name:                 "Get balance error",
			productSlug:          "valid_product",
			mockBehavior:         commonMockBehavior(customerrors.ErrGetBalance),
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"errors":"%s"}`, customerrors.ErrGetBalance),
		},
		{
			name:                 "Update balance error",
			productSlug:          "valid_product",
			mockBehavior:         commonMockBehavior(customerrors.ErrUpdateBalance),
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

			if testCase.invalidToken {
				req.Header.Set("Authorization", "Bearer invalid_token")

			} else if !testCase.noToken {
				req.Header.Set("Authorization", "Bearer valid_token")
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)

			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
