package server

import (
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func New(debug, cors bool, resolver *graph_resolver.Resolver, logger *logger.Adapter) *HttpServer {
	return &HttpServer{
		debug:    debug,
		cors:     cors,
		log:      logger,
		resolver: resolver,
	}
}
