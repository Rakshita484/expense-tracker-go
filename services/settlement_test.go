package services

import (
	"testing"

	"github.com/raksh/expense-tracker/models"
	"github.com/shopspring/decimal"
)

// TestCalculateSettlements_TwoUsers tests settlement between two users
// where one user paid for the other.
//
// Scenario: Alice paid $100, split equally between Alice and Bob.
//   - Alice's balance: +$50 (she overpaid by $50)
//   - Bob's balance:   -$50 (he owes $50)
//   - Expected: Bob pays Alice $50 (1 transaction)
func TestCalculateSettlements_TwoUsers(t *testing.T) {
	balances := []models.UserBalance{
		{UserID: 1, Name: "Alice", Balance: decimal.NewFromFloat(50)},
		{UserID: 2, Name: "Bob", Balance: decimal.NewFromFloat(-50)},
	}

	settlements := CalculateSettlements(balances)

	if len(settlements) != 1 {
		t.Fatalf("expected 1 settlement, got %d", len(settlements))
	}

	s := settlements[0]
	if s.FromUserID != 2 || s.ToUserID != 1 {
		t.Errorf("expected Bob(2) -> Alice(1), got %d -> %d", s.FromUserID, s.ToUserID)
	}
	if !s.Amount.Equal(decimal.NewFromFloat(50)) {
		t.Errorf("expected amount $50.00, got %s", s.Amount.String())
	}
}

// TestCalculateSettlements_ThreeUsersEqual tests settlement among three users
// with equal split.
//
// Scenario: Alice paid $300, split equally among Alice, Bob, Charlie.
//   - Each person's share: $100
//   - Alice's balance: +$200 (she paid $300, owes $100)
//   - Bob's balance:   -$100 (he paid $0, owes $100)
//   - Charlie's balance: -$100 (he paid $0, owes $100)
//   - Expected: Bob pays Alice $100, Charlie pays Alice $100 (2 transactions)
func TestCalculateSettlements_ThreeUsersEqual(t *testing.T) {
	balances := []models.UserBalance{
		{UserID: 1, Name: "Alice", Balance: decimal.NewFromFloat(200)},
		{UserID: 2, Name: "Bob", Balance: decimal.NewFromFloat(-100)},
		{UserID: 3, Name: "Charlie", Balance: decimal.NewFromFloat(-100)},
	}

	settlements := CalculateSettlements(balances)

	if len(settlements) != 2 {
		t.Fatalf("expected 2 settlements, got %d", len(settlements))
	}

	// Total settlement amount should equal $200.
	total := decimal.Zero
	for _, s := range settlements {
		total = total.Add(s.Amount)
	}
	if !total.Equal(decimal.NewFromFloat(200)) {
		t.Errorf("expected total settlement $200, got %s", total.String())
	}
}

// TestCalculateSettlements_ZeroBalances tests that no settlements are generated
// when all balances are zero (everyone is settled).
func TestCalculateSettlements_ZeroBalances(t *testing.T) {
	balances := []models.UserBalance{
		{UserID: 1, Name: "Alice", Balance: decimal.Zero},
		{UserID: 2, Name: "Bob", Balance: decimal.Zero},
		{UserID: 3, Name: "Charlie", Balance: decimal.Zero},
	}

	settlements := CalculateSettlements(balances)

	if len(settlements) != 0 {
		t.Fatalf("expected 0 settlements, got %d", len(settlements))
	}
}

// TestCalculateSettlements_SingleUser tests that a single user generates no settlements.
func TestCalculateSettlements_SingleUser(t *testing.T) {
	balances := []models.UserBalance{
		{UserID: 1, Name: "Alice", Balance: decimal.Zero},
	}

	settlements := CalculateSettlements(balances)

	if len(settlements) != 0 {
		t.Fatalf("expected 0 settlements, got %d", len(settlements))
	}
}

// TestCalculateSettlements_ComplexMultiUser tests a complex scenario with
// multiple creditors and debtors.
//
// Scenario: 4 users with the following balances after multiple expenses:
//   - Alice: +$60 (owed $60)
//   - Bob:   +$30 (owed $30)
//   - Charlie: -$50 (owes $50)
//   - Dave:    -$40 (owes $40)
//
// Expected settlements (greedy, largest first):
//   1. Charlie pays Alice $50 (Charlie settled, Alice has $10 left)
//   2. Dave pays Alice $10 (Alice settled, Dave has $30 left)
//   3. Dave pays Bob $30 (Both settled)
// Total: 3 transactions
func TestCalculateSettlements_ComplexMultiUser(t *testing.T) {
	balances := []models.UserBalance{
		{UserID: 1, Name: "Alice", Balance: decimal.NewFromFloat(60)},
		{UserID: 2, Name: "Bob", Balance: decimal.NewFromFloat(30)},
		{UserID: 3, Name: "Charlie", Balance: decimal.NewFromFloat(-50)},
		{UserID: 4, Name: "Dave", Balance: decimal.NewFromFloat(-40)},
	}

	settlements := CalculateSettlements(balances)

	if len(settlements) != 3 {
		t.Fatalf("expected 3 settlements, got %d", len(settlements))
	}

	// Verify total amounts balance out.
	totalFromDebtors := decimal.Zero
	for _, s := range settlements {
		totalFromDebtors = totalFromDebtors.Add(s.Amount)
	}
	expectedTotal := decimal.NewFromFloat(90) // $60 + $30 = $90
	if !totalFromDebtors.Equal(expectedTotal) {
		t.Errorf("expected total settlement %s, got %s", expectedTotal.String(), totalFromDebtors.String())
	}

	// Verify first settlement: Charlie -> Alice $50.
	if settlements[0].FromUserID != 3 || settlements[0].ToUserID != 1 {
		t.Errorf("settlement[0]: expected Charlie(3) -> Alice(1), got %d -> %d",
			settlements[0].FromUserID, settlements[0].ToUserID)
	}
	if !settlements[0].Amount.Equal(decimal.NewFromFloat(50)) {
		t.Errorf("settlement[0]: expected $50, got %s", settlements[0].Amount.String())
	}
}

// TestCalculateSettlements_Empty tests that an empty balance list returns no settlements.
func TestCalculateSettlements_Empty(t *testing.T) {
	settlements := CalculateSettlements(nil)

	if len(settlements) != 0 {
		t.Fatalf("expected 0 settlements, got %d", len(settlements))
	}
}

// TestCalculateSettlements_PennyPrecision tests that the algorithm correctly handles
// small amounts typical of penny-rounding scenarios.
func TestCalculateSettlements_PennyPrecision(t *testing.T) {
	balances := []models.UserBalance{
		{UserID: 1, Name: "Alice", Balance: decimal.NewFromFloat(0.01)},
		{UserID: 2, Name: "Bob", Balance: decimal.NewFromFloat(-0.01)},
	}

	settlements := CalculateSettlements(balances)

	if len(settlements) != 1 {
		t.Fatalf("expected 1 settlement, got %d", len(settlements))
	}
	if !settlements[0].Amount.Equal(decimal.NewFromFloat(0.01)) {
		t.Errorf("expected $0.01, got %s", settlements[0].Amount.String())
	}
}
