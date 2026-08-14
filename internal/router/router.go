package router

import (
	"github.com/gin-gonic/gin"

	"github.com/coolkit-org/coolkit/internal/handler"
	"github.com/coolkit-org/coolkit/internal/middleware"
)

func Setup(healthHandler *handler.HealthHandler, authHandler *handler.AuthHandler, clubHandler *handler.ClubHandler, eventHandler *handler.EventHandler, eventRoleHandler *handler.EventRoleHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)
		v1.GET("/health/db", healthHandler.DBHealth)
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwtSecret))
		{
			protected.GET("/auth/me", authHandler.Me)
			protected.POST("/clubs", clubHandler.Create)
			protected.GET("/clubs", clubHandler.List)
			protected.GET("/clubs/:id", clubHandler.Get)
			protected.POST("/clubs/:id/join", clubHandler.Join)
			protected.GET("/clubs/:id/members", clubHandler.Members)
			protected.POST("/clubs/:id/events", eventHandler.Create)
			protected.GET("/clubs/:id/events", eventHandler.List)
			protected.GET("/events/:id", eventRoleHandler.GetEventDetails)
			protected.GET("/events/:id/roles", eventRoleHandler.GetRoles)
			protected.POST("/events/:id/roles", eventRoleHandler.AssignRole)
		}
	}

	r.Static("/assets", "./frontend/assets")
	r.Static("/css", "./frontend/css")
	r.Static("/js", "./frontend/js")

	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/dashboard", "./frontend/dashboard.html")
	r.StaticFile("/club", "./frontend/club.html")
	r.StaticFile("/event", "./frontend/event.html")

	return r
}
