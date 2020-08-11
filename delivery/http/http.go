package http

import (
	"fmt"
	netHttp "net/http"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/docs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Http struct {
	Config *configs.Config
	DB     *infras.MysqlConn
	Router *chi.Mux
}

func (h *Http) Serve() {
	h.Router.Get("/health", h.HealthCheck)

	if h.Config.Env == shared.EnvirontmentDev {
		docs.SwaggerInfo.Title = shared.ServiceName
		docs.SwaggerInfo.Version = shared.ServiceVersion
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", h.Config.AppURL)
		h.Router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))
	}

	netHttp.ListenAndServe(":"+h.Config.Port, h.Router)
}

func (h *Http) HealthCheck(w netHttp.ResponseWriter, r *netHttp.Request) {
	if err := h.DB.Read.Ping(); err != nil {
		shared.Failed(w, r, shared.WriteFailed(500, "Server is unhealthy", nil))
		return
	}
	shared.Success(w, r, shared.WriteSuccess(200, "Server is alive", nil))
}
