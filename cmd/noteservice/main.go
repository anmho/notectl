package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/anmho/notectl/notes"
	"github.com/caarlos0/env/v6"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)
import pb "github.com/anmho/notectl/gen/proto/notes"

type Config struct {
	DbHost string `env:"DB_HOST"`
	DbUser string `env:"DB_USER"`
	DbPass string `env:"DB_PASS"`
	DbName string `env:"DB_NAME"`
	DbPort int    `env:"DB_PORT"`
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields)
	})
}

func main() {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	logger := slog.Default()

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}
	logger.Info("config: ", config)
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPass,
		config.DbName)

	db, err := sql.Open("pgx", connString)

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...)),
	)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(50)

	service := notes.NewService(db)
	pb.RegisterNoteServiceServer(s, &service)
	log.Printf("server listening at %s\n", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}

}
