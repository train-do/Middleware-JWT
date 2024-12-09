package repository_test

import (
	"testing"
	"voucher_system/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTest() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	return gormDB, mock
}

// func TestRedeemRepository_FindUsersByVoucherCode(t *testing.T) {
// 	db, mock := SetupTest()
// 	defer func() { _ = mock.ExpectationsWereMet() }()

// 	logger := zap.NewNop()
// 	repo := repository.NewRedeemRepository(db, logger)

// 	voucherCode := "ABC123"
// 	userID := 1
// 	redeems := []models.Redeem{
// 		{ID: 1, UserID: userID, VoucherID: 1, RedeemDate: time.Now()},
// 		{ID: 2, UserID: userID, VoucherID: 1, RedeemDate: time.Now()},
// 	}

// 	mock.ExpectQuery(`SELECT redeems.* FROM "redeems" JOIN "vouchers" ON "vouchers".id = "redeems".voucher_id WHERE "vouchers".voucher_code = $1`).
// 		WithArgs(voucherCode).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "voucher_id", "redeem_date"}).
// 			AddRow(1, 1, 1, time.Now()).
// 			AddRow(2, 1, 1, time.Now()))

// 	result, err := repo.FindUsersByVoucherCode(voucherCode)

// 	fmt.Println("Result:", result)
// 	fmt.Println("Expected redeems:", redeems)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, len(result), len(redeems))
// 	if len(result) > 0 {
// 		assert.Equal(t, result[0].UserID, redeems[0].UserID)
// 	}
// }

func TestRedeemRepository_FindUsersByVoucherCode_NoResults(t *testing.T) {
	db, mock := SetupTest()
	defer func() { _ = mock.ExpectationsWereMet() }()

	logger := zap.NewNop()
	repo := repository.NewRedeemRepository(db, logger)

	voucherCode := "XYZ789"

	mock.ExpectQuery(`SELECT redeems.* FROM redeems JOIN vouchers ON vouchers.id = redeems.voucher_id WHERE vouchers.voucher_code = $1`).
		WithArgs(voucherCode).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "voucher_id"}))

	result, err := repo.FindUsersByVoucherCode(voucherCode)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "no users found for the given voucher code")
}

func TestFindRedeemHistoryByUser(t *testing.T) {
	db, mock := SetupTest()
	defer func() { _ = mock.ExpectationsWereMet() }()
	logger := zap.NewNop()
	repo := repository.NewRedeemRepository(db, logger)

	userID := 123
	rows := sqlmock.NewRows([]string{"id", "user_id", "voucher_id", "usage_date"}).
		AddRow(1, userID, 456, "2024-12-01").
		AddRow(2, userID, 457, "2024-12-02")

	mock.ExpectQuery(`SELECT \* FROM "redeems" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	redeems, err := repo.FindRedeemHistoryByUser(userID)

	assert.NoError(t, err)
	assert.Len(t, redeems, 2)
}
