package main

import (
	"context"
	"flag"
	"log"

	"homework/internal/app"
	"homework/internal/config"
	"homework/internal/storage/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dataFlag = flag.String("data", "{}", "data in JSON format")
)

func main() {
	// const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.Connect(ctx, cfg.PsqlDSN)
	// pool, err := pgxpool.Connect(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// err = app.ClearTables(pool)
	// err = app.GenerateFakeOrders(pool, 100)
	// if err != nil {
	// 	log.Fatalf("Error generating fake orders: %v", err)
	// }

	//fmt.Println("Successfully generated test data!")
	storageFacade := newStorageFacade(pool)
	service := app.NewService(storageFacade)

	if err := app.RunCLI(ctx, service, dataFlag); err != nil {
		log.Fatal(err)
	}
}

func newStorageFacade(pool *pgxpool.Pool) app.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgRepository(txManager)

	return app.NewStorageFacade(txManager, pgRepository)
}
