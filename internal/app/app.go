package app

import (
	"time"

	"go.uber.org/zap"

	grpcapp "texts/internal/app/grpc"
	"texts/internal/service/texts"
	"texts/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *zap.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	secret string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	textsService := texts.New(log, storage, storage, storage, tokenTTL)

	grpcApp, err := grpcapp.New(log, textsService, grpcPort, secret)
	if err != nil {
		panic(err)
	}

	return &App{
		GRPCServer: grpcApp,
	}
}
