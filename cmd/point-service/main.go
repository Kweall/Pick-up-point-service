package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"homework/internal/app"
	"homework/internal/app/mw"
	"homework/internal/app/point_service"
	"homework/internal/storage/postgres"
	desc "homework/pkg/point-service/v1"

	"github.com/go-chi/chi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	psqlDSN   = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	grpcHost  = "127.0.0.1:7001"
	httpHost  = "127.0.0.1:7000"
	adminHost = "127.0.0.1:7002"
)

func main() {

	// cfg, err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// pool, err := pgxpool.Connect(ctx, cfg.PsqlDSN)
	pool, err := pgxpool.Connect(ctx, psqlDSN)
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
	// service := app.NewService(storageFacade)

	// if err := app.RunCLI(ctx, service, dataFlag); err != nil {
	// 	log.Fatal(err)
	// }
	pointService := point_service.NewImplementation(storageFacade)

	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.Logging, mw.Auth),
	)
	reflection.Register(grpcServer)

	desc.RegisterPointServiceServer(grpcServer, pointService)

	mux := runtime.NewServeMux()
	err = desc.RegisterPointServiceHandlerFromEndpoint(ctx, mux, grpcHost, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatalf("failed to register point service handler: %v", err)
	}

	go func() {
		if err = http.ListenAndServe(httpHost, mux); err != nil {
			log.Fatalf("failed to listen and serve point service handler: %v", err)
		}
	}()

	go func() {
		adminServer := chi.NewMux()
		adminServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			b, err := os.ReadFile("./pkg/point-service/v1/point_service.swagger.json")
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				log.Printf("failed to read swagger.json: %v", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		})
		if err = http.ListenAndServe(adminHost, adminServer); err != nil {
			log.Fatalf("failed to listen and serve admin server: %v", err)
		}
	}()

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newStorageFacade(pool *pgxpool.Pool) app.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgRepository(txManager)

	return app.NewStorageFacade(txManager, pgRepository)
}
