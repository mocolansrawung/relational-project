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
	"github.com/evermos/boilerplate-go/di"
	"github.com/evermos/boilerplate-go/docs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Http struct {
	Config *configs.Config
	DB     *infras.MysqlConn
	Router *chi.Mux
}

func (h *Http) Shutdown(done chan os.Signal, svc di.ServiceRegistry) {
	<-done
	defer os.Exit(0)
	log.Info().Msg("received signal shutdown...")
	time.Sleep(time.Duration(h.Config.ShutdownPeriod) * time.Second)
	log.Info().Msg("Clean up all resources...")
	svc.Shutdown()
	log.Info().Msg("Server shutdown properly...")
}

func (h *Http) GracefulShutdown(svc di.ServiceRegistry) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go h.Shutdown(done, svc)
}

func (h *Http) Serve() {
	log.Info().Str("port", h.Config.Port).Msg("Running HTTP server...")
	h.Router.Get("/health", h.HealthCheck)
	if h.Config.Env == shared.DevEnvironment {
		docs.SwaggerInfo.Title = shared.ServiceName
		docs.SwaggerInfo.Version = shared.ServiceVersion
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", h.Config.AppURL)
		h.Router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))
	}

	netHttp.ListenAndServe(":"+h.Config.Port, h.Router)
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
func (h *Http) HealthCheck(w netHttp.ResponseWriter, r *netHttp.Request) {
	if err := h.DB.Read.Ping(); err != nil {
		log.Error().Err(err).Msg("")
		response.WithUnhealthy(w)
		return
	}
	response.WithMessage(w, http.StatusOK, "OK")
}
