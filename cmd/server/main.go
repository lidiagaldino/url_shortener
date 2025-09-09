package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/infra"
	"url-shortener/internal/services"
	"url-shortener/pkg"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(cfg.DBName)

	repo := infra.NewMongoURLRepository(db)
	idGen := &pkg.ShortIDGenerator{}
	service := services.NewURLService(repo, idGen)
	handler := handlers.NewURLHandler(service)

	r := chi.NewRouter()
	r.Post("/shorten", handler.Shorten)
	r.Get("/{id}", handler.Redirect)

	log.Printf("🚀 Servidor rodando em %s\n", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
		log.Fatal(err)
	}
}
