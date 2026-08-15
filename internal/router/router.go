package router

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coolkit-org/coolkit/internal/handler"
	"github.com/coolkit-org/coolkit/internal/middleware"
)

func Setup(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	clubHandler *handler.ClubHandler,
	memberHandler *handler.MemberHandler,
	eventHandler *handler.EventHandler,
	eventRoleHandler *handler.EventRoleHandler,
	announcementHandler *handler.AnnouncementHandler,
	taskHandler *handler.TaskHandler,
	financeHandler *handler.FinanceHandler,
	hierarchyHandler *handler.HierarchyHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)
		v1.GET("/health/db", healthHandler.DBHealth)
		
		// Apply rate limiting to auth endpoints (5 requests per minute)
		authGroup := v1.Group("")
		authGroup.Use(middleware.RateLimit(5, time.Minute))
		{
			authGroup.POST("/auth/register", authHandler.Register)
			authGroup.POST("/auth/login", authHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwtSecret))
		{
			// Auth
			protected.GET("/auth/me", authHandler.Me)

			// Clubs
			protected.POST("/clubs", clubHandler.Create)
			protected.GET("/clubs", clubHandler.List)
			protected.GET("/clubs/:id", clubHandler.Get)
			protected.PUT("/clubs/:id/settings", clubHandler.UpdateSettings)
			protected.POST("/clubs/join", clubHandler.Join)
			protected.GET("/clubs/:id/members", clubHandler.Members)
			
			// Join Requests
			protected.GET("/clubs/:id/requests", clubHandler.GetJoinRequests)
			protected.POST("/clubs/:id/requests/:request_id/approve", clubHandler.ApproveJoinRequest)
			protected.POST("/clubs/:id/requests/:request_id/reject", clubHandler.RejectJoinRequest)

			// Member Management
			protected.PUT("/clubs/:id/members/:user_id/role", memberHandler.UpdateRole)
			protected.DELETE("/clubs/:id/members/:user_id", memberHandler.Remove)
			protected.DELETE("/clubs/:id/members/me", memberHandler.Leave)

			// Announcements
			protected.POST("/clubs/:id/announcements", announcementHandler.Create)
			protected.GET("/clubs/:id/announcements", announcementHandler.List)
			protected.DELETE("/announcements/:id", announcementHandler.Delete)

			// Events (under clubs)
			protected.POST("/clubs/:id/events", eventHandler.Create)
			protected.GET("/clubs/:id/events", eventHandler.List)

			// Events (standalone)
			protected.GET("/events/:id", eventRoleHandler.GetEventDetails)
			protected.PUT("/events/:id", eventHandler.Update)
			protected.DELETE("/events/:id", eventHandler.Delete)

			// Event Roles
			protected.GET("/events/:id/roles", eventRoleHandler.GetRoles)
			protected.POST("/events/:id/roles", eventRoleHandler.AssignRole)
			protected.DELETE("/events/:id/roles/:role_id", eventRoleHandler.RemoveRole)

			// Tasks (per event)
			protected.POST("/events/:id/tasks", taskHandler.Create)
			protected.GET("/events/:id/tasks", taskHandler.List)
			protected.PATCH("/tasks/:id/status", taskHandler.UpdateStatus)
			protected.DELETE("/tasks/:id", taskHandler.Delete)

			// Finance (per event)
			protected.POST("/events/:id/finance", financeHandler.Create)
			protected.GET("/events/:id/finance", financeHandler.List)
			protected.DELETE("/finance/:id", financeHandler.Delete)

			// Hierarchy & Domains
			protected.POST("/clubs/:id/hierarchy", hierarchyHandler.SetHierarchy)
			protected.GET("/clubs/:id/hierarchy", hierarchyHandler.GetHierarchy)
			protected.POST("/clubs/:id/domains", hierarchyHandler.CreateDomain)
			protected.GET("/clubs/:id/domains", hierarchyHandler.ListDomains)
			protected.PUT("/clubs/:id/domains/:domain_id", hierarchyHandler.UpdateDomain)
			protected.DELETE("/clubs/:id/domains/:domain_id", hierarchyHandler.DeleteDomain)
			protected.PUT("/clubs/:id/members/:user_id/organization", hierarchyHandler.AssignMemberOrganization)
			protected.GET("/clubs/:id/organization", hierarchyHandler.GetOrganizationTree)
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
