package repository

import (
	"errors"
	"voucher_system/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RedeemRepository interface {
	FindUsersByVoucherCode(voucherCode string) ([]models.Redeem, error)
	FindRedeemHistoryByUser(userID int) ([]models.Redeem, error)
}

type redeemRepository struct {
	DB  *gorm.DB
	log *zap.Logger
}

func NewRedeemRepository(db *gorm.DB, log *zap.Logger) RedeemRepository {
	return &redeemRepository{DB: db, log: log}
}

func (r *redeemRepository) FindUsersByVoucherCode(voucherCode string) ([]models.Redeem, error) {
	var redeems []models.Redeem

	r.log.Info("Fetching users by voucher code", zap.String("voucher_code", voucherCode))

	err := r.DB.Table("redeems").
		Select("redeems.*").
		Joins("JOIN vouchers ON vouchers.id = redeems.voucher_id").
		Where("vouchers.voucher_code = ?", voucherCode).
		Find(&redeems).Error

	if err != nil {
		r.log.Error("Error fetching users by voucher code", zap.Error(err))
	}

	if len(redeems) == 0 {
		return nil, errors.New("no users found for the given voucher code")
	}

	r.log.Info("Query result", zap.Any("redeems", redeems))

	return redeems, err
}

func (r *redeemRepository) FindRedeemHistoryByUser(userID int) ([]models.Redeem, error) {
	var redeems []models.Redeem
	err := r.DB.Where("user_id = ?", userID).Find(&redeems).Error
	if err != nil {
		r.log.Error("Error fetching redeem voucher by users", zap.Error(err))
	}

	if len(redeems) == 0 {
		return nil, errors.New("no voucher exchange history found")
	}
	return redeems, err
}
