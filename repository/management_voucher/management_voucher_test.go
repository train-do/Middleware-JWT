package managementvoucher_test

import (
	"fmt"
	"testing"
	"time"
	"voucher_system/models"
	managementvoucher "voucher_system/repository/management_voucher"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	return gormDB, mock
}

func TestCreateVoucher(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()

	customerRepo := managementvoucher.NewManagementVoucherRepo(db, &log)

	t.Run("Succesfully create a voucher", func(t *testing.T) {

		voucher := &models.Voucher{
			VoucherName:     "Promo December",
			VoucherCode:     "DESC2024",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "Get more discount on december",
			VoucherCategory: "discount",
			DiscountValue:   10,
			MinimumPurchase: 200000,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now().AddDate(0, 0, -5),
			EndDate:         time.Now().AddDate(0, 0, -1),
			ApplicableAreas: []string{"US", "Canada"},
			Quota:           100,
			Status:          false,
			CreatedAt:       time.Now().AddDate(0, 0, 1),
			UpdatedAt:       time.Now().AddDate(0, 0, -1),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "vouchers"`).
			WithArgs(
				voucher.VoucherName,
				voucher.VoucherCode,
				voucher.VoucherType,
				voucher.PointsRequired,
				voucher.Description,
				voucher.VoucherCategory,
				voucher.DiscountValue,
				voucher.MinimumPurchase,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				voucher.Quota,
				voucher.Status,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()
		err := customerRepo.CreateVoucher(voucher)
		assert.NoError(t, err)
		assert.Equal(t, 1, voucher.ID)
		assert.NotEmpty(t, voucher.VoucherName)
	})

	t.Run("Failed to create a voucher", func(t *testing.T) {
		voucher := &models.Voucher{
			VoucherName:     "Promo December",
			VoucherCode:     "DESC2024",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "Get more discount on december",
			VoucherCategory: "discount",
			DiscountValue:   10,
			MinimumPurchase: 200000,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now().AddDate(0, 0, -5),
			EndDate:         time.Now().AddDate(0, 0, -1),
			ApplicableAreas: []string{"US", "Canada"},
			Quota:           100,
			Status:          false,
			CreatedAt:       time.Now().AddDate(0, 0, 1),
			UpdatedAt:       time.Now().AddDate(0, 0, -1),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "vouchers"`).
			WithArgs(
				voucher.VoucherName,
				voucher.VoucherCode,
				voucher.VoucherType,
				voucher.PointsRequired,
				voucher.Description,
				voucher.VoucherCategory,
				voucher.DiscountValue,
				voucher.MinimumPurchase,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				voucher.Quota,
				voucher.Status,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("database error"))

		mock.ExpectRollback()
		err := customerRepo.CreateVoucher(voucher)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestSoftDeleteVoucher(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()

	voucherRepo := managementvoucher.NewManagementVoucherRepo(db, &log)

	t.Run("Successfully soft delete a voucher", func(t *testing.T) {
		voucherID := 1

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "vouchers" SET "deleted_at"=`).
			WithArgs(sqlmock.AnyArg(), voucherID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := voucherRepo.SoftDeleteVoucher(voucherID)
		assert.NoError(t, err)
	})

	t.Run("Failed to soft delete a voucher", func(t *testing.T) {
		voucherID := 2

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "vouchers" SET "deleted_at"=`).
			WithArgs(sqlmock.AnyArg(), voucherID).
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		err := voucherRepo.SoftDeleteVoucher(voucherID)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})

}

func TestUpdateVoucher(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	voucherRepo := managementvoucher.NewManagementVoucherRepo(db, &log)

	t.Run("Successfully update a voucher", func(t *testing.T) {
		voucherID := 1
		voucher := &models.Voucher{
			VoucherName:     "Promo Updated",
			VoucherCode:     "UPDATED2024",
			VoucherType:     "e-commerce",
			PointsRequired:  10,
			Description:     "Updated discount",
			VoucherCategory: "discount",
			DiscountValue:   15,
			MinimumPurchase: 250000,
			Quota:           50,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "vouchers" SET`).
			WithArgs(
				voucher.VoucherName,
				voucher.VoucherCode,
				voucher.VoucherType,
				voucher.PointsRequired,
				voucher.Description,
				voucher.VoucherCategory,
				voucher.DiscountValue,
				voucher.MinimumPurchase,
				voucher.Quota,
				sqlmock.AnyArg(), // updated_at
				voucherID,
			).
			WillReturnResult(sqlmock.NewResult(1, int64(voucherID)))
		mock.ExpectCommit()

		err := voucherRepo.UpdateVoucher(voucher, voucherID)
		assert.NoError(t, err)
	})

	t.Run("Failed to update due to no matching record", func(t *testing.T) {
		voucherID := 2
		voucher := &models.Voucher{
			VoucherName: "Promo Not Found",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "vouchers" SET`).
			WithArgs(voucher.VoucherName, sqlmock.AnyArg(), voucherID).
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
		mock.ExpectCommit()

		err := voucherRepo.UpdateVoucher(voucher, voucherID)
		assert.Error(t, err)
		assert.EqualError(t, err, "no record found with shipping_id 2")
	})

	t.Run("Failed to update due to database error", func(t *testing.T) {
		voucherID := 3
		voucher := &models.Voucher{
			VoucherName: "Promo Error",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "vouchers" SET`).
			WithArgs(voucher.VoucherName, sqlmock.AnyArg(), voucherID).
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		err := voucherRepo.UpdateVoucher(voucher, voucherID)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestShowRedeemPoints(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	voucherRepo := managementvoucher.NewManagementVoucherRepo(db, &log)

	t.Run("Successfully show redeem points", func(t *testing.T) {

		mockRows := sqlmock.NewRows([]string{"voucher_name", "discount_value", "points_required"}).
			AddRow("Promo A", 20.0, 50).
			AddRow("Promo B", 15.0, 30)

		mock.ExpectQuery(`SELECT v.voucher_name, v.discount_value, v.points_required FROM vouchers as v WHERE`).
			WithArgs("redeem points").
			WillReturnRows(mockRows)

		result, err := voucherRepo.ShowRedeemPoints()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)

		assert.Equal(t, "Promo A", (*result)[0].VoucherName)
		assert.Equal(t, 20.0, (*result)[0].DiscountValue)
		assert.Equal(t, 50, (*result)[0].PointsRequired)

		assert.Equal(t, "Promo B", (*result)[1].VoucherName)
		assert.Equal(t, 15.0, (*result)[1].DiscountValue)
		assert.Equal(t, 30, (*result)[1].PointsRequired)
	})

	t.Run("Failed to show redeem points due to database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT v.voucher_name, v.discount_value, v.points_required FROM vouchers as v WHERE`).
			WithArgs("redeem points").
			WillReturnError(fmt.Errorf("database error"))

		result, err := voucherRepo.ShowRedeemPoints()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "database error")
	})
}

func TestGetVouchersByQueryParams(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	voucherRepo := managementvoucher.NewManagementVoucherRepo(db, &log)

	t.Run("Successfully fetch vouchers with all query parameters", func(t *testing.T) {
		rawPaymentMethods := `["Credit Card","PayPal"]`
		rawApplicableAreas := `["US","Canada"]`

		mockRows := sqlmock.NewRows([]string{
			"id", "voucher_name", "voucher_code", "voucher_type", "points_required",
			"description", "voucher_category", "discount_value", "minimum_purchase",
			"payment_methods", "start_date", "end_date", "applicable_areas", "quota",
			"status", "created_at", "updated_at",
		}).
			AddRow(1, "Promo A", "CODE123", "e-commerce", 0, "Description A", "discount", 10.0, 200000,
				rawPaymentMethods, time.Now(), time.Now().AddDate(0, 0, 1), rawApplicableAreas, 100, true, time.Now(), time.Now()).
			AddRow(2, "Promo B", "CODE456", "retail", 0, "Description B", "voucher", 15.0, 150000,
				rawPaymentMethods, time.Now(), time.Now().AddDate(0, 0, 1), rawApplicableAreas, 200, true, time.Now(), time.Now())

		mock.ExpectQuery(`SELECT (.+) FROM "vouchers"`).
			WillReturnRows(mockRows)

		status := "active"
		area := "US"
		voucherType := "e-commerce"

		result, err := voucherRepo.GetVouchersByQueryParams(status, area, voucherType)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)

		// Validate first voucher
		assert.Equal(t, "Promo A", (*result)[0].VoucherName)
		assert.Equal(t, "CODE123", (*result)[0].VoucherCode)
		assert.Equal(t, []string{"Credit Card", "PayPal"}, (*result)[0].PaymentMethods)
		assert.Equal(t, []string{"US", "Canada"}, (*result)[0].ApplicableAreas)

		// Validate second voucher
		assert.Equal(t, "Promo B", (*result)[1].VoucherName)
		assert.Equal(t, "CODE456", (*result)[1].VoucherCode)
	})

	t.Run("No vouchers match query parameters", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM "vouchers"`).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "voucher_name", "voucher_code", "voucher_type", "points_required",
				"description", "voucher_category", "discount_value", "minimum_purchase",
				"payment_methods", "start_date", "end_date", "applicable_areas", "quota",
				"status", "created_at", "updated_at",
			}))

		status := "non-active"
		area := "Europe"
		voucherType := "retail"

		result, err := voucherRepo.GetVouchersByQueryParams(status, area, voucherType)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 0)
	})

	t.Run("Database error while fetching vouchers", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM "vouchers"`).
			WillReturnError(fmt.Errorf("database error"))

		status := "active"
		area := "Asia"
		voucherType := "redeem points"

		result, err := voucherRepo.GetVouchersByQueryParams(status, area, voucherType)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "database error")
	})
}

func TestCreateRedeemVoucher(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	voucherRepo := managementvoucher.NewManagementVoucherRepo(db, &log)

	today := time.Now()

	t.Run("Successfully create redeem voucher", func(t *testing.T) {
		redeem := &models.Redeem{
			UserID:    1,
			VoucherID: 100,
		}
		points := 50

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "redeems" WHERE user_id = \$1 AND voucher_id = \$2 ORDER BY "redeems"."id" LIMIT \$3`).
			WithArgs(redeem.UserID, redeem.VoucherID, 1).
			WillReturnRows(sqlmock.NewRows(nil)) // Tidak ada baris ditemukan

		// Mock fetch voucher data
		mock.ExpectQuery(`SELECT quota, points_required, start_date, end_date FROM "vouchers" WHERE id = \$1`).
			WithArgs(redeem.VoucherID).
			WillReturnRows(sqlmock.NewRows([]string{"quota", "points_required", "start_date", "end_date"}).
				AddRow(10, 50, today.AddDate(0, 0, -5), today.AddDate(0, 0, 5)))

		// Mock create redeem
		mock.ExpectQuery(`INSERT INTO "redeems"`).
			WithArgs(
				redeem.UserID,
				redeem.VoucherID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			// WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock decrement quota
		mock.ExpectExec(`UPDATE "vouchers" SET "quota"=quota - \$1 WHERE id = \$2 AND "vouchers"."deleted_at" IS NULL`).
			WithArgs(1, redeem.VoucherID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := voucherRepo.CreateRedeemVoucher(redeem, points)
		assert.NoError(t, err)
	})

	t.Run("Redeem already exists", func(t *testing.T) {
		redeem := &models.Redeem{
			UserID:    2,
			VoucherID: 101,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "redeems" WHERE user_id = \$1 AND voucher_id = \$2 ORDER BY "redeems"."id" LIMIT \$3`).
			WithArgs(redeem.UserID, redeem.VoucherID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "voucher_id"}).
				AddRow(redeem.UserID, redeem.VoucherID))

		mock.ExpectRollback()

		err := voucherRepo.CreateRedeemVoucher(redeem, 50)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("user_id %d already claimed voucher_id %d", redeem.UserID, redeem.VoucherID))
	})

	t.Run("Insufficient quota", func(t *testing.T) {
		redeem := &models.Redeem{
			UserID:    3,
			VoucherID: 102,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "redeems" WHERE user_id = \$1 AND voucher_id = \$2 ORDER BY "redeems"."id" LIMIT \$3`).
			WithArgs(redeem.UserID, redeem.VoucherID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(`SELECT quota, points_required, start_date, end_date FROM "vouchers"`).
			WithArgs(redeem.VoucherID).
			WillReturnRows(sqlmock.NewRows([]string{"quota", "points_required", "start_date", "end_date"}).
				AddRow(0, 50, today.AddDate(0, 0, -5), today.AddDate(0, 0, 5)))

		mock.ExpectRollback()

		err := voucherRepo.CreateRedeemVoucher(redeem, 50)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("quota for voucher ID %d is not sufficient", redeem.VoucherID))
	})

	t.Run("Points mismatch", func(t *testing.T) {
		redeem := &models.Redeem{
			UserID:    4,
			VoucherID: 103,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "redeems" WHERE user_id = \$1 AND voucher_id = \$2 ORDER BY "redeems"."id" LIMIT \$3`).
			WithArgs(redeem.UserID, redeem.VoucherID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(`SELECT quota, points_required, start_date, end_date FROM "vouchers"`).
			WithArgs(redeem.VoucherID).
			WillReturnRows(sqlmock.NewRows([]string{"quota", "points_required", "start_date", "end_date"}).
				AddRow(10, 100, today.AddDate(0, 0, -5), today.AddDate(0, 0, 5)))

		mock.ExpectRollback()

		err := voucherRepo.CreateRedeemVoucher(redeem, 50)
		assert.Error(t, err)
		assert.EqualError(t, err, "required points (100) do not match provided points (50)")
	})

	t.Run("Voucher expired", func(t *testing.T) {
		redeem := &models.Redeem{
			UserID:    5,
			VoucherID: 104,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "redeems" WHERE user_id = \$1 AND voucher_id = \$2 ORDER BY "redeems"."id" LIMIT \$3`).
			WithArgs(redeem.UserID, redeem.VoucherID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(`SELECT quota, points_required, start_date, end_date FROM "vouchers"`).
			WithArgs(redeem.VoucherID).
			WillReturnRows(sqlmock.NewRows([]string{"quota", "points_required", "start_date", "end_date"}).
				AddRow(10, 50, today.AddDate(0, 0, -10), today.AddDate(0, 0, -1)))

		mock.ExpectRollback()

		err := voucherRepo.CreateRedeemVoucher(redeem, 50)
		assert.Error(t, err)
		assert.EqualError(t, err, "voucher expired")
	})
}
