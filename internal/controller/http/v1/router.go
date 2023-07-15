package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"id-maker/internal/usecase"
	"id-maker/pkg/logger"
	"net/http"
)

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Segment) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")

	handler.GET("/swagger/*any", swaggerHandler)
	handler.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.Group("/v1")
	{
		newSegmentRoutes(h, t, l)
	}
}
