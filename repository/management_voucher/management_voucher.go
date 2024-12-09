package managementvoucher

import (
	"encoding/json"
	"fmt"
	"time"
	"voucher_system/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ManagementVoucherInterface interface {
	CreateVoucher(voucher *models.Voucher) error
	SoftDeleteVoucher(voucherID int) error
	UpdateVoucher(voucher *models.Voucher, voucherID int) error
	ShowRedeemPoints() (*[]RedeemPoint, error)
	GetVouchersByQueryParams(status, area, voucher_type string) (*[]models.Voucher, error)
	CreateRedeemVoucher(redeem *models.Redeem, points int) error
}

type ManagementVoucherRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewManagementVoucherRepo(db *gorm.DB, log *zap.Logger) ManagementVoucherInterface {
	return &ManagementVoucherRepo{DB: db, Log: log}
}

func (m *ManagementVoucherRepo) CreateVoucher(voucher *models.Voucher) error {
	err := m.DB.Create(voucher).Error
	if err != nil {
		m.Log.Error("Error from repo creating voucher:", zap.Error(err))
		return err
	}

	return nil
}

func (m *ManagementVoucherRepo) SoftDeleteVoucher(voucherID int) error {

	err := m.DB.Delete(&models.Voucher{}, voucherID).Error
	if err != nil {
		m.Log.Error("Error from repo soft deleting voucher:", zap.Error(err))
		return err
	}

	return nil
}

func (m *ManagementVoucherRepo) UpdateVoucher(voucher *models.Voucher, voucherID int) error {

	result := m.DB.Model(&voucher).
		Where("id = ?", voucherID).
		Updates(voucher)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found with shipping_id %d", voucherID)
	}

	return nil
}

type RedeemPoint struct {
	VoucherName    string  `json:"voucher_name"`
	PointsRequired int     `json:"points_required"`
	DiscountValue  float64 `json:"discount_value"`
}

func (m *ManagementVoucherRepo) ShowRedeemPoints() (*[]RedeemPoint, error) {

	voucher := []RedeemPoint{}

	query := m.DB.
		Table("vouchers as v").
		Select(`v.voucher_name, v.discount_value, v.points_required`).
		Where("voucher_type = ? AND start_date <= NOW() AND end_date >= NOW()", "redeem points")

	err := query.Find(&voucher).Error
	if err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (m *ManagementVoucherRepo) GetVouchersByQueryParams(status, area, voucher_type string) (*[]models.Voucher, error) {

	var rawVouchers []struct {
		models.Voucher
		RawPaymentMethods  []byte `gorm:"column:payment_methods"`
		RawApplicableAreas []byte `gorm:"column:applicable_areas"`
	}

	query := m.DB.Model(&models.Voucher{})

	if area != "" {
		query = query.Where("applicable_areas @> ?", fmt.Sprintf(`["%s"]`, area))
	}

	if status != "" {
		if status == "active" {
			query = query.Where("start_date <= NOW() AND end_date >= NOW()")
		} else if status == "non-active" {
			query = query.Where("end_date < NOW()")
		}
	}

	if voucher_type != "" {
		query = query.Where("voucher_type = ?", voucher_type)
	}

	err := query.Find(&rawVouchers).Error
	if err != nil {
		return nil, err
	}

	vouchers := make([]models.Voucher, 0, len(rawVouchers))
	for _, rawVoucher := range rawVouchers {
		v := rawVoucher.Voucher

		if len(rawVoucher.RawPaymentMethods) > 0 {
			if err := json.Unmarshal(rawVoucher.RawPaymentMethods, &v.PaymentMethods); err != nil {
				return nil, err
			}
		}

		if len(rawVoucher.RawApplicableAreas) > 0 {
			if err := json.Unmarshal(rawVoucher.RawApplicableAreas, &v.ApplicableAreas); err != nil {
				return nil, err
			}
		}
		vouchers = append(vouchers, v)

	}

	return &vouchers, nil
}

func (m *ManagementVoucherRepo) CreateRedeemVoucher(redeem *models.Redeem, points int) error {

	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			m.Log.Panic("Transaction rolled back due to panic", zap.Any("reason", r))
		}
	}()

	err := tx.Model(&models.Redeem{}).
		Where("user_id = ? AND voucher_id = ?", redeem.UserID, redeem.VoucherID).
		First(&models.Redeem{}).Error

	if err == nil {
		tx.Rollback()
		return fmt.Errorf("user_id %d already claimed voucher_id %d", redeem.UserID, redeem.VoucherID)
	}

	var voucher struct {
		Quota          int
		PointsRequired int
		StartDate      time.Time
		EndDate        time.Time
	}

	today := time.Now()

	err = tx.Model(&models.Voucher{}).
		Where("id = ?", redeem.VoucherID).
		Select("quota, points_required, start_date, end_date").
		Scan(&voucher).Error
	if err != nil {
		tx.Rollback()
		m.Log.Error("Failed to fetch voucher data: ", zap.Error(err))
		return err
	}
	fmt.Println(voucher.PointsRequired)

	if voucher.Quota <= 0 {
		tx.Rollback()
		return fmt.Errorf("quota for voucher ID %d is not sufficient", redeem.VoucherID)
	}

	if points != voucher.PointsRequired {
		tx.Rollback()
		return fmt.Errorf("required points (%d) do not match provided points (%d)", voucher.PointsRequired, points)
	}

	if voucher.StartDate.After(today) {
		tx.Rollback()
		return fmt.Errorf("voucher cannot be used before its start date: %s", voucher.StartDate.Format("2006-01-02"))
	}

	if voucher.EndDate.Before(today) {
		tx.Rollback()
		return fmt.Errorf("voucher expired")
	}

	err = tx.Create(redeem).Error
	if err != nil {
		tx.Rollback()
		m.Log.Error("Failed to create redeem: ", zap.Error(err))
		return err
	}

	err = tx.Model(&models.Voucher{}).
		Where("id = ?", redeem.VoucherID).
		UpdateColumn("quota", gorm.Expr("quota - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		m.Log.Error("Failed to decrement voucher quota: ", zap.Error(err))
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		m.Log.Error("Failed to commit transaction: ", zap.Error(err))
		return err
	}

	return nil
}
