CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance INT NOT NULL DEFAULT 1000
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    product_type VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user INT REFERENCES users(id),
    to_user INT REFERENCES users(id),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE inventory (
    user_id INT REFERENCES users(id),
    item_id INT REFERENCES goods(id),
    quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (user_id, item_id)
);

INSERT INTO goods (product_type, price)
VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);
    