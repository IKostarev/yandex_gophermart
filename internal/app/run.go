package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"yandex_gophermart/internal/config"
	"yandex_gophermart/internal/loaylty_system"
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
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cfg := config.NewConfig()

	db, err := sql.Open("pgx", cfg.Database.ConnectionString)
	if err != nil {
		log.Sugar().Errorf("error while init db: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Sugar().Errorf("error while closing db: %s", err.Error())
			os.Exit(1)
		}
	}()

	dbManager, err := storage.NewPostgres(ctx, db)
	if err != nil {
		log.Sugar().Errorf("error while init db: %s", err.Error())
		os.Exit(1)
	}

	appServer := server.NewServer(cfg.Server.Address, NewApp(dbManager, log.Sugar()))

	loyaltyPointsSystem := loaylty_system.NewLSystem(cfg.AccrualSystem.Address, dbManager, log.Sugar())

	run := runner.NewRun(appServer, loyaltyPointsSystem, log.Sugar())
	if err = run.Run(ctx); err != nil {
		log.Sugar().Errorf("error while running runner: %s", err.Error())
		return
	}
}
