package http

import (
	"fmt"
	"net/http"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/docs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

type HTTP struct {
	Config *configs.Config
	DB     *infras.MySQLConn
	Router *chi.Mux
}

// Shutdown shuts down the service.
func (h *HTTP) Shutdown(done chan os.Signal) {
	<-done
	defer os.Exit(0)
	log.Info().Msg("Received shutdown signal.")
	log.Info().Int64("seconds", h.Config.Server.ShutdownPeriodSeconds).Msg("Entering pre-shutdown period.")
	time.Sleep(time.Duration(h.Config.Server.ShutdownPeriodSeconds) * time.Second)
	log.Info().Msg("Cleaning up all resources.")
	// TODO: next PR
	log.Info().Msg("Cleaning up completed.")
}

// SetupGracefulShutdown sets up graceful shutdown procedure for this service.
func (h *HTTP) SetupGracefulShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go h.Shutdown(done)
}

func (h *HTTP) Serve() {
	log.Info().Str("port", h.Config.Server.Port).Msg("HTTP server is running.")
	h.Router.Get("/health", h.HealthCheck)
	if h.Config.Server.Env == shared.DevEnvironment {
		docs.SwaggerInfo.Title = shared.ServiceName
		docs.SwaggerInfo.Version = shared.ServiceVersion
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", h.Config.App.URL)
		h.Router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))
		log.Info().Str("url", swaggerURL).Msg("Swagger documentation is available.")
	}

	err := netHttp.ListenAndServe(":"+h.Config.Server.Port, h.Router)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}

// HealthCheck performs a health check on the server. Usually required by
// Kubernetes to check if the service is healthy.
// @Summary Health Check
// @Description Health Check Endpoint
// @Tags service
// @Produce json
// @Accept json
// @Success 200 {object} response.Base
// @Failure 503 {object} response.Base
// @Router /health [get]
func (h *HTTP) HealthCheck(w netHttp.ResponseWriter, r *netHttp.Request) {
	if err := h.DB.Read.Ping(); err != nil {
		log.Error().Err(err).Msg("")
		response.WithUnhealthy(w)
		return
	}
	response.WithMessage(w, http.StatusOK, "OK")
}
