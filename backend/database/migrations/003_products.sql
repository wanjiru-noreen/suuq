CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    business_id INTEGER NOT NULL,

    name TEXT NOT NULL,

    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,

    UNIQUE (business_id, name)
);