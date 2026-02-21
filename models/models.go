package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// User represents a user in the expense tracking system.
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"uniqueIndex;not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Group represents a group where members can share expenses.
type Group struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`

	// Many-to-many relationship with User via the "group_members" join table.
	Members []User `json:"members,omitempty" gorm:"many2many:group_members;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Expense represents a single expense paid by one user in a group.
// Amount uses shopspring/decimal to avoid floating-point precision errors.
type Expense struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	GroupID     uint            `json:"group_id" gorm:"not null;index"`
	Description string          `json:"description" gorm:"not null"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:numeric(15,2);not null"`
	PaidByID    uint            `json:"paid_by_id" gorm:"not null"`

	// Relationships
	Group  Group  `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	PaidBy User   `json:"paid_by,omitempty" gorm:"foreignKey:PaidByID"`
	Splits []ExpenseSplit `json:"splits,omitempty" gorm:"foreignKey:ExpenseID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExpenseSplit represents how much each user owes for an expense.
// This is stored as a separate record to support flexible splitting strategies.
type ExpenseSplit struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	ExpenseID uint            `json:"expense_id" gorm:"not null;index"`
	UserID    uint            `json:"user_id" gorm:"not null;index"`
	Amount    decimal.Decimal `json:"amount" gorm:"type:numeric(15,2);not null"`

	// Relationships
	Expense Expense `json:"-" gorm:"foreignKey:ExpenseID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
}

// UserBalance represents the net balance of a user within a group.
// A positive balance means the user is owed money; negative means the user owes money.
type UserBalance struct {
	UserID  uint            `json:"user_id"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

// Settlement represents a single settlement transaction between two users.
type Settlement struct {
	FromUserID   uint            `json:"from_user_id"`
	FromUserName string          `json:"from_user_name"`
	ToUserID     uint            `json:"to_user_id"`
	ToUserName   string          `json:"to_user_name"`
	Amount       decimal.Decimal `json:"amount"`
}

// --- Request DTOs ---

// CreateUserRequest is the request body for creating a new user.
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// CreateGroupRequest is the request body for creating a new group.
type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

// AddMemberRequest is the request body for adding a user to a group.
type AddMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// CreateExpenseRequest is the request body for creating a new expense.
type CreateExpenseRequest struct {
	Description string          `json:"description" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	PaidByID    uint            `json:"paid_by_id" binding:"required"`
}

// --- Response DTOs ---

// ErrorResponse provides a structured JSON error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse provides a structured JSON success response.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
