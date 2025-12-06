CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role_id INT NOT NULL REFERENCES roles(id),
    "created_by" varchar(200) NOT NULL,
    "created_date" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_by" varchar(200) NOT NULL,
    "modified_date" timestamptz,
    is_deleted boolean NOT NULL DEFAULT false
);



CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL
);


CREATE TYPE status_type AS ENUM (
  'editable',
  'confirm',
  'shipping',
  'completed',
  'cancel'
);


CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    total NUMERIC(12, 2) NOT NULL,
    status status_type NOT NULL,
    reason TEXT,
    "created_by" varchar(200) NOT NULL,
    "created_date" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_by" varchar(200) NOT NULL,
    "modified_date" timestamptz
);


CREATE TABLE IF NOT EXISTS order_details (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    sub_total NUMERIC(12, 2) NOT NULL
);