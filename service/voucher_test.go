package service_test

import (
	"testing"
	"time"
	"voucher_system/models"
	"voucher_system/repository"
	"voucher_system/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockVoucherRepository struct {
	mock.Mock
}

func (m *MockVoucherRepository) FindAll(userID int, voucherType string) ([]*models.Voucher, error) {
	args := m.Called(userID, voucherType)
	return args.Get(0).([]*models.Voucher), args.Error(1)
}

func (m *MockVoucherRepository) FindValidVoucher(userID int, voucherCode, area string, transactionAmount, shippingAmount float64, paymentMethod string, transactionDate time.Time) (*models.Voucher, error) {
	args := m.Called(userID, voucherCode, area, transactionAmount, shippingAmount, paymentMethod, transactionDate)
	return args.Get(0).(*models.Voucher), args.Error(1)
}
func (m *MockVoucherRepository) UpdateVoucherQuota(voucherID int, quota int) error {
	args := m.Called(voucherID, quota)
	return args.Error(0)
}

type MockHistoryRepository struct {
	mock.Mock
}

func (m *MockHistoryRepository) CreateHistory(history *models.History) error {
	args := m.Called(history)
	return args.Error(0)
}
func (m *MockHistoryRepository) FindUsageHistoryByUser(userID int) ([]models.History, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.History), args.Error(1)
}

type MockRedeemRepository struct {
	mock.Mock
}

func (m *MockRedeemRepository) FindUsersByVoucherCode(voucherCode string) ([]models.Redeem, error) {
	args := m.Called(voucherCode)
	return args.Get(0).([]models.Redeem), args.Error(1)
}
func (m *MockRedeemRepository) FindRedeemHistoryByUser(userID int) ([]models.Redeem, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Redeem), args.Error(1)
}

type MockRepository struct {
	mock.Mock
	VoucherRepo repository.VoucherRepository
	HistoryRepo repository.HistoryRepository
	RedeemRepo  repository.RedeemRepository
}

func (m *MockRepository) Voucher() repository.VoucherRepository {
	return m.VoucherRepo
}

func (m *MockRepository) Redeem() repository.RedeemRepository {
	return m.RedeemRepo
}
func (m *MockRepository) History() repository.HistoryRepository {
	return m.HistoryRepo
}

func TestVoucherService_FindVouchers(t *testing.T) {
	mockVoucherRepo := new(MockVoucherRepository)
	mockHistoryRepo := new(MockHistoryRepository)

	mockRepo := &repository.Repository{
		Voucher: mockVoucherRepo,
		History: mockHistoryRepo,
	}

	logger := zap.NewNop()
	service := service.NewVoucherService(*mockRepo, logger)

	userID := 1
	voucherType := "e-commerce"

	mockVoucher := &models.Voucher{
		ID:            1,
		VoucherCode:   "VOUCHER1",
		VoucherType:   voucherType,
		Quota:         100,
		DiscountValue: 10,
	}

	mockVoucherRepo.On("FindAll", userID, voucherType).Return([]*models.Voucher{mockVoucher}, nil)

	vouchers, err := service.FindVouchers(userID, voucherType)

	mockVoucherRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Len(t, vouchers, 1)
	assert.Equal(t, "VOUCHER1", vouchers[0].VoucherCode)
}

func TestVoucherService_FindVouchers_NoVouchers(t *testing.T) {
	mockVoucherRepo := new(MockVoucherRepository)
	mockHistoryRepo := new(MockHistoryRepository)

	mockRepo := &repository.Repository{
		Voucher: mockVoucherRepo,
		History: mockHistoryRepo,
	}

	logger := zap.NewNop()
	service := service.NewVoucherService(*mockRepo, logger)

	userID := 1
	voucherType := "e-commerce"

	mockVoucherRepo.On("FindAll", userID, voucherType).Return([]*models.Voucher{}, nil)

	vouchers, err := service.FindVouchers(userID, voucherType)

	mockVoucherRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Nil(t, vouchers)
	assert.Equal(t, "no vouchers available", err.Error())
}

func TestVoucherService_ValidateVoucher_Success(t *testing.T) {
	mockVoucherRepo := new(MockVoucherRepository)
	mockHistoryRepo := new(MockHistoryRepository)

	mockRepo := &repository.Repository{
		Voucher: mockVoucherRepo,
		History: mockHistoryRepo,
	}

	logger := zap.NewNop()
	service := service.NewVoucherService(*mockRepo, logger)

	userID := 1
	voucherCode := "VOUCHER1"
	transactionAmount := 100.0
	shippingAmount := 20.0
	area := "area1"
	paymentMethod := "credit"
	transactionDate := time.Now()

	mockVoucher := &models.Voucher{
		ID:              1,
		VoucherCode:     voucherCode,
		VoucherCategory: "Discount",
		DiscountValue:   10.0,
		Quota:           100,
	}

	mockVoucherRepo.On("FindValidVoucher", userID, voucherCode, area, transactionAmount, shippingAmount, paymentMethod, transactionDate).
		Return(mockVoucher, nil)

	voucher, benefitValue, err := service.ValidateVoucher(userID, voucherCode, transactionAmount, shippingAmount, area, paymentMethod, transactionDate)

	mockVoucherRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)

	assert.NoError(t, err)
	assert.NotNil(t, voucher)
	assert.Equal(t, 10.0, benefitValue)
}

// func TestVoucherService_UseVoucher_Success(t *testing.T) {
// 	mockVoucherRepo := new(MockVoucherRepository)
// 	mockHistoryRepo := new(MockHistoryRepository)
// 	mockRedeemRepo := new(repository.MockRedeemRepository)

// 	mockRepo := &repository.Repository{
// 		Voucher: mockVoucherRepo,
// 		History: mockHistoryRepo,
// 		Redeem:  mockRedeemRepo,
// 	}

// 	logger := zap.NewNop()
// 	service := service.NewVoucherService(*mockRepo, logger)

// 	userID := 1
// 	voucherCode := "VOUCHER1"
// 	transactionAmount := 100.0
// 	paymentMethod := "credit"
// 	area := "area1"

// 	mockVoucher := &models.Voucher{
// 		ID:              1,
// 		VoucherCode:     voucherCode,
// 		VoucherCategory: "Discount",
// 		DiscountValue:   10.0,
// 		Quota:           100,
// 	}

// 	benefitValue := 10.0

// 	mockVoucherRepo.On("FindValidVoucher", userID, voucherCode, area, transactionAmount, 0.0, paymentMethod, mock.Anything).
// 		Return(mockVoucher, benefitValue, nil)

// 	mockHistoryRepo.On("CreateHistory", mock.AnythingOfType("*models.History")).Return(nil)

// 	mockVoucherRepo.On("UpdateVoucherQuota", mock.Anything, mock.Anything).Return(nil)

// 	err := service.UseVoucher(userID, voucherCode, transactionAmount, paymentMethod, area)

// 	mockVoucherRepo.AssertExpectations(t)
// 	mockHistoryRepo.AssertExpectations(t)
// 	mockRedeemRepo.AssertExpectations(t)

// 	assert.NoError(t, err)
// }

// func TestVoucherService_UseVoucher_QuotaExceeded(t *testing.T) {
// 	mockVoucherRepo := new(MockVoucherRepository)
// 	mockHistoryRepo := new(MockHistoryRepository)
// 	mockRedeemRepo := new(repository.MockRedeemRepository)

// 	mockRepo := &repository.Repository{
// 		Voucher: mockVoucherRepo,
// 		History: mockHistoryRepo,
// 		Redeem:  mockRedeemRepo,
// 	}

// 	logger := zap.NewNop()
// 	service := service.NewVoucherService(*mockRepo, logger)

// 	userID := 1
// 	voucherCode := "VOUCHER1"
// 	transactionAmount := 100.0
// 	paymentMethod := "credit"
// 	area := "area1"

// 	mockVoucher := &models.Voucher{
// 		ID:              1,
// 		VoucherCode:     voucherCode,
// 		VoucherCategory: "Discount",
// 		DiscountValue:   10.0,
// 		Quota:           0, // This will trigger the "quota exceeded" scenario
// 	}

// 	mockVoucherRepo.On("FindValidVoucher", userID, voucherCode, area, transactionAmount, 0.0, paymentMethod, mock.Anything).
//     Return(mockVoucher, 10.0, nil)

// 	mockHistoryRepo.On("CreateHistory", mock.AnythingOfType("*models.History")).Return(nil)

// 	mockVoucherRepo.On("UpdateVoucherQuota", mock.Anything, mock.Anything).Return(nil)

// 	err := service.UseVoucher(userID, voucherCode, transactionAmount, paymentMethod, area)

// 	mockVoucherRepo.AssertExpectations(t)
// 	mockHistoryRepo.AssertExpectations(t)
// 	mockRedeemRepo.AssertExpectations(t)

// 	assert.Error(t, err)
// 	assert.Equal(t, "voucher quota exceeded", err.Error())
// }

func TestHistoryService_GetRedeemHistoryByUser(t *testing.T) {
	mockRedeemRepo := new(MockRedeemRepository)
	mockHistoryRepo := new(MockHistoryRepository)

	mockRepo := &repository.Repository{
		Redeem:  mockRedeemRepo,
		History: mockHistoryRepo,
	}

	logger := zap.NewNop()
	service := service.NewHistoryService(*mockRepo, logger)

	userID := 1
	mockRedeemHistory := []models.Redeem{
		{ID: 1, VoucherID: 1, UserID: userID},
	}

	mockRedeemRepo.On("FindRedeemHistoryByUser", userID).Return(mockRedeemHistory, nil)

	history, err := service.GetRedeemHistoryByUser(userID)

	mockRedeemRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, 1, history[0].VoucherID)
}

func TestHistoryService_GetUsageHistoryByUser(t *testing.T) {
	mockRedeemRepo := new(MockRedeemRepository)
	mockHistoryRepo := new(MockHistoryRepository)

	mockRepo := &repository.Repository{
		Redeem:  mockRedeemRepo,
		History: mockHistoryRepo,
	}

	logger := zap.NewNop()
	service := service.NewHistoryService(*mockRepo, logger)

	userID := 1
	mockUsageHistory := []models.History{
		{ID: 1, UserID: userID, VoucherID: 1, UsageDate: time.Now()},
	}

	mockHistoryRepo.On("FindUsageHistoryByUser", userID).Return(mockUsageHistory, nil)

	history, err := service.GetUsageHistoryByUser(userID)

	mockHistoryRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, 1, history[0].VoucherID)
}

func TestHistoryService_GetUsersByVoucherCode(t *testing.T) {
	mockRedeemRepo := new(MockRedeemRepository)
	mockHistoryRepo := new(MockHistoryRepository)

	mockRepo := &repository.Repository{
		Redeem:  mockRedeemRepo,
		History: mockHistoryRepo,
	}

	logger := zap.NewNop()
	service := service.NewHistoryService(*mockRepo, logger)

	voucherCode := "VOUCHER1"
	mockUsers := []models.Redeem{
		{ID: 1, UserID: 1, VoucherID: 1, RedeemDate: time.Now()},
	}

	mockRedeemRepo.On("FindUsersByVoucherCode", voucherCode).Return(mockUsers, nil)

	users, err := service.GetUsersByVoucherCode(voucherCode)

	mockRedeemRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, 1, users[0].VoucherID)
	assert.Equal(t, 1, users[0].UserID)
}
