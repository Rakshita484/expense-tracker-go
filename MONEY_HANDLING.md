# 💵 Money Handling Approach

This document explains the strategies and technologies used to ensure 100% financial accuracy in the Expense Tracker API.

## 1. The Core Problem: Floating Point Math

In most programming languages, `float64` or `double` types use binary floating-point representation (IEEE 754). This is unsuitable for money because it cannot precisely represent many decimal fractions.

**Example of the problem:**
```go
0.1 + 0.2 // Results in 0.30000000000000004, not 0.3
```

## 2. Our Solution: Arbitrary-Precision Decimals

We use the [`shopspring/decimal`](https://github.com/shopspring/decimal) library for all monetary calculations in Go.

### Why this works:
- It stores numbers as integers internally (coefficient * 10^exponent).
- It provides precise control over rounding and truncation.
- It prevents precision accumulation errors during complex splits.

## 3. Database Storage

We use the `NUMERIC(15,2)` type in the database (SQLite/PostgreSQL).

- **Precision (15)**: Total number of digits stored.
- **Scale (2)**: Number of digits after the decimal point.

Unlike `FLOAT` or `REAL`, `NUMERIC` stores values exactly as provided, matching the precision of our Go code.

## 4. Equal Split Algorithm

Splitting $100 among 3 people is mathematically $33.3333...$ which cannot be represented in currency. We solve this using a **Truncation + Remainder Distribution** strategy:

### The Steps:
1. **Divide**: $100.00 / 3 = 33.3333...$
2. **Truncate**: Floor the value to 2 decimal places → **$33.33** per person.
3. **Calculate Remainder**: $100.00 - (33.33 \times 3) = **$0.01**.
4. **Distribute**: Give the $0.01 remainder to the first member.

### Final Result:
- Member 1: **$33.34**
- Member 2: **$33.33**
- Member 3: **$33.33**
- **Total**: $100.00 (Exactly matches original amount)

## 5. API Transmission

Monetary values are transmitted as **Strings** in JSON responses (e.g., `"120.00"`).

**Why?**
JSON numbers are typically parsed as floating-point by browsers and other clients. By using strings, we force the client to handle the number explicitly, preserving the precision we calculated on the server.
