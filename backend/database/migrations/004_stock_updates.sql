CREATE TABLE IF NOT EXISTS stock_updates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    product_id INTEGER NOT NULL,

    quantity_change INTEGER NOT NULL,

    update_type TEXT NOT NULL CHECK (
        update_type IN ('addition', 'sale', 'adjustment')
    ),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);