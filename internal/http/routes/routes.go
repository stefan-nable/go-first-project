package routes

import (
	"Main/internal/worker"
	"database/sql"
	"github.com/emicklei/go-restful"
)

type Router struct {
	db     *sql.DB
	Worker worker.W
}

func NewRouter(db *sql.DB) *Router {
	return &Router{db: db, Worker: &worker.Worker{}}
}

func (r *Router) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/api/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.POST("").To(r.PostWorkersHandler).Doc("Starts processing with the given number of workers."))
}
