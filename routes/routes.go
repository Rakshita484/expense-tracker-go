package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raksh/expense-tracker/handlers"
)

// RegisterRoutes sets up all API routes on the Gin engine.
// Routes are grouped under /api for clean namespacing.
func RegisterRoutes(router *gin.Engine, h *handlers.Handler) {
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", h.CreateUser)              // POST /api/users
			users.GET("", h.ListUsers)                // GET  /api/users
			users.GET("/:id/groups", h.GetUserGroups) // GET /api/users/:id/groups
		}

		// Group routes
		groups := api.Group("/groups")
		{
			groups.GET("", h.ListGroups)                     // GET  /api/groups
			groups.POST("", h.CreateGroup)                   // POST /api/groups
			groups.POST("/:id/members", h.AddUserToGroup)    // POST /api/groups/:id/members
			groups.GET("/:id/members", h.GetGroupMembers)    // GET  /api/groups/:id/members
			groups.POST("/:id/expenses", h.AddExpense)       // POST /api/groups/:id/expenses
			groups.GET("/:id/expenses", h.GetGroupExpenses)  // GET  /api/groups/:id/expenses
			groups.GET("/:id/balances", h.GetGroupBalances)  // GET  /api/groups/:id/balances
			groups.GET("/:id/settlements", h.GetSettlements) // GET  /api/groups/:id/settlements
		}

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/stats", h.GetDashboardStats) // GET /api/dashboard/stats
		}
	}
}
