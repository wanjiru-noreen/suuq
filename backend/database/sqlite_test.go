package database

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	tables := []string{
		"users",
		"businesses",
		"products",
		"stock_updates",
		"debtors",
		"creditors",
	}

	for _, table := range tables {
		var name string

		err := db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?",
			table,
		).Scan(&name)

		if err != nil {
			t.Errorf("table %s was not created: %v", table, err)
		}
	}
}

func TestForeignKeys(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var enabled int

	err = db.QueryRow("PRAGMA foreign_keys").Scan(&enabled)
	if err != nil {
		t.Fatal(err)
	}

	if enabled != 1 {
		t.Fatal("foreign keys are not enabled")
	}
}

func TestProductConstraints(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO products (business_id, name, quantity)
		VALUES (999, 'Test Product', 10)
	`)
	if err == nil {
		t.Fatal("expected foreign key constraint to reject invalid business_id")
	}
}

func TestProductQuantityAndNameConstraints(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES ('Test User', 'test@example.com', 'hash')
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO businesses (user_id, name)
		VALUES (1, 'Test Business')
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO products (business_id, name, quantity)
		VALUES (1, 'Milk', 10)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Negative quantity should be rejected.
	_, err = db.Exec(`
		INSERT INTO products (business_id, name, quantity)
		VALUES (1, 'Bread', -5)
	`)
	if err == nil {
		t.Fatal("expected negative quantity to be rejected")
	}

	// Duplicate product name within the same business should be rejected.
	_, err = db.Exec(`
		INSERT INTO products (business_id, name, quantity)
		VALUES (1, 'Milk', 20)
	`)
	if err == nil {
		t.Fatal("expected duplicate product name to be rejected")
	}
}

func TestDebtorConstraints(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES ('Debtor User', 'debtor@example.com', 'hash')
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO businesses (user_id, name)
		VALUES (1, 'Debtor Business')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// A valid debtor should work.
	_, err = db.Exec(`
		INSERT INTO debtors (business_id, name, phone, amount_owed)
		VALUES (1, 'John', '0712345678', 500)
	`)
	if err != nil {
		t.Fatalf("valid debtor was rejected: %v", err)
	}

	// A negative amount should be rejected.
	_, err = db.Exec(`
		INSERT INTO debtors (business_id, name, amount_owed)
		VALUES (1, 'Jane', -100)
	`)
	if err == nil {
		t.Fatal("expected negative debtor amount to be rejected")
	}

	// A nonexistent business should be rejected.
	_, err = db.Exec(`
		INSERT INTO debtors (business_id, name, amount_owed)
		VALUES (999, 'Invalid', 100)
	`)
	if err == nil {
		t.Fatal("expected invalid business_id to be rejected")
	}
}

func TestCreditorConstraints(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES ('Creditor User', 'creditor@example.com', 'hash')
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO businesses (user_id, name)
		VALUES (1, 'Creditor Business')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// A valid creditor should work.
	_, err = db.Exec(`
		INSERT INTO creditors (business_id, name, phone, amount_owed)
		VALUES (1, 'Supplier Ltd', '0798765432', 1000)
	`)
	if err != nil {
		t.Fatalf("valid creditor was rejected: %v", err)
	}

	// A negative amount should be rejected.
	_, err = db.Exec(`
		INSERT INTO creditors (business_id, name, amount_owed)
		VALUES (1, 'Another Supplier', -200)
	`)
	if err == nil {
		t.Fatal("expected negative creditor amount to be rejected")
	}

	// A nonexistent business should be rejected.
	_, err = db.Exec(`
		INSERT INTO creditors (business_id, name, amount_owed)
		VALUES (999, 'Invalid Supplier', 200)
	`)
	if err == nil {
		t.Fatal("expected invalid business_id to be rejected")
	}
}

func TestTransactionItemConstraints(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	// Create a user.
	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES ('Transaction User', 'transaction@example.com', 'hash')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create a business.
	_, err = db.Exec(`
		INSERT INTO businesses (user_id, name)
		VALUES (1, 'Transaction Business')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create a product.
	_, err = db.Exec(`
		INSERT INTO products (business_id, name, quantity)
		VALUES (1, 'Milk', 20)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create a transaction.
	_, err = db.Exec(`
		INSERT INTO transactions (business_id)
		VALUES (1)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// A valid transaction item should work.
	_, err = db.Exec(`
		INSERT INTO transaction_items (transaction_id, product_id, quantity)
		VALUES (1, 1, 2)
	`)
	if err != nil {
		t.Fatalf("valid transaction item was rejected: %v", err)
	}

	// Zero quantity should be rejected.
	_, err = db.Exec(`
		INSERT INTO transaction_items (transaction_id, product_id, quantity)
		VALUES (1, 1, 0)
	`)
	if err == nil {
		t.Fatal("expected zero quantity to be rejected")
	}

	// A nonexistent transaction should be rejected.
	_, err = db.Exec(`
		INSERT INTO transaction_items (transaction_id, product_id, quantity)
		VALUES (999, 1, 2)
	`)
	if err == nil {
		t.Fatal("expected invalid transaction_id to be rejected")
	}

	// A nonexistent product should be rejected.
	_, err = db.Exec(`
		INSERT INTO transaction_items (transaction_id, product_id, quantity)
		VALUES (1, 999, 2)
	`)
	if err == nil {
		t.Fatal("expected invalid product_id to be rejected")
	}
}
