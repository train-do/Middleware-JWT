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

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	return gormDB, mock
}

func TestFindUsageHistoryByUser(t *testing.T) {
	db, mock := setupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()
	logger := zap.NewNop()
	repo := repository.NewHistoryRepository(db, logger)

	userID := 123
	rows := sqlmock.NewRows([]string{"id", "user_id", "voucher_id", "transaction_amount", "benefit_value", "usage_date"}).
		AddRow(1, userID, 456, 100.0, 10.0, "2024-12-01").
		AddRow(2, userID, 457, 200.0, 20.0, "2024-12-02")

	mock.ExpectQuery(`SELECT \* FROM "histories" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	histories, err := repo.FindUsageHistoryByUser(userID)

	assert.NoError(t, err)
	assert.Len(t, histories, 2)
}


// func TestCreateHistory(t *testing.T) {
// 	db, mock := setupTestDB()
// 	defer func() { _ = mock.ExpectationsWereMet() }()
// 	logger := zap.NewNop()
// 	repo := repository.NewHistoryRepository(db, logger)

// 	history := &models.History{
// 		ID:                1,
// 		UserID:            123,
// 		VoucherID:         456,
// 		TransactionAmount: 100.0,
// 		BenefitValue:      10.0,
// 		UsageDate:         time.Now(),
// 	}

// 	mock.ExpectExec(`INSERT INTO "histories" \("id","user_id","voucher_id","transaction_amount","benefit_value","usage_date"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
// 		WithArgs(history.ID, history.UserID, history.VoucherID, history.TransactionAmount, history.BenefitValue, sqlmock.AnyArg()).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	err := repo.CreateHistory(history)

// 	assert.NoError(t, err)
// }

