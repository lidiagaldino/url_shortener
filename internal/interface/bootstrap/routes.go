package bootstrap

import (
	"net/http"

	"url-shortener/internal/config"
	"url-shortener/internal/infra/persistence"
	"url-shortener/internal/infra/security"
	"url-shortener/internal/interface/handler"
	"url-shortener/internal/interface/middleware"
	"url-shortener/internal/services"
	"url-shortener/pkg"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database, cfg *config.Config) http.Handler {
	urlRepo := persistence.NewMongoURLRepository(db)
	userRepo := persistence.NewMongoUserRepository(db)
	statsRepo := persistence.NewMongoURLStatsRepository(db)

	idGen := &pkg.ShortIDGenerator{}
	hasher := &pkg.PassowrdHasher{}
	tokenGen := security.NewJWTService(cfg.SecretKey)

	urlService := services.NewURLService(urlRepo, idGen, statsRepo)
	userService := services.NewUserService(userRepo, hasher, tokenGen)

	urlHandler := handler.NewURLHandler(urlService)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Post("/users", userHandler.Save)
	r.Post("/users/signin", userHandler.Login)

	r.Get("/urls/{id}", urlHandler.Redirect)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware(tokenGen))
		protected.Post("/urls/shorten", urlHandler.Shorten)
	})

	return r
}
