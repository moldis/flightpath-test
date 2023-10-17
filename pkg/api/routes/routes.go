package routes

import (
	"artemb/flights-path/pkg/api/controller"
	"artemb/flights-path/pkg/config"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	baseRoute = "/"
	calculate = "/calculate"
)

type dependencies struct {
	logger *zap.Logger
}

func MakeRoutes(router chi.Router, cfg *config.Config, logger *zap.Logger) error {
	return makeV1Routes(router, cfg, logger)
}

func makeV1Routes(router chi.Router, cfg *config.Config, logger *zap.Logger) error {
	deps, err := makeDeps(cfg, logger)
	if err != nil {
		return err
	}

	searchController := makeSearchController(deps)
	router.
		Route(baseRoute, func(r chi.Router) {
			r.Route(calculate, makeSearchRoutes(searchController))
		})

	return nil
}

func makeSearchRoutes(ctrl *controller.SearchController) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get(baseRoute, ctrl.Search)
	}
}

func makeSearchController(deps *dependencies) *controller.SearchController {
	return &controller.SearchController{
		Logger: deps.logger,
	}
}

func makeDeps(cfg *config.Config, logger *zap.Logger) (*dependencies, error) {
	return &dependencies{logger: logger}, nil
}
