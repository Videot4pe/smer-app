package app

import (
	_ "backend/docs"
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/domain/files"
	"backend/internal/domain/smer"
	"backend/internal/domain/user"
	"backend/pkg/logging"
	"backend/pkg/metric"
	"backend/pkg/oauth"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Handler interface {
	Register(router *httprouter.Router)
}

func NewRouter(ctx context.Context, config *config.Config, logger *logging.Logger, pgClient *pgxpool.Pool) *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	logger.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Println("heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	filesStorage := files.NewFilesStorage(ctx, pgClient, logger)

	userStorage := user.NewUserStorage(ctx, pgClient, logger)
	userHandler := user.NewUserHandler(ctx, userStorage, logger, filesStorage)
	userHandler.Register(router)

	authHandler := auth.NewAuthHandler(ctx, userStorage, logger, config)
	authHandler.Register(router)

	oauthProvider := oauth.GetOAuthProvider(logger, config, userStorage)
	oauthProvider.UseVKAuth(router)
	oauthProvider.UseGoogleAuth(router)

	smerStorage := smer.NewSmerStorage(ctx, pgClient, logger)
	smerHandler := smer.NewSmerHandler(ctx, smerStorage, logger)
	smerHandler.Register(router)

	return router
}
