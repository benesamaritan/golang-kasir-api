CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    price INT NOT NULL,
    stock INT NOT NULL,
    category_id INT,
    active BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    subtotal INT NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- Insert example data into categories
INSERT INTO categories (name, description) VALUES
('Minuman', 'Berbagai jenis minuman dingin dan panas'),
('Makanan', 'Makanan berat untuk makan siang atau malam'),
('Snack', 'Camilan ringan untuk menemani waktu luang');

-- Insert example data into products (assuming category IDs are 1, 2, 3 respectively)
INSERT INTO products (name, price, stock, category_id, active) VALUES
('Es Teh', 5000, 100, 1, TRUE),
('Es Jeruk', 6000, 80, 1, TRUE),
('Nasi Goreng', 15000, 50, 2, TRUE),
('Mie Ayam', 12000, 60, 2, TRUE),
('Keripik Kentang', 8000, 120, 3, FALSE); -- One inactive product for testing

-- Insert example data into transactions and transaction_details
-- Transaction 1
INSERT INTO transactions (total_amount) VALUES (25000);
INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES
((SELECT id FROM transactions ORDER BY id DESC LIMIT 1), (SELECT id FROM products WHERE name = 'Es Teh'), 2, 10000),
((SELECT id FROM transactions ORDER BY id DESC LIMIT 1), (SELECT id FROM products WHERE name = 'Nasi Goreng'), 1, 15000);

-- Transaction 2
INSERT INTO transactions (total_amount) VALUES (22000);
INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES
((SELECT id FROM transactions ORDER BY id DESC LIMIT 1), (SELECT id FROM products WHERE name = 'Es Jeruk'), 1, 6000),
((SELECT id FROM transactions ORDER BY id DESC LIMIT 1), (SELECT id FROM products WHERE name = 'Keripik Kentang'), 2, 16000);