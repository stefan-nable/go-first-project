package routes

import (
	"database/sql"
	"github.com/emicklei/go-restful"
)

type Router struct {
	db *sql.DB
}

func NewRouter(db *sql.DB) *Router {
	return &Router{db: db}
}

func (r *Router) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/api/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.POST("").To(r.PostWorkersHandler).Doc("Starts processing with the given number of workers."))
}
