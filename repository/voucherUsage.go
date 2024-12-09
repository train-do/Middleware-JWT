package repository

import (
	"encoding/json"
	"fmt"
	"time"
	"voucher_system/helper"
	"voucher_system/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VoucherRepository interface {
	FindAll(userID int, voucherType string) ([]*models.Voucher, error)
	FindValidVoucher(userID int, voucherCode, area string, transactionAmount, shippingAmount float64, paymentMethod string, transactionDate time.Time) (*models.Voucher, error)
	UpdateVoucherQuota(voucherID int, quota int) error
}

type voucherRepository struct {
	DB  *gorm.DB
	log *zap.Logger
}

func NewVoucherRepository(db *gorm.DB, log *zap.Logger) VoucherRepository {
	return &voucherRepository{DB: db, log: log}
}

func (r *voucherRepository) FindAll(userID int, voucherType string) ([]*models.Voucher, error) {
	r.log.Info("Fetching all vouchers", zap.Int("userID", userID), zap.String("voucherType", voucherType))
	var rawVouchers []struct {
		models.Voucher
		RawPaymentMethods  []byte `gorm:"column:payment_methods"`
		RawApplicableAreas []byte `gorm:"column:applicable_areas"`
	}

	query := r.DB.
		Table("vouchers").
		Select(`vouchers.*, vouchers.payment_methods AS raw_payment_methods, vouchers.applicable_areas AS raw_applicable_areas`).
		Joins("JOIN redeems ON redeems.voucher_id = vouchers.id").
		Where("redeems.user_id = ? AND vouchers.status = ?", userID, true)

	if voucherType != "" {
		query = query.Where("vouchers.voucher_type = ?", voucherType)
	}

	err := query.Find(&rawVouchers).Error
	if err != nil {
		r.log.Error("Error fetching vouchers", zap.Error(err))
		return nil, err
	}

	r.log.Info("Fetched vouchers successfully", zap.Int("voucherCount", len(rawVouchers)))
	vouchers := make([]*models.Voucher, 0, len(rawVouchers))
	for _, rawVoucher := range rawVouchers {
		v := rawVoucher.Voucher

		if len(rawVoucher.RawPaymentMethods) > 0 {
			if err := json.Unmarshal(rawVoucher.RawPaymentMethods, &v.PaymentMethods); err != nil {
				r.log.Error("Error unmarshalling payment methods", zap.Error(err))
				return nil, err
			}
		}

		if len(rawVoucher.RawApplicableAreas) > 0 {
			if err := json.Unmarshal(rawVoucher.RawApplicableAreas, &v.ApplicableAreas); err != nil {
				r.log.Error("Error unmarshalling applicable areas", zap.Error(err))
				return nil, err
			}
		}
		vouchers = append(vouchers, &v)

	}

	return vouchers, nil
}

func (r *voucherRepository) FindValidVoucher(userID int, voucherCode, area string, transactionAmount, shippingAmount float64, paymentMethod string, transactionDate time.Time) (*models.Voucher, error) {
	r.log.Info("Finding valid voucher", zap.Int("userID", userID), zap.String("voucherCode", voucherCode), zap.String("area", area))

	var rawVoucher struct {
		models.Voucher
		RawPaymentMethods  []byte `gorm:"column:payment_methods"`
		RawApplicableAreas []byte `gorm:"column:applicable_areas"`
	}

	err := r.DB.Table("vouchers").
		Select(`vouchers.*`).
		Joins("JOIN redeems ON redeems.voucher_id = vouchers.id").
		Where(`
			redeems.user_id = ? AND
			vouchers.voucher_code = ? AND 
			quota > 0`,
			userID, voucherCode).First(&rawVoucher).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("voucher not found")
		}
		r.log.Error("Error finding valid voucher", zap.String("voucherCode", voucherCode), zap.Error(err))
		return nil, err
	}

	if transactionAmount < rawVoucher.Voucher.MinimumPurchase {
		return nil, fmt.Errorf("transaction amount must be at least %.2f", rawVoucher.Voucher.MinimumPurchase)
	}

	if len(rawVoucher.RawApplicableAreas) > 0 {
		var areas []string
		if err := json.Unmarshal(rawVoucher.RawApplicableAreas, &areas); err != nil {
			r.log.Error("Error unmarshalling applicable areas", zap.Error(err))
			return nil, fmt.Errorf("failed to parse applicable_areas: %w", err)
		}
		if !helper.Contains(areas, area) {
			return nil, fmt.Errorf("area not found")
		}
	}

	if len(rawVoucher.RawPaymentMethods) > 0 {
		var paymentMethods []string
		if err := json.Unmarshal(rawVoucher.RawPaymentMethods, &paymentMethods); err != nil {
			r.log.Error("Error unmarshalling payment methods", zap.Error(err))
			return nil, fmt.Errorf("failed to parse payment_methods: %w", err)
		}
		if !helper.Contains(paymentMethods, paymentMethod) {
			return nil, fmt.Errorf("payment method not found")
		}
	}

	if transactionDate.Before(rawVoucher.Voucher.StartDate) {
		return nil, fmt.Errorf("voucher not available yet")
	}
	if transactionDate.After(rawVoucher.Voucher.EndDate) {
		return nil, fmt.Errorf("voucher expired")
	}

	if len(rawVoucher.RawPaymentMethods) > 0 {
		if err := json.Unmarshal(rawVoucher.RawPaymentMethods, &rawVoucher.PaymentMethods); err != nil {
			r.log.Error("Error unmarshalling payment methods", zap.Error(err))
			return nil, fmt.Errorf("failed to parse payment_methods: %w", err)
		}
	}

	if len(rawVoucher.RawApplicableAreas) > 0 {
		if err := json.Unmarshal(rawVoucher.RawApplicableAreas, &rawVoucher.ApplicableAreas); err != nil {
			r.log.Error("Error unmarshalling applicable areas", zap.Error(err))
			return nil, fmt.Errorf("failed to parse applicable_areas: %w", err)
		}
	}

	r.log.Info("Valid voucher found", zap.String("voucherCode", voucherCode))
	return &rawVoucher.Voucher, nil
}

func (r *voucherRepository) UpdateVoucherQuota(voucherID int, quota int) error {
	r.log.Info("Updating voucher quota", zap.Int("voucherID", voucherID), zap.Int("quota", quota))
	err := r.DB.Model(&models.Voucher{}).Where("id = ?", voucherID).Update("quota", quota).Error
	if err != nil {
		r.log.Error("Error updating voucher quota", zap.Int("voucherID", voucherID), zap.Error(err))
		return err
	}
	r.log.Info("Voucher quota updated", zap.Int("voucherID", voucherID), zap.Int("newQuota", quota))
	return nil
}
