CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance INT NOT NULL DEFAULT 1000 CHECK (balance >= 0)
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL CHECK (price > 0)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user INT REFERENCES users(id),
    to_user INT REFERENCES users(id),
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE purchases (
    user_id INT REFERENCES users(id),
    item_id INT REFERENCES goods(id),
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    purchase_price INT NOT NULL,
    purchased_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, item_id)
);