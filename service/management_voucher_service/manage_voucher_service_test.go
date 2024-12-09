package managementvoucherservice_test

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"voucher_system/models"
	"voucher_system/repository"
	managementvoucher "voucher_system/repository/management_voucher"
	managementvoucherservice "voucher_system/service/management_voucher_service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateVoucher(t *testing.T) {
	log := *zap.NewNop()
	mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
	repo := repository.Repository{
		Manage: mockRepo,
	}
	service := managementvoucherservice.NewManagementVoucherService(repo, &log)

	voucher := &models.Voucher{
		VoucherName:    "Discount 10%",
		Quota:          100,
		PointsRequired: 50,
		StartDate:      time.Now().AddDate(0, 0, -5),
		EndDate:        time.Now().AddDate(0, 0, 5),
	}

	t.Run("Successfully create voucher", func(t *testing.T) {

		mockRepo.On("CreateVoucher", voucher).Return(nil)

		err := service.CreateVoucher(voucher)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "CreateVoucher", voucher)
	})

	t.Run("Fail to create voucher - insufficient data", func(t *testing.T) {

		invalidVoucher := &models.Voucher{
			VoucherName: "",
		}

		mockRepo.On("CreateVoucher", invalidVoucher).Return(fmt.Errorf("voucher data is invalid"))

		err := service.CreateVoucher(invalidVoucher)

		assert.Error(t, err)
		assert.EqualError(t, err, "voucher data is invalid")
		mockRepo.AssertCalled(t, "CreateVoucher", invalidVoucher)
	})
}

func TestSoftDeleteVoucher(t *testing.T) {
	log := *zap.NewNop() // Logger yang tidak menghasilkan output apapun untuk testing

	voucherID := 100

	t.Run("Successfully soft delete voucher", func(t *testing.T) {
		// Mock baru untuk setiap test case
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)

		// Setup mock untuk SoftDeleteVoucher yang berhasil (tidak mengembalikan error)
		mockRepo.On("SoftDeleteVoucher", voucherID).Return(nil)

		// Panggil method SoftDeleteVoucher di service
		err := service.SoftDeleteVoucher(voucherID)

		// Pastikan tidak ada error
		assert.NoError(t, err)

		// Pastikan bahwa metode SoftDeleteVoucher dipanggil dengan benar
		mockRepo.AssertCalled(t, "SoftDeleteVoucher", voucherID)
	})

	t.Run("Fail to soft delete voucher", func(t *testing.T) {
		// Mock baru untuk setiap test case
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)

		// Setup mock untuk mengembalikan error
		mockRepo.On("SoftDeleteVoucher", voucherID).Return(errors.New("voucher not found"))

		// Panggil service
		err := service.SoftDeleteVoucher(voucherID)

		// Verifikasi hasil
		assert.Error(t, err)
		assert.Equal(t, "voucher not found", err.Error())

		mockRepo.AssertCalled(t, "SoftDeleteVoucher", voucherID)
	})
}

func TestUpdateVoucher(t *testing.T) {
	log := *zap.NewNop()

	voucherID := 100
	voucher := &models.Voucher{
		VoucherName:    "Updated Discount 10%",
		Quota:          50,
		PointsRequired: 25,
	}

	t.Run("Successfully update voucher", func(t *testing.T) {
		// Mock baru untuk test case ini
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)

		// Setup mock
		mockRepo.On("UpdateVoucher", voucher, voucherID).Return(nil)

		// Panggil service
		err := service.UpdateVoucher(voucher, voucherID)

		// Verifikasi hasil
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "UpdateVoucher", voucher, voucherID)
	})

	t.Run("Fail to update voucher", func(t *testing.T) {
		// Mock baru untuk test case ini
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)

		// Setup mock
		mockRepo.On("UpdateVoucher", voucher, voucherID).Return(fmt.Errorf("voucher update failed"))

		// Panggil service
		err := service.UpdateVoucher(voucher, voucherID)

		// Verifikasi hasil
		assert.Error(t, err)
		assert.EqualError(t, err, "voucher update failed")
		mockRepo.AssertCalled(t, "UpdateVoucher", voucher, voucherID)
	})
}

func TestShowRedeemPoints(t *testing.T) {
	log := *zap.NewNop()

	expectedRedeemPoints := &[]managementvoucher.RedeemPoint{
		{VoucherName: "Discount 10%", PointsRequired: 50, DiscountValue: 10.0},
	}

	t.Run("Successfully show redeem points", func(t *testing.T) {
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)

		mockRepo.On("ShowRedeemPoints").Return(expectedRedeemPoints, nil)

		redeemPoints, err := service.ShowRedeemPoints()

		assert.NoError(t, err)
		assert.Equal(t, expectedRedeemPoints, redeemPoints)
		mockRepo.AssertCalled(t, "ShowRedeemPoints")
	})

	t.Run("Fail to show redeem points", func(t *testing.T) {
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)

		mockRepo.On("ShowRedeemPoints").Return(nil, fmt.Errorf("failed to fetch redeem points"))

		redeemPoints, err := service.ShowRedeemPoints()

		assert.Error(t, err)
		assert.Nil(t, redeemPoints)
		assert.EqualError(t, err, "failed to fetch redeem points")
		mockRepo.AssertCalled(t, "ShowRedeemPoints")
	})
}

func TestGetVouchersByQueryParams(t *testing.T) {
	log := *zap.NewNop()

	status := "active"
	area := "Jakarta"
	voucherType := "discount"

	expectedVouchers := &[]models.Voucher{
		{VoucherName: "Discount 10%", Quota: 100, PointsRequired: 50},
	}

	t.Run("Successfully get vouchers by query params", func(t *testing.T) {
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)
		mockRepo.On("GetVouchersByQueryParams", status, area, voucherType).Return(expectedVouchers, nil)

		vouchers, err := service.GetVouchersByQueryParams(status, area, voucherType)

		assert.NoError(t, err)
		assert.Equal(t, expectedVouchers, vouchers)
		mockRepo.AssertCalled(t, "GetVouchersByQueryParams", status, area, voucherType)
	})

	t.Run("Fail to get vouchers by query params", func(t *testing.T) {
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)
		mockRepo.On("GetVouchersByQueryParams", status, area, voucherType).Return(nil, fmt.Errorf("no vouchers found"))

		vouchers, err := service.GetVouchersByQueryParams(status, area, voucherType)

		assert.Error(t, err)
		assert.Nil(t, vouchers)
		assert.EqualError(t, err, "no vouchers found")
		mockRepo.AssertCalled(t, "GetVouchersByQueryParams", status, area, voucherType)
	})
}

func TestCreateRedeemVoucher(t *testing.T) {
	log := *zap.NewNop()

	redeem := &models.Redeem{
		UserID:    1,
		VoucherID: 100,
	}
	points := 50

	t.Run("Successfully create redeem voucher", func(t *testing.T) {
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)
		mockRepo.On("CreateRedeemVoucher", redeem, points).Return(nil)

		err := service.CreateRedeemVoucher(redeem, points)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "CreateRedeemVoucher", redeem, points)
	})

	t.Run("Fail to create redeem voucher", func(t *testing.T) {
		mockRepo := &managementvoucher.ManagementVoucherRepoMock{}
		repo := repository.Repository{
			Manage: mockRepo,
		}
		service := managementvoucherservice.NewManagementVoucherService(repo, &log)
		mockRepo.On("CreateRedeemVoucher", redeem, points).Return(fmt.Errorf("user already claimed this voucher"))

		err := service.CreateRedeemVoucher(redeem, points)

		assert.Error(t, err)
		assert.EqualError(t, err, "user already claimed this voucher")
		mockRepo.AssertCalled(t, "CreateRedeemVoucher", redeem, points)
	})
}
