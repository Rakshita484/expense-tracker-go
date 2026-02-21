package repository

import (
	"fmt"
	"log/slog"

	"github.com/raksh/expense-tracker/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Repository provides data access methods for the expense tracker.
// It encapsulates all database operations behind a clean interface.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository with the given GORM database connection.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// --- User Operations ---

// CreateUser inserts a new user into the database.
func (r *Repository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	slog.Info("User created", "id", user.ID, "name", user.Name)
	return nil
}

// GetUsers retrieves all users from the database.
func (r *Repository) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

// GetUserByID retrieves a single user by their ID.
func (r *Repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found with id %d: %w", id, err)
	}
	return &user, nil
}

// --- Group Operations ---

// CreateGroup inserts a new group into the database.
func (r *Repository) CreateGroup(group *models.Group) error {
	if err := r.db.Create(group).Error; err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	slog.Info("Group created", "id", group.ID, "name", group.Name)
	return nil
}

// GetGroupByID retrieves a group by ID with its members preloaded.
func (r *Repository) GetGroupByID(id uint) (*models.Group, error) {
	var group models.Group
	if err := r.db.Preload("Members").First(&group, id).Error; err != nil {
		return nil, fmt.Errorf("group not found with id %d: %w", id, err)
	}
	return &group, nil
}

// AddUserToGroup adds a user to a group's member list.
func (r *Repository) AddUserToGroup(group *models.Group, user *models.User) error {
	if err := r.db.Model(group).Association("Members").Append(user); err != nil {
		return fmt.Errorf("failed to add user %d to group %d: %w", user.ID, group.ID, err)
	}
	slog.Info("User added to group", "user_id", user.ID, "group_id", group.ID)
	return nil
}

// IsUserInGroup checks whether a user is a member of a given group.
func (r *Repository) IsUserInGroup(groupID, userID uint) (bool, error) {
	var count int64
	err := r.db.Table("group_members").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}
	return count > 0, nil
}

// --- Expense Operations ---

// CreateExpenseWithSplits creates an expense and its splits within a single database
// transaction. This ensures atomicity — either both the expense and all splits are
// persisted, or neither is. This prevents data inconsistency from partial failures
// and protects against race conditions through PostgreSQL's MVCC isolation.
func (r *Repository) CreateExpenseWithSplits(expense *models.Expense, splits []models.ExpenseSplit) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Insert the expense record.
		if err := tx.Create(expense).Error; err != nil {
			return fmt.Errorf("failed to create expense: %w", err)
		}

		// Step 2: Assign the new expense ID to each split, then insert them.
		for i := range splits {
			splits[i].ExpenseID = expense.ID
		}
		if err := tx.Create(&splits).Error; err != nil {
			return fmt.Errorf("failed to create expense splits: %w", err)
		}

		slog.Info("Expense created with splits",
			"expense_id", expense.ID,
			"splits_count", len(splits),
			"amount", expense.Amount.String(),
		)
		return nil
	})
}

// GetGroupBalances calculates the net balance for each user in a group.
//
// The balance is computed as:
//   total_paid (sum of expenses the user paid for)
//   - total_owed (sum of splits assigned to the user)
//
// A positive balance means the user is owed money by others.
// A negative balance means the user owes money to others.
//
// This uses a single efficient SQL query with aggregation rather than
// loading all expenses into memory.
func (r *Repository) GetGroupBalances(groupID uint) ([]models.UserBalance, error) {
	var balances []models.UserBalance

	// This query:
	// 1. Starts with all members of the group (group_members join table).
	// 2. LEFT JOINs expenses to find how much each member paid.
	// 3. LEFT JOINs expense_splits to find how much each member owes.
	// 4. Computes paid - owed as net balance.
	// 5. Groups by user to get one row per member.
	query := `
		SELECT 
			u.id AS user_id,
			u.name AS name,
			COALESCE(paid.total_paid, 0) - COALESCE(owed.total_owed, 0) AS balance
		FROM group_members gm
		JOIN users u ON u.id = gm.user_id
		LEFT JOIN (
			SELECT paid_by_id, SUM(amount) AS total_paid
			FROM expenses
			WHERE group_id = ?
			GROUP BY paid_by_id
		) paid ON paid.paid_by_id = u.id
		LEFT JOIN (
			SELECT es.user_id, SUM(es.amount) AS total_owed
			FROM expense_splits es
			JOIN expenses e ON e.id = es.expense_id
			WHERE e.group_id = ?
			GROUP BY es.user_id
		) owed ON owed.user_id = u.id
		WHERE gm.group_id = ?
		ORDER BY balance DESC
	`

	if err := r.db.Raw(query, groupID, groupID, groupID).Scan(&balances).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate balances: %w", err)
	}

	// Convert raw numeric values to proper decimal.Decimal for precision.
	for i := range balances {
		balances[i].Balance = balances[i].Balance.Round(2)
	}

	return balances, nil
}

// GetGroupMembers returns the count of members in a group.
func (r *Repository) GetGroupMembers(groupID uint) ([]models.User, error) {
	var group models.Group
	if err := r.db.Preload("Members").First(&group, groupID).Error; err != nil {
		return nil, fmt.Errorf("group not found with id %d: %w", groupID, err)
	}
	return group.Members, nil
}

// DB returns the underlying *gorm.DB for advanced use cases (e.g., migrations).
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// --- Helper: Check if decimal is positive ---

func isPositive(d decimal.Decimal) bool {
	return d.GreaterThan(decimal.Zero)
}
