CREATE TABLE IF NOT EXISTS transaction_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    transaction_id INTEGER NOT NULL,

    product_id INTEGER NOT NULL,

    quantity INTEGER NOT NULL CHECK (quantity > 0),

    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);