package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raksh/expense-tracker/models"
	"github.com/raksh/expense-tracker/services"
)

// Handler contains HTTP handlers that translate HTTP requests into service calls.
// It follows the adapter pattern — converting between HTTP and domain concerns.
type Handler struct {
	service *services.Service
}

// NewHandler creates a new Handler with the given service.
func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

// --- User Handlers ---

// CreateUser handles POST /api/users
// Creates a new user with the provided name and email.
func (h *Handler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid create user request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		slog.Error("Failed to create user", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create user",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// ListUsers handles GET /api/users
// Returns a list of all users.
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		slog.Error("Failed to list users", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to list users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// GetUserGroups handles GET /api/users/:id/groups
// Returns all groups the user belongs to.
func (h *Handler) GetUserGroups(c *gin.Context) {
	userID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	groups, err := h.service.GetUserGroups(userID)
	if err != nil {
		slog.Error("Failed to get user groups", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get user groups",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User groups retrieved successfully",
		Data:    groups,
	})
}

// --- Group Handlers ---

// CreateGroup handles POST /api/groups
// Creates a new group with the provided name.
func (h *Handler) CreateGroup(c *gin.Context) {
	var req models.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid create group request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	group, err := h.service.CreateGroup(req)
	if err != nil {
		slog.Error("Failed to create group", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create group",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Group created successfully",
		Data:    group,
	})
}

// AddUserToGroup handles POST /api/groups/:id/members
// Adds an existing user to an existing group.
func (h *Handler) AddUserToGroup(c *gin.Context) {
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req models.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid add member request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	group, err := h.service.AddUserToGroup(groupID, req.UserID)
	if err != nil {
		slog.Error("Failed to add user to group", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to add user to group",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User added to group successfully",
		Data:    group,
	})
}

// GetGroupMembers handles GET /api/groups/:id/members
// Returns all members of a specific group.
func (h *Handler) GetGroupMembers(c *gin.Context) {
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	members, err := h.service.GetGroupMembers(groupID)
	if err != nil {
		slog.Error("Failed to get group members", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get group members",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Group members retrieved successfully",
		Data:    members,
	})
}

// ListGroups handles GET /api/groups
// Returns all groups with member counts.
func (h *Handler) ListGroups(c *gin.Context) {
	groups, err := h.service.ListGroups()
	if err != nil {
		slog.Error("Failed to list groups", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to list groups",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Groups retrieved successfully",
		Data:    groups,
	})
}

// --- Expense Handlers ---

// AddExpense handles POST /api/groups/:id/expenses
// Creates a new expense in the group, split equally among all members.
func (h *Handler) AddExpense(c *gin.Context) {
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req models.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid create expense request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	expense, err := h.service.CreateExpense(groupID, req)
	if err != nil {
		slog.Error("Failed to create expense", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to create expense",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Expense created successfully",
		Data:    expense,
	})
}

// --- Balance & Settlement Handlers ---

// GetGroupBalances handles GET /api/groups/:id/balances
// Returns the net balance for each user in the group.
func (h *Handler) GetGroupBalances(c *gin.Context) {
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	balances, err := h.service.GetGroupBalances(groupID)
	if err != nil {
		slog.Error("Failed to get group balances", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to get group balances",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Balances retrieved successfully",
		Data:    balances,
	})
}

// GetSettlements handles GET /api/groups/:id/settlements
// Returns the minimum set of transactions needed to settle all debts.
func (h *Handler) GetSettlements(c *gin.Context) {
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	settlements, err := h.service.GetSettlements(groupID)
	if err != nil {
		slog.Error("Failed to get settlements", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to get settlements",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Settlements calculated successfully",
		Data:    settlements,
	})
}

// GetGroupExpenses handles GET /api/groups/:id/expenses
// Returns all expenses for a group with who paid and splits.
func (h *Handler) GetGroupExpenses(c *gin.Context) {
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	expenses, err := h.service.GetGroupExpenses(groupID)
	if err != nil {
		slog.Error("Failed to get group expenses", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to get group expenses",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Expenses retrieved successfully",
		Data:    expenses,
	})
}

// GetDashboardStats handles GET /api/dashboard/stats
// Returns aggregate statistics for the dashboard.
func (h *Handler) GetDashboardStats(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		slog.Error("Failed to get dashboard stats", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get dashboard stats",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Dashboard stats retrieved successfully",
		Data:    stats,
	})
}

// --- Helpers ---

// parseUintParam extracts a uint path parameter from the Gin context.
// It writes an error response and returns 0 if the parameter is invalid.
func parseUintParam(c *gin.Context, name string) (uint, error) {
	idStr := c.Param(name)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		slog.Warn("Invalid path parameter", "param", name, "value", idStr)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid " + name + " parameter",
			Details: "Must be a positive integer",
		})
		return 0, err
	}
	return uint(id), nil
}
