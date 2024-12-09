package service

import (
	"voucher_system/models"
	"voucher_system/repository"

	"go.uber.org/zap"
)

type HistoryService interface {
	GetRedeemHistoryByUser(userID int) ([]models.Redeem, error)
	GetUsageHistoryByUser(userID int) ([]models.History, error)
	GetUsersByVoucherCode(voucherCode string) ([]models.Redeem, error)
}

type historyService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewHistoryService(repo repository.Repository, log *zap.Logger) HistoryService {
	return &historyService{
		repo: repo,
		log:  log,
	}
}

func (s *historyService) GetRedeemHistoryByUser(userID int) ([]models.Redeem, error) {
    return s.repo.Redeem.FindRedeemHistoryByUser(userID)
}

func (s *historyService) GetUsageHistoryByUser(userID int) ([]models.History, error) {
	return s.repo.History.FindUsageHistoryByUser(userID)
}

func (s *historyService) GetUsersByVoucherCode(voucherCode string) ([]models.Redeem, error) {
	return s.repo.Redeem.FindUsersByVoucherCode(voucherCode)
}
