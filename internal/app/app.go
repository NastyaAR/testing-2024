package app

import (
	"avito-test-task/config"
	"avito-test-task/internal/delivery/handlers"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"time"
)

func Run(cfg *config.Config) {
	lg, err := pkg.CreateLogger(cfg.LogFile, "prod")
	if err != nil {
		log.Fatal("can't create logger")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.Db.Db)
	pool, err := pgx.Connect(ctx, connString)
	defer pool.Close(context.Background())
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	houseRepo := repo.NewPostgresHouseRepo(pool)
	houseUsecase := usecase.NewHouseUsecase(houseRepo, 5*time.Second)
	houseHandler := handlers.NewHouseHandler(houseUsecase, time.Duration(cfg.DbTimeoutSec), lg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/house/create", houseHandler.Create)
	r.Get("/house/", houseHandler.GetFlatsByID)

	err = http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("final")
}
