package app

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"yandex_gophermart/internal/config"
	"yandex_gophermart/internal/loayltysystem"
	"yandex_gophermart/internal/runner"
	"yandex_gophermart/internal/server"
	storage "yandex_gophermart/internal/storage/db"
	"yandex_gophermart/pkg/logger"
)

const logLevel = "info"

func Run() {
	ctx := context.Background()

	log, err := logger.New(logLevel)
	if err != nil {
		log.Sugar().Fatalf("error init log level: %s", err.Error())
	}

	cfg := config.NewConfig()

	//db, err := sql.Open("pgx", cfg.Database.ConnectionString)
	//if err != nil {
	//	log.Sugar().Fatalf("error while init db: %s", err.Error())
	//}
	//defer func() {
	//	if err := db.Close(); err != nil {
	//		log.Sugar().Fatalf("error while closing db: %s", err.Error())
	//	}
	//}()

	dbManager, err := storage.NewPostgres(ctx, cfg, log.Sugar())
	if err != nil {
		log.Sugar().Fatalf("error while init db: %s", err.Error())
	}
	//defer func() {
	//	dbManager.Close()
	//}()

	appServer := server.NewServer(cfg.Server.Address, NewApp(dbManager, log.Sugar()))

	loyaltyPointsSystem := loayltysystem.NewLSystem(cfg.AccrualSystem.Address, dbManager, log.Sugar())

	run := runner.NewRun(appServer, loyaltyPointsSystem, log.Sugar())
	if err = run.Run(ctx); err != nil {
		log.Sugar().Fatalf("error while running runner: %s", err.Error())
		return
	}
}
