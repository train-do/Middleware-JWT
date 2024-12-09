package service

import (
	"errors"
	"fmt"
	"time"
	"voucher_system/models"
	"voucher_system/repository"

	"go.uber.org/zap"
)

type VoucherService interface {
	FindVouchers(userID int, voucherType string) ([]*models.Voucher, error)
	ValidateVoucher(userID int, voucherCode string, transactionAmount float64, shippingAmount float64, area string, paymentMethod string, transactionDate time.Time) (*models.Voucher, float64, error)
	UseVoucher(userID int, voucherCode string, transactionAmount float64, paymentMethod string, area string) error
}

type voucherService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewVoucherService(repo repository.Repository, log *zap.Logger) VoucherService {
	return &voucherService{
		repo: repo,
		log:  log,
	}
}

func (s *voucherService) FindVouchers(userID int, voucherType string) ([]*models.Voucher, error) {
	s.log.Info("Finding vouchers", zap.Int("userID", userID), zap.String("voucherType", voucherType))
	vouchers, err := s.repo.Voucher.FindAll(userID, voucherType)
	if err != nil {
		s.log.Error("Error finding vouchers", zap.Error(err))
		return nil, err
	}
	if len(vouchers) == 0 {
		s.log.Info("No vouchers available", zap.Int("userID", userID))
		return nil, errors.New("no vouchers available")
	}

	s.log.Info("Vouchers found", zap.Int("voucherCount", len(vouchers)))
	return vouchers, nil
}

func (s *voucherService) ValidateVoucher(userID int, voucherCode string, transactionAmount float64, shippingAmount float64, area string, paymentMethod string, transactionDate time.Time) (*models.Voucher, float64, error) {
	s.log.Info("Validating voucher", zap.Int("userID", userID), zap.String("voucherCode", voucherCode), zap.Float64("transactionAmount", transactionAmount))

	voucher, err := s.repo.Voucher.FindValidVoucher(userID, voucherCode, area, transactionAmount, shippingAmount, paymentMethod, transactionDate)
	if err != nil {
		s.log.Error("Voucher validation failed", zap.String("voucherCode", voucherCode), zap.Error(err))
		return nil, 0, err
	}

	var benefitValue float64
	discountValueStr := fmt.Sprintf("%.0f", voucher.DiscountValue)
	if voucher.VoucherCategory == "Free Shipping" {
		benefitValue = shippingAmount
	} else if len(discountValueStr) > 4 {
		benefitValue = voucher.DiscountValue
	} else {
		benefitValue = (voucher.DiscountValue / 100) * transactionAmount
	}

	s.log.Info("Voucher validated successfully", zap.String("voucherCode", voucherCode), zap.Float64("benefitValue", benefitValue))
	return voucher, benefitValue, nil
}

func (s *voucherService) UseVoucher(userID int, voucherCode string, transactionAmount float64, paymentMethod string, area string) error {

	voucher, benefitValue, err := s.ValidateVoucher(userID, voucherCode, transactionAmount, 0, area, paymentMethod, time.Now())
	if err != nil {
		return err
	}

	history := &models.History{
		UserID:            userID,
		VoucherID:         voucher.ID,
		TransactionAmount: transactionAmount,
		BenefitValue:      benefitValue,
		UsageDate:         time.Now(),
	}

	err = s.repo.History.CreateHistory(history)
	if err != nil {
		return err
	}

	newQuota := voucher.Quota - 1
	if newQuota < 0 {
		return fmt.Errorf("voucher quota exceeded")
	}

	err = s.repo.Voucher.UpdateVoucherQuota(voucher.ID, newQuota)
	if err != nil {
		return err
	}

	return nil
}
