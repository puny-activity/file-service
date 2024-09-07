package httprouter

import (
	"github.com/go-chi/chi/v5"
	httpcontroller "github.com/puny-activity/file-service/api/http/controller"
	httpmiddleware "github.com/puny-activity/file-service/api/http/middleware"
	"github.com/puny-activity/file-service/config"
	"github.com/rs/zerolog"
)

type Router struct {
	cfg        *config.HTTP
	router     *chi.Mux
	middleware *httpmiddleware.Middleware
	wrapper    *Wrapper
	controller *httpcontroller.Controller
	log        *zerolog.Logger
}

func New(cfg *config.HTTP, router *chi.Mux, middleware *httpmiddleware.Middleware,
	wrapper *Wrapper, controller *httpcontroller.Controller, log *zerolog.Logger) *Router {
	return &Router{
		cfg:        cfg,
		router:     router,
		middleware: middleware,
		wrapper:    wrapper,
		controller: controller,
		log:        log,
	}
}

func (r *Router) Setup() {
	r.router.Group(func(router chi.Router) {
		router.Route("/stream", func(router chi.Router) {
			router.Get("/{file_id}", r.wrapper.Wrap(r.controller.StreamFile))
		})
	})
}
