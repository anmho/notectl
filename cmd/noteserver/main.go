package main

import (
	"database/sql"
	"fmt"
	"github.com/anmho/notectl/notes"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"log"
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

func main() {

	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	log.Println(config)
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

	s := grpc.NewServer()
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
