package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/raksh/expense-tracker/models"
	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// main seeds the database with sample data for development and testing.
// Run this with: go run seed/seed.go
func main() {
	slog.Info("Starting database seeder")

	// Load .env configuration.
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate to ensure tables exist.
	if err := db.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Expense{},
		&models.ExpenseSplit{},
	); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	slog.Info("Creating seed data...")

	// --- Create Users ---
	users := []models.User{
		{Name: "Alice Johnson", Email: "alice@example.com"},
		{Name: "Bob Smith", Email: "bob@example.com"},
		{Name: "Charlie Brown", Email: "charlie@example.com"},
		{Name: "Diana Prince", Email: "diana@example.com"},
	}

	for i := range users {
		result := db.Where("email = ?", users[i].Email).FirstOrCreate(&users[i])
		if result.Error != nil {
			log.Fatalf("Failed to create user %s: %v", users[i].Name, result.Error)
		}
		slog.Info("User ready", "name", users[i].Name, "id", users[i].ID)
	}

	// --- Create Group ---
	group := models.Group{Name: "Weekend Trip"}
	result := db.Where("name = ?", group.Name).FirstOrCreate(&group)
	if result.Error != nil {
		log.Fatalf("Failed to create group: %v", result.Error)
	}
	slog.Info("Group ready", "name", group.Name, "id", group.ID)

	// --- Add Members to Group ---
	for i := range users {
		if err := db.Model(&group).Association("Members").Append(&users[i]); err != nil {
			slog.Warn("Failed to add member (may already exist)", "user", users[i].Name, "error", err)
		}
	}
	slog.Info("Members added to group", "count", len(users))

	// --- Create Sample Expenses ---
	// Expense 1: Alice paid $120 for dinner, split 4 ways.
	expense1 := models.Expense{
		GroupID:     group.ID,
		Description: "Dinner at Italian Restaurant",
		Amount:      decimal.NewFromFloat(120.00),
		PaidByID:    users[0].ID, // Alice
	}
	if err := createExpenseIfNotExists(db, &expense1, users); err != nil {
		log.Fatalf("Failed to create expense 1: %v", err)
	}

	// Expense 2: Bob paid $80 for gas, split 4 ways.
	expense2 := models.Expense{
		GroupID:     group.ID,
		Description: "Gas for Road Trip",
		Amount:      decimal.NewFromFloat(80.00),
		PaidByID:    users[1].ID, // Bob
	}
	if err := createExpenseIfNotExists(db, &expense2, users); err != nil {
		log.Fatalf("Failed to create expense 2: %v", err)
	}

	// Expense 3: Charlie paid $45.99 for snacks, split 4 ways.
	expense3 := models.Expense{
		GroupID:     group.ID,
		Description: "Snacks and Drinks",
		Amount:      decimal.NewFromFloat(45.99),
		PaidByID:    users[2].ID, // Charlie
	}
	if err := createExpenseIfNotExists(db, &expense3, users); err != nil {
		log.Fatalf("Failed to create expense 3: %v", err)
	}

	slog.Info("Seed data created successfully!")
	fmt.Println("\n=== Seed Summary ===")
	fmt.Println("Users: Alice, Bob, Charlie, Diana")
	fmt.Println("Group: Weekend Trip (all 4 members)")
	fmt.Println("Expenses:")
	fmt.Println("  1. Dinner ($120.00) - paid by Alice")
	fmt.Println("  2. Gas ($80.00) - paid by Bob")
	fmt.Println("  3. Snacks ($45.99) - paid by Charlie")
	fmt.Println("\nRun the API server to check balances and settlements!")
}

// connectDB connects to the configured database (SQLite or PostgreSQL).
func connectDB() (*gorm.DB, error) {
	driver := getEnv("DB_DRIVER", "sqlite")
	dbName := getEnv("DB_NAME", "expense_tracker")

	var dialector gorm.Dialector
	switch driver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			dbName,
		)
		dialector = postgres.Open(dsn)
	default:
		dialector = sqlite.Open(dbName + ".db")
	}

	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// createExpenseIfNotExists creates an expense with equal splits if it doesn't
// already exist (checked by description and group).
func createExpenseIfNotExists(db *gorm.DB, expense *models.Expense, members []models.User) error {
	var count int64
	db.Model(&models.Expense{}).
		Where("description = ? AND group_id = ?", expense.Description, expense.GroupID).
		Count(&count)
	if count > 0 {
		slog.Info("Expense already exists, skipping", "description", expense.Description)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(expense).Error; err != nil {
			return err
		}

		memberCount := decimal.NewFromInt(int64(len(members)))
		perPerson := expense.Amount.Div(memberCount).Truncate(2)
		remainder := expense.Amount.Sub(perPerson.Mul(memberCount))

		for i, m := range members {
			amount := perPerson
			if i == 0 && remainder.GreaterThan(decimal.Zero) {
				amount = amount.Add(remainder)
			}
			split := models.ExpenseSplit{
				ExpenseID: expense.ID,
				UserID:    m.ID,
				Amount:    amount,
			}
			if err := tx.Create(&split).Error; err != nil {
				return err
			}
		}

		slog.Info("Expense created", "description", expense.Description, "amount", expense.Amount.String())
		return nil
	})
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
