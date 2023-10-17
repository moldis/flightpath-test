package main

import (
	mw "artemb/flights-path/pkg/api/middleware"
	"artemb/flights-path/pkg/api/routes"
	"artemb/flights-path/pkg/config"
	"artemb/flights-path/pkg/logging"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

const (
	cmdServer = "server"
)

func main() {
	app := kingpin.New("api", "Flights path API")
	configFile := app.
		Flag("config", "path to config file").
		Short('c').
		Required().
		PlaceHolder("./path/config.yaml").
		String()

	app.Command(cmdServer, "runs the API server")
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	cfg := config.Read(*configFile)

	logger, undo := initLogger(cfg)
	defer undo()

	router, err := initRouter(cfg, logger)
	if err != nil {
		log.Fatalln(err)
	}

	switch command {
	case cmdServer:
		err = runServer(router, cfg)
		if err != nil {
			logger.Fatal("Server fatal error", zap.Error(err))
		}
	}
}

func initRouter(cfg *config.Config, logger *zap.Logger) (*chi.Mux, error) {
	// TODO consider to replace with https://github.com/gin-gonic/gin
	// robust framework with build-in validator
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger(logger))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.SetHeader("Content-type", "application/json"))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Api.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Api.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Api.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Api.Cors.ExposedHeaders,
		AllowCredentials: cfg.Api.Cors.AllowCredentials,
		MaxAge:           cfg.Api.Cors.MaxAge,
	}))

	compressor := middleware.NewCompressor(5)
	r.Use(compressor.Handler)

	if err := routes.MakeRoutes(r, cfg, logger); err != nil {
		return nil, err
	}

	return r, nil
}

func initLogger(cfg *config.Config) (*zap.Logger, func()) {
	logger, err := logging.NewLogger(cfg.Logging)
	if err != nil {
		panic(fmt.Sprintf("Can't initialize logger: %s", err.Error()))
	}

	logger = logger.Named(cfg.AppName)
	undo := zap.ReplaceGlobals(logger)

	return logger, func() {
		undo()
		_ = logger.Sync()
	}
}

func runServer(r *chi.Mux, cfg *config.Config) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Api.Port), r)
}
