# 🧮 Settlement Scenarios & Optimization

This document demonstrates how the greedy settlement algorithm simplifies complex debts into the minimum number of transactions.

## Algorithm Overview

The goal is to move money from people who owe (debtors) to people who are owed (creditors) using the fewest steps. We follow a **Greedy Two-Pointer Strategy**:

1. Sort creditors descending (most owed first).
2. Sort debtors descending by absolute debt (most debt first).
3. Settle the largest possible amount between the top creditor and top debtor.
4. Update balances and repeat.

---

## Scenario 1: Simple Pairwise
*Three users where two owe one.*

- **Alice**: +$60 (Paid for dinner)
- **Bob**: -$30 (Owes share)
- **Charlie**: -$30 (Owes share)

**Settlements:**
1. **Bob** pays **Alice** $30.
2. **Charlie** pays **Alice** $30.

**Total: 2 Transactions** (Optimal)

---

## Scenario 2: The Chain Reaction
*Naive approach would require 3 payments, but optimization reduces it to 2.*

- **Alice**: +$50 (Lent to Bob)
- **Bob**: $0 (Owes Alice $50, but Charlie owes Bob $50)
- **Charlie**: -$50 (Owes Bob)

**Settlements:**
1. **Charlie** pays **Alice** $50.

**Total: 1 Transaction** (Optimized from 2)

---

## Scenario 3: Complex Group (The "Trip" Scenario)
*Five users with multiple overlapping expenses. Total expenses: $500.*

### Final Balances:
- **Alice**: +$120
- **Bob**: +$80
- **Charlie**: -$100
- **Dave**: -$70
- **Eve**: -$30

### Greedy Execution:
1. **Charlie** pays **Alice** $100 (Charlie settled, Alice needs $20 more)
2. **Dave** pays **Alice** $20 (Alice settled, Dave needs to pay $50 more)
3. **Dave** pays **Bob** $50 (Dave settled, Bob needs $30 more)
4. **Eve** pays **Bob** $30 (All settled)

**Total: 4 Transactions**
*Without optimization (pairwise), this could have taken up to 10 transactions.*

---

## Why this matters
1. **Reduces Friction**: Fewer bank transfers or digital payments between friends.
2. **Clarity**: Shows exactly who needs to pay whom to clear all debt.
3. **Exactness**: Using decimal math ensures that the sum of all payments exactly equals the total debt (Zero-sum).
