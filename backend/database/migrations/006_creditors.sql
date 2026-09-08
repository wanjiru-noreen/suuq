CREATE TABLE IF NOT EXISTS creditors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    business_id INTEGER NOT NULL,

    name TEXT NOT NULL,

    phone TEXT,

    amount_owed REAL NOT NULL DEFAULT 0 CHECK (amount_owed >= 0),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
);