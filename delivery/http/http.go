package http

import (
	netHttp "net/http"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared"

	"github.com/go-chi/chi"
)

type Http struct {
	Config *configs.Config
	DB     *infras.MysqlConn
	Router *chi.Mux
}

func (h *Http) Serve() {
	h.Router.Get("/health", h.HealthCheck)
	netHttp.ListenAndServe(":"+h.Config.Port, h.Router)
}

func (h *Http) HealthCheck(w netHttp.ResponseWriter, r *netHttp.Request) {
	if err := h.DB.Read.Ping(); err != nil {
		shared.Failed(w, r, shared.WriteFailed(500, "Server is unhealthy", nil))
		return
	}
	shared.Success(w, r, shared.WriteSuccess(200, "Server is alive", nil))
}
