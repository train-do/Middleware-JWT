package infra

import (
	"voucher_system/config"
	"voucher_system/controller"
	"voucher_system/database"
	"voucher_system/helper"
	"voucher_system/middleware"
	"voucher_system/repository"
	"voucher_system/service"
	"voucher_system/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Cfg config.Configuration
	DB  *gorm.DB
	Ctl controller.Controller
	Log *zap.Logger
	Cacher     database.Cacher
	Middleware middleware.Middleware
}

func NewServiceContext() (*ServiceContext, error) {

	handlerError := func(err error) (*ServiceContext, error) {
		return nil, err
	}

	// instance config
	config, err := config.ReadConfig()
	if err != nil {
		handlerError(err)
	}

	if err := utils.InitJwtKey(config); err != nil {
		return handlerError(err)
	}

	// instance looger
	log, err := helper.InitZapLogger()
	if err != nil {
		handlerError(err)
	}

	// instance database
	db, err := database.InitDB(config)
	if err != nil {
		handlerError(err)
	}

	rdb := database.NewCacher(config, 60*60)

	middleware := middleware.NewMiddleware(log, rdb)

	// instance repository
	repository := repository.NewRepository(db, log)

	// instance service
	service := service.NewService(repository, log)

	// instance controller
	Ctl := controller.NewController(service, log, rdb)



	return &ServiceContext{Cfg: config, DB: db, Ctl: *Ctl, Log: log, Cacher: rdb, Middleware: middleware}, nil
}
