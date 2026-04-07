package server

import (
	"net/http"
	"time"

	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

const AUTH_TIMEOUT = 24 * time.Hour

type HttpServer struct {
	debug bool
	cors  bool

	resolver *graph_resolver.Resolver
	log      *logger.Adapter
	engine   *gin.Engine
	server   http.Server
}
