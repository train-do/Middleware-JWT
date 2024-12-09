package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"voucher_system/controller"
	"voucher_system/models"
	"voucher_system/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockVoucherService struct {
	mock.Mock
}

func (m *MockVoucherService) FindVouchers(userID int, voucherType string) ([]*models.Voucher, error) {
	args := m.Called(userID, voucherType)
	return args.Get(0).([]*models.Voucher), args.Error(1)
}

func (m *MockVoucherService) ValidateVoucher(userID int, voucherCode string, transactionAmount float64, shippingAmount float64, area string, paymentMethod string, transactionDate time.Time) (*models.Voucher, float64, error) {
	args := m.Called(userID, voucherCode, transactionAmount, shippingAmount, area, paymentMethod, transactionDate)
	return args.Get(0).(*models.Voucher), args.Get(1).(float64), args.Error(2)
}

func (m *MockVoucherService) UseVoucher(userID int, voucherCode string, transactionAmount float64, paymentMethod string, area string) error {
	args := m.Called(userID, voucherCode, transactionAmount, paymentMethod, area)
	return args.Error(0)
}

type MockHistoryService struct {
	mock.Mock
}

func (m *MockHistoryService) GetRedeemHistoryByUser(userID int) ([]models.Redeem, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Redeem), args.Error(1)
}

func (m *MockHistoryService) GetUsageHistoryByUser(userID int) ([]models.History, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.History), args.Error(1)
}

func (m *MockHistoryService) GetUsersByVoucherCode(voucherCode string) ([]models.Redeem, error) {
	args := m.Called(voucherCode)
	return args.Get(0).([]models.Redeem), args.Error(1)
}

type MockService struct {
	mock.Mock
	VoucherService service.VoucherService
	HistoryService service.HistoryService
}

func (m *MockService) Voucher() service.VoucherService {
	return m.VoucherService
}
func (m *MockService) History() service.HistoryService {
	return m.HistoryService
}

func TestVoucherController_FindVouchers(t *testing.T) {
	mockVoucherService := new(MockVoucherService)
	mockHistoryService := new(MockHistoryService)

	mockService := &service.Service{
		Voucher: mockVoucherService,
		History: mockHistoryService,
	}

	logger := zap.NewNop()

	controller := controller.NewVoucherController(*mockService, logger)

	userID := 1
	voucherType := "e-commerce"

	mockVoucherService.On("FindVouchers", userID, voucherType).Return([]*models.Voucher{
		{ID: 1, VoucherCode: "VOUCHER1", VoucherType: "e-commerce", DiscountValue: 10},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/vouchers?type=e-commerce", nil)

	controller.FindVouchers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockVoucherService.AssertExpectations(t)
}

func TestVoucherController_ValidateVoucher(t *testing.T) {
	mockVoucherService := new(MockVoucherService)
	mockHistoryService := new(MockHistoryService)

	mockService := &service.Service{
		Voucher: mockVoucherService,
		History: mockHistoryService,
	}

	logger := zap.NewNop()

	controller := controller.NewVoucherController(*mockService, logger)

	userID := 1
	voucherCode := "VOUCHER1"
	transactionAmount := 100.0
	shippingAmount := 10.0
	area := "area1"
	paymentMethod := "credit_card"
	transactionDate := time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC)
	voucher := models.Voucher{ID: 1, VoucherCode: voucherCode, Status: true}
	benefitValue := 20.0

	mockVoucherService.On("ValidateVoucher", userID, voucherCode, transactionAmount, shippingAmount, area, paymentMethod, transactionDate).Return(&voucher, benefitValue, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/validate-voucher", strings.NewReader(`{
		"voucher_code": "VOUCHER1",
		"transaction_amount": 100,
		"shipping_amount": 10,
		"area": "area1",
		"payment_method": "credit_card",
		"transaction_date": "2024-12-01"
	}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ValidateVoucher(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockVoucherService.AssertExpectations(t)
}

func TestVoucherController_UseVoucher(t *testing.T) {
	mockVoucherService := new(MockVoucherService)
	mockHistoryService := new(MockHistoryService)

	mockService := &service.Service{
		Voucher: mockVoucherService,
		History: mockHistoryService,
	}

	logger := zap.NewNop()

	controller := controller.NewVoucherController(*mockService, logger)

	userID := 1
	voucherCode := "VOUCHER1"
	transactionAmount := 100.0
	area := "area1"
	paymentMethod := "credit_card"

	mockVoucherService.On("UseVoucher", userID, voucherCode, transactionAmount, paymentMethod, area).
		Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/use-voucher", strings.NewReader(`{
		"user_id": 1,
		"voucher_code": "VOUCHER1",
		"transaction_amount": 100,
		"payment_method": "credit_card",  
		"area": "area1"
	}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UseVoucher(c)

	fmt.Println("Response Code:", w.Code)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "voucher used successfully", response["message"])

	mockVoucherService.AssertExpectations(t)
}

func TestVoucherController_GetRedeemHistoryByUser(t *testing.T) {
	mockHistoryService := new(MockHistoryService)
	mockVoucherService := new(MockVoucherService)

	mockService := &service.Service{
		Voucher: mockVoucherService,
		History: mockHistoryService,
	}

	logger := zap.NewNop()

	controller := controller.NewVoucherController(*mockService, logger)

	userID := 1
	redeemHistory := []models.Redeem{
		{ID: 1, UserID: userID, VoucherID: 1, RedeemDate: time.Now()},
	}

	mockHistoryService.On("GetRedeemHistoryByUser", userID).Return(redeemHistory, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/redeem-history", nil)

	controller.GetRedeemHistoryByUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockHistoryService.AssertExpectations(t)
}

func TestVoucherController_GetUsersByVoucherCode(t *testing.T) {

	mockHistoryService := new(MockHistoryService)
	mockService := &service.Service{
		History: mockHistoryService,
	}

	logger := zap.NewNop()

	controller := controller.NewVoucherController(*mockService, logger)

	voucherCode := "VOUCHER1"
	expectedUsers := []models.Redeem{
		{ID: 1, UserID: 1, VoucherID: 1, RedeemDate: time.Now()},
		{ID: 2, UserID: 1, VoucherID: 2, RedeemDate: time.Now()},
	}

	mockHistoryService.On("GetUsersByVoucherCode", voucherCode).
		Return(expectedUsers, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "voucher_code", Value: voucherCode}}

	controller.GetUsersByVoucherCode(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Users fetched successfully", response["message"])
	assert.NotNil(t, response["data"])

	mockHistoryService.AssertExpectations(t)
}
