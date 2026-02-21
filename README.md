# 💰 Expense Tracker & Bill Splitting API

A production-ready REST API built in Go for expense tracking and bill splitting, inspired by Splitwise. Features a greedy settlement algorithm that minimizes the number of transactions needed to settle debts.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clearly separated layers:

```
expense-tracker/
├── main.go                 # Entry point: wires dependencies, starts server
├── config/config.go        # Environment & database configuration
├── models/models.go        # Domain models & DTOs
├── repository/repository.go # Data access layer (GORM + PostgreSQL)
├── services/services.go    # Business logic & settlement algorithm
├── services/settlement_test.go # Unit tests for settlement
├── handlers/handlers.go    # HTTP handlers (Gin)
├── routes/routes.go        # Route definitions
├── seed/seed.go            # Database seeder
├── postman/                # Postman collection
├── .env.example            # Environment template
└── README.md
```

**Dependency flow**: `main.go` → `handlers` → `services` → `repository` → `database`

Each layer only depends on the layer directly below it — never the reverse.

---

## 🚀 Setup Instructions

### Prerequisites

- **Go** 1.21+ installed
- **PostgreSQL** running locally or remotely

### 1. Clone & Configure

```bash
cd expense-tracker
cp .env.example .env
# Edit .env with your PostgreSQL credentials
```

### 2. Create the Database

```sql
CREATE DATABASE expense_tracker;
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Server

```bash
go run main.go
```

The server starts on `http://localhost:8080`. Tables are auto-migrated on startup.

### 5. Seed Sample Data (Optional)

```bash
go run seed/seed.go
```

### 6. Run Tests

```bash
go test ./services/ -v
```

---

## 📡 API Documentation

### Base URL: `http://localhost:8080/api`

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/users` | Create a new user |
| `GET` | `/users` | List all users |
| `POST` | `/groups` | Create a new group |
| `POST` | `/groups/:id/members` | Add a user to a group |
| `POST` | `/groups/:id/expenses` | Add an expense to a group |
| `GET` | `/groups/:id/balances` | Get net balances for a group |
| `GET` | `/groups/:id/settlements` | Get minimum settlement transactions |

### Request/Response Examples

#### Create User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@example.com"}'
```
Response (201):
```json
{
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "Alice Johnson",
    "email": "alice@example.com"
  }
}
```

#### Create Group
```bash
curl -X POST http://localhost:8080/api/groups \
  -H "Content-Type: application/json" \
  -d '{"name": "Weekend Trip"}'
```

#### Add User to Group
```bash
curl -X POST http://localhost:8080/api/groups/1/members \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'
```

#### Add Expense
```bash
curl -X POST http://localhost:8080/api/groups/1/expenses \
  -H "Content-Type: application/json" \
  -d '{"description": "Dinner", "amount": "120.00", "paid_by_id": 1}'
```

#### Get Balances
```bash
curl http://localhost:8080/api/groups/1/balances
```

#### Get Settlements
```bash
curl http://localhost:8080/api/groups/1/settlements
```

### Error Responses

All errors return structured JSON:
```json
{
  "error": "Failed to create expense",
  "details": "payer (user 99) is not a member of group 1"
}
```

| Status Code | Meaning |
|------------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## 💵 Money Handling

### Why Not `float64`?

Floating-point arithmetic causes precision errors in financial calculations:

```
// float64 (WRONG):
0.1 + 0.2 = 0.30000000000000004

// decimal (CORRECT):
0.1 + 0.2 = 0.3
```

### Our Approach

- **Go code**: All monetary values use [`shopspring/decimal`](https://github.com/shopspring/decimal), which provides arbitrary-precision decimal arithmetic.
- **Database**: Money columns use PostgreSQL's `NUMERIC(15,2)` type, which stores exact decimal values.
- **API**: Amounts are transmitted as strings (`"120.00"`) to avoid JSON floating-point issues.
- **Splitting**: Equal splits are calculated using truncation with remainder distribution:
  ```
  $100.00 ÷ 3 users:
  Base:      $33.33 per person (truncated)
  Total:     $33.33 × 3 = $99.99
  Remainder: $100.00 - $99.99 = $0.01
  Result:    User 1 gets $33.34, Users 2-3 get $33.33
  Sum:       $33.34 + $33.33 + $33.33 = $100.00 ✓
  ```

**Zero floating-point errors at every layer.**

---

## 🧮 Settlement Algorithm

### Greedy Debt Simplification

The settlement algorithm minimizes the number of transactions needed to clear all debts:

```
Input:  User balances (positive = owed money, negative = owes money)
Output: Minimum set of (from, to, amount) transactions
```

### Algorithm Steps

1. **Separate** users into creditors (+) and debtors (−)
2. **Sort** creditors descending, debtors by absolute value descending
3. **Match** largest creditor with largest debtor:
   - Transfer `min(credit, |debt|)`
   - Reduce both balances
   - Move to next when one reaches zero
4. **Repeat** until all balances are settled

### Time Complexity Analysis

| Operation | Complexity |
|-----------|-----------|
| Separate creditors/debtors | O(n) |
| Sort both lists | O(n log n) |
| Match and settle | O(n) |
| **Total** | **O(n log n)** |

Where `n` = number of users with non-zero balances.

**Space Complexity**: O(n) for storing creditor and debtor lists.

### Example Walkthrough

```
Users: Alice (+$60), Bob (+$30), Charlie (-$50), Dave (-$40)

Step 1: Charlie pays Alice min($50, $60) = $50
        → Alice: +$10, Charlie: $0 (settled)

Step 2: Dave pays Alice min($40, $10) = $10
        → Alice: $0 (settled), Dave: -$30

Step 3: Dave pays Bob min($30, $30) = $30
        → Bob: $0 (settled), Dave: $0 (settled)

Result: 3 transactions (down from potential 6 pairwise)
```

---

## 🔒 Concurrency Handling

### How Race Conditions Are Prevented

#### 1. Database Transactions (ACID)

Expense creation uses GORM's `Transaction()` method, which wraps the expense insert and all split inserts in a single SQL transaction:

```go
func (r *Repository) CreateExpenseWithSplits(expense *models.Expense, splits []models.ExpenseSplit) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Both operations succeed or both fail — no partial writes.
        tx.Create(expense)
        tx.Create(&splits)
        return nil
    })
}
```

**Guarantees:**
- **Atomicity**: Either the expense AND all splits are saved, or nothing is.
- **Consistency**: The database is never in a state with an expense but no splits.

#### 2. PostgreSQL MVCC

PostgreSQL uses Multi-Version Concurrency Control (MVCC) at the `READ COMMITTED` isolation level (default). This means:

- Concurrent expense creations in the same group don't block each other.
- Each transaction sees a consistent snapshot of the data.
- Balance calculations always use committed data — never dirty reads.

#### 3. Stateless API Design

The Gin handlers are stateless — no shared mutable state between requests. Each request gets its own:
- Gin context
- GORM session (via connection pool)
- Local variables

This eliminates in-memory race conditions entirely. All shared state lives in PostgreSQL, which handles concurrency natively.

---

## 📋 Example Scenario

### Setup

Three friends — Alice, Bob, Charlie — go on a weekend trip.

| Expense | Amount | Paid By | Per Person |
|---------|--------|---------|-----------|
| Dinner | $120.00 | Alice | $40.00 |
| Gas | $80.00 | Bob | $26.67 |
| Snacks | $45.99 | Charlie | $15.33 |

### Balance Calculation

| User | Total Paid | Total Owed | Net Balance |
|------|-----------|------------|-------------|
| Alice | $120.00 | $82.00 | **+$38.00** |
| Bob | $80.00 | $82.00 | **−$2.00** |
| Charlie | $45.99 | $81.99 | **−$36.00** |

### Settlements

```
1. Charlie pays Alice $36.00  (Charlie settled)
2. Bob pays Alice $2.00       (Both settled)

Total: 2 transactions ✓
Net sum: $0.00 ✓
```

---

## 📬 Postman Collection

Import `postman/Expense_Tracker.postman_collection.json` into Postman. The collection includes all API endpoints with example request bodies. Set the `baseUrl` variable to your server address (default: `http://localhost:8080`).

---

## 🛠️ Tech Stack

| Technology | Purpose |
|-----------|---------|
| Go 1.21+ | Language |
| Gin | HTTP routing framework |
| GORM | ORM for PostgreSQL |
| shopspring/decimal | Precise decimal arithmetic |
| godotenv | Environment variable loading |
| log/slog | Structured logging |
| PostgreSQL | Database |

---

## 📝 License

MIT
