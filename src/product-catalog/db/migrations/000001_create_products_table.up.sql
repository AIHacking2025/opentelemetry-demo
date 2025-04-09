CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    picture VARCHAR(255) NOT NULL,
    price_units INTEGER NOT NULL,
    price_nanos INTEGER NOT NULL,
    price_currency_code VARCHAR(3) NOT NULL,
    categories TEXT[] NOT NULL
); 