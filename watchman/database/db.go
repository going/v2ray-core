package database

import (
	"context"

	"v2ray.com/core/watchman/database/connector"
	"v2ray.com/core/watchman/proto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Connect(ctx context.Context, logger *zap.Logger, cfg *proto.DBConfig) {
	conn, err := sqlx.ConnectContext(ctx, "mysql", cfg.Master)
	if err != nil {
		logger.Panic(err.Error())
	}
	conn = conn.Unsafe()

	conn.SetMaxOpenConns(cfg.MaxOpen)
	conn.SetMaxIdleConns(cfg.MaxIdle)

	if err := conn.PingContext(ctx); err != nil {
		logger.Panic(err.Error())
	}

	connector.RegisterDB(connector.Connectors.Default, conn)
}
