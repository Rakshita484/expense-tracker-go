package services

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/raksh/expense-tracker/models"
	"github.com/raksh/expense-tracker/repository"
	"github.com/shopspring/decimal"
)

// Service contains the business logic for the expense tracker.
// It sits between the HTTP handlers and the repository layer,
// enforcing validation rules and implementing domain logic.
type Service struct {
	repo *repository.Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// --- User Operations ---

// CreateUser validates and creates a new user.
func (s *Service) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers returns all users.
func (s *Service) GetUsers() ([]models.User, error) {
	return s.repo.GetUsers()
}

// --- Group Operations ---

// CreateGroup validates and creates a new group.
func (s *Service) CreateGroup(req models.CreateGroupRequest) (*models.Group, error) {
	group := &models.Group{
		Name: req.Name,
	}
	if err := s.repo.CreateGroup(group); err != nil {
		return nil, err
	}
	return group, nil
}

// AddUserToGroup adds an existing user to an existing group.
func (s *Service) AddUserToGroup(groupID, userID uint) (*models.Group, error) {
	group, err := s.repo.GetGroupByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is already in the group.
	isMember, err := s.repo.IsUserInGroup(groupID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, fmt.Errorf("user %d is already a member of group %d", userID, groupID)
	}

	if err := s.repo.AddUserToGroup(group, user); err != nil {
		return nil, err
	}

	// Reload group with updated members.
	return s.repo.GetGroupByID(groupID)
}

// --- Expense Operations ---

// CreateExpense validates inputs, computes equal splits, and creates the expense
// with all splits in a single database transaction.
//
// Split Calculation:
//   For an amount of $100.00 split among 3 users:
//   - Base split: 100.00 / 3 = 33.33 (truncated to 2 decimal places)
//   - Remainder: 100.00 - (33.33 × 3) = 100.00 - 99.99 = 0.01
//   - First user gets 33.34, second and third get 33.33
//   - This ensures the splits always sum exactly to the original amount.
func (s *Service) CreateExpense(groupID uint, req models.CreateExpenseRequest) (*models.Expense, error) {
	// Validate the amount is positive.
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("amount must be positive, got %s", req.Amount.String())
	}

	// Validate the group exists.
	group, err := s.repo.GetGroupByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	// Validate there are members in the group.
	if len(group.Members) == 0 {
		return nil, fmt.Errorf("group has no members to split the expense")
	}

	// Validate the payer is a member of the group.
	payerIsMember := false
	for _, m := range group.Members {
		if m.ID == req.PaidByID {
			payerIsMember = true
			break
		}
	}
	if !payerIsMember {
		return nil, fmt.Errorf("payer (user %d) is not a member of group %d", req.PaidByID, groupID)
	}

	// Calculate equal splits using integer-safe decimal arithmetic.
	memberCount := decimal.NewFromInt(int64(len(group.Members)))

	// Truncate to 2 decimal places (floor division for currency).
	perPersonAmount := req.Amount.Div(memberCount).Truncate(2)

	// Calculate the remainder to ensure splits sum exactly to the total amount.
	// This handles cases like $100 / 3 = $33.33 each, with $0.01 remainder.
	totalSplitSoFar := perPersonAmount.Mul(memberCount)
	remainder := req.Amount.Sub(totalSplitSoFar)

	slog.Info("Calculating expense splits",
		"total", req.Amount.String(),
		"members", len(group.Members),
		"per_person", perPersonAmount.String(),
		"remainder", remainder.String(),
	)

	// Build the split records.
	splits := make([]models.ExpenseSplit, len(group.Members))
	for i, member := range group.Members {
		splitAmount := perPersonAmount

		// Distribute the remainder (penny by penny) to the first users.
		// This ensures the total always matches exactly.
		if remainder.GreaterThan(decimal.Zero) {
			splitAmount = splitAmount.Add(decimal.NewFromFloat(0.01))
			remainder = remainder.Sub(decimal.NewFromFloat(0.01))
		}

		splits[i] = models.ExpenseSplit{
			UserID: member.ID,
			Amount: splitAmount,
		}
	}

	// Create the expense record.
	expense := &models.Expense{
		GroupID:      groupID,
		Description:  req.Description,
		Amount:       req.Amount,
		PaidByID:     req.PaidByID,
	}

	// Persist expense and splits atomically within a transaction.
	// This prevents partial writes — either everything succeeds or nothing does.
	if err := s.repo.CreateExpenseWithSplits(expense, splits); err != nil {
		return nil, err
	}

	expense.Splits = splits
	return expense, nil
}

// --- Balance & Settlement Operations ---

// GetGroupBalances returns the net balance for each user in a group.
func (s *Service) GetGroupBalances(groupID uint) ([]models.UserBalance, error) {
	// Validate group exists.
	if _, err := s.repo.GetGroupByID(groupID); err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}
	return s.repo.GetGroupBalances(groupID)
}

// GetSettlements calculates the minimum number of transactions needed
// to settle all debts within a group.
func (s *Service) GetSettlements(groupID uint) ([]models.Settlement, error) {
	balances, err := s.GetGroupBalances(groupID)
	if err != nil {
		return nil, err
	}
	return CalculateSettlements(balances), nil
}

// CalculateSettlements implements a greedy algorithm to minimize the number
// of settlement transactions needed to resolve all debts.
//
// Algorithm Overview:
//   1. Separate users into creditors (positive balance = owed money)
//      and debtors (negative balance = owe money).
//   2. Sort creditors descending by balance and debtors ascending (most debt first).
//   3. Match the largest creditor with the largest debtor:
//      - Transfer amount = min(credit, |debt|)
//      - Reduce both balances by the transfer amount.
//      - If a balance reaches zero, move to the next creditor/debtor.
//   4. Repeat until all balances are settled.
//
// Time Complexity: O(n log n) for sorting + O(n) for settlement = O(n log n) overall,
// where n is the number of users with non-zero balances.
//
// Space Complexity: O(n) for storing creditors and debtors.
//
// This is a pure function with no side effects, making it easy to test independently.
func CalculateSettlements(balances []models.UserBalance) []models.Settlement {
	// Separate users into creditors and debtors based on their net balance.
	var creditors []models.UserBalance // Users who are owed money (positive balance)
	var debtors []models.UserBalance   // Users who owe money (negative balance)

	// A small threshold to handle rounding errors (values below $0.01 are zero).
	threshold := decimal.NewFromFloat(0.01)

	for _, b := range balances {
		if b.Balance.GreaterThanOrEqual(threshold) {
			creditors = append(creditors, b)
		} else if b.Balance.LessThanOrEqual(threshold.Neg()) {
			debtors = append(debtors, b)
		}
		// Users with zero balance require no settlement.
	}

	// Sort creditors in descending order (largest credit first).
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].Balance.GreaterThan(creditors[j].Balance)
	})

	// Sort debtors in ascending order (largest debt first, since debts are negative).
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].Balance.LessThan(debtors[j].Balance)
	})

	var settlements []models.Settlement

	// Use two pointers to greedily match the largest creditor with the largest debtor.
	// This minimizes the number of transactions needed.
	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		credit := creditors[i].Balance
		debt := debtors[j].Balance.Abs() // Convert negative to positive for comparison.

		// The settlement amount is the smaller of the two absolute values.
		// This ensures we don't over-settle either party.
		var amount decimal.Decimal
		if credit.LessThanOrEqual(debt) {
			amount = credit
		} else {
			amount = debt
		}

		// Only create a settlement if the amount is meaningful (>= $0.01).
		if amount.GreaterThanOrEqual(threshold) {
			settlements = append(settlements, models.Settlement{
				FromUserID:   debtors[j].UserID,
				FromUserName: debtors[j].Name,
				ToUserID:     creditors[i].UserID,
				ToUserName:   creditors[i].Name,
				Amount:       amount.Round(2),
			})
		}

		// Reduce both balances by the settlement amount.
		creditors[i].Balance = creditors[i].Balance.Sub(amount)
		debtors[j].Balance = debtors[j].Balance.Add(amount)

		// If the creditor is fully settled, move to the next one.
		if creditors[i].Balance.LessThan(threshold) {
			i++
		}

		// If the debtor is fully settled, move to the next one.
		if debtors[j].Balance.Abs().LessThan(threshold) {
			j++
		}
	}

	slog.Info("Settlements calculated",
		"creditors", len(creditors),
		"debtors", len(debtors),
		"transactions", len(settlements),
	)

	return settlements
}
