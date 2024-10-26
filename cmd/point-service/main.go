package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"homework/internal/app"
	"homework/internal/app/mw"
	"homework/internal/app/point_service"
	"homework/internal/cache"
	"homework/internal/config"
	"homework/internal/storage/postgres"
	desc "homework/pkg/point-service/v1"

	"github.com/go-chi/chi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.InitMetrics()
	go app.StartMetricsEndpoint()
	app.IncrementOrdersGiven()

	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.Connect(ctx, cfg.PsqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	redisCache := cache.NewRedisCache("localhost:6379", "qwerty", 0, 10*time.Minute)
	storageFacade := newStorageFacade(pool)

	pointService := point_service.NewImplementation(storageFacade, redisCache)

	lis, err := net.Listen("tcp", cfg.GrpcHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.Logging, mw.Auth),
	)
	reflection.Register(grpcServer)

	desc.RegisterPointServiceServer(grpcServer, pointService)

	mux := runtime.NewServeMux()
	err = desc.RegisterPointServiceHandlerFromEndpoint(ctx, mux, cfg.GrpcHost, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatalf("failed to register point service handler: %v", err)
	}

	go func() {
		if err = http.ListenAndServe(cfg.HttpHost, mux); err != nil {
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
		if err = http.ListenAndServe(cfg.AdminHost, adminServer); err != nil {
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
