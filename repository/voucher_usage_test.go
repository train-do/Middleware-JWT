package repository_test

import (
	"testing"
	"time"
	"voucher_system/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	return gormDB, mock
}

func TestVoucherRepository_FindAll(t *testing.T) {
	db, mock := SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewVoucherRepository(db, logger)

	rows := sqlmock.NewRows([]string{"id", "voucher_code", "voucher_type", "quota", "start_date", "end_date", "minimum_purchase", "payment_methods", "applicable_areas"}).
		AddRow(1, "VOUCHER1", "e-commerce", 100, time.Now(), time.Now().Add(24*time.Hour), 50.0, `["credit"]`, `["area1"]`)

	mock.ExpectQuery("SELECT vouchers.*, vouchers.payment_methods AS raw_payment_methods, vouchers.applicable_areas AS raw_applicable_areas").
		WillReturnRows(rows)

	vouchers, err := repo.FindAll(1, "e-commerce")

	assert.NoError(t, err)
	assert.Len(t, vouchers, 1)
	assert.Equal(t, "VOUCHER1", vouchers[0].VoucherCode)
}

func TestVoucherRepository_FindValidVoucher_VoucherNotFound(t *testing.T) {
	db, mock := SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewVoucherRepository(db, logger)

	mock.ExpectQuery("SELECT vouchers.*").
		WillReturnError(gorm.ErrRecordNotFound)

	voucher, err := repo.FindValidVoucher(1, "INVALIDCODE", "area1", 100.0, 10.0, "credit", time.Now())

	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Equal(t, "voucher not found", err.Error())
}

func TestVoucherRepository_FindValidVoucher_MinimumPurchaseNotMet(t *testing.T) {
	db, mock := SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewVoucherRepository(db, logger)

	rows := sqlmock.NewRows([]string{"id", "voucher_code", "voucher_type", "quota", "start_date", "end_date", "minimum_purchase", "payment_methods", "applicable_areas"}).
		AddRow(1, "VOUCHER1", "e-commerce", 100, time.Now(), time.Now().Add(24*time.Hour), 50.0, `["credit"]`, `["area1"]`)

	mock.ExpectQuery("SELECT vouchers.*").
		WillReturnRows(rows)

	voucher, err := repo.FindValidVoucher(1, "VOUCHER1", "area1", 40.0, 10.0, "credit", time.Now())

	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Equal(t, "transaction amount must be at least 50.00", err.Error())
}

func TestVoucherRepository_FindValidVoucher_AreaNotApplicable(t *testing.T) {
	db, mock := SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewVoucherRepository(db, logger)

	rows := sqlmock.NewRows([]string{"id", "voucher_code", "voucher_type", "quota", "start_date", "end_date", "minimum_purchase", "payment_methods", "applicable_areas"}).
		AddRow(1, "VOUCHER1", "e-commerce", 100, time.Now(), time.Now().Add(24*time.Hour), 50.0, `["credit"]`, `["area2"]`)

	mock.ExpectQuery("SELECT vouchers.*").
		WillReturnRows(rows)

	voucher, err := repo.FindValidVoucher(1, "VOUCHER1", "area1", 100.0, 10.0, "credit", time.Now())

	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Equal(t, "area not found", err.Error())
}

func TestVoucherRepository_FindValidVoucher_PaymentMethodNotApplicable(t *testing.T) {
	db, mock := SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewVoucherRepository(db, logger)

	rows := sqlmock.NewRows([]string{"id", "voucher_code", "voucher_type", "quota", "start_date", "end_date", "minimum_purchase", "payment_methods", "applicable_areas"}).
		AddRow(1, "VOUCHER1", "e-commerce", 100, time.Now(), time.Now().Add(24*time.Hour), 50.0, `["credit"]`, `["area1"]`)

	mock.ExpectQuery("SELECT vouchers.*").
		WillReturnRows(rows)

	voucher, err := repo.FindValidVoucher(1, "VOUCHER1", "area1", 100.0, 10.0, "debit", time.Now())

	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Equal(t, "payment method not found", err.Error())
}

func TestVoucherRepository_FindValidVoucher_VoucherExpired(t *testing.T) {
	db, mock := SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewVoucherRepository(db, logger)

	rows := sqlmock.NewRows([]string{"id", "voucher_code", "voucher_type", "quota", "start_date", "end_date", "minimum_purchase", "payment_methods", "applicable_areas"}).
		AddRow(1, "VOUCHER1", "e-commerce", 100, time.Now().Add(-48*time.Hour), time.Now().Add(-24*time.Hour), 50.0, `["credit"]`, `["area1"]`)

	mock.ExpectQuery("SELECT vouchers.*").
		WillReturnRows(rows)

	voucher, err := repo.FindValidVoucher(1, "VOUCHER1", "area1", 100.0, 10.0, "credit", time.Now())

	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Equal(t, "voucher expired", err.Error())
}

// func TestVoucherRepository_UpdateVoucherQuota(t *testing.T) {
//     db, mock := SetupTestDB()
//     defer func() { _ = mock.ExpectationsWereMet() }()

//     logger := zap.NewNop()
//     repo := repository.NewVoucherRepository(db, logger)

//     voucherID := 1
//     newQuota := 50

//     mock.ExpectBegin() // Expect transaction to start
//     mock.ExpectExec("UPDATE \"vouchers\" SET \"quota\"=$1, \"updated_at\"=$2 WHERE \"id\" = $3 AND \"vouchers\".\"deleted_at\" IS NULL").
//         WithArgs(newQuota, sqlmock.AnyArg(), voucherID).
//         WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful update
//     mock.ExpectCommit() // Expect transaction commit

//     err := repo.UpdateVoucherQuota(voucherID, newQuota)

//     assert.NoError(t, err)
// }

// func TestVoucherRepository_UpdateVoucherQuota_Failed(t *testing.T) {
//     db, mock := SetupTestDB()
//     defer func() { _ = mock.ExpectationsWereMet() }()

//     logger := zap.NewNop()
//     repo := repository.NewVoucherRepository(db, logger)

//     voucherID := 1
//     newQuota := 50

//     mock.ExpectBegin() // Expect transaction to start
//     mock.ExpectExec("UPDATE \"vouchers\" SET \"quota\"=$1, \"updated_at\"=$2 WHERE \"id\" = $3 AND \"vouchers\".\"deleted_at\" IS NULL").
//         WithArgs(newQuota, sqlmock.AnyArg(), voucherID).
//         WillReturnError(fmt.Errorf("update failed")) // Simulate failure
//     mock.ExpectRollback() // Expect rollback after failure

//     err := repo.UpdateVoucherQuota(voucherID, newQuota)

//     assert.Error(t, err)
//     assert.Equal(t, "update failed", err.Error()) // Ensure the error is returned
// }
