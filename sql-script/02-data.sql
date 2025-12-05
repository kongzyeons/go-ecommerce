INSERT INTO roles (name)
VALUES ('admin'),('user');

-- Insert admin user
INSERT INTO users (name,password,email,role_id ,created_by,modified_by)
VALUES ('admin','8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918','admin@example.com',1,'admin','admin');



-- Insert sample products
INSERT INTO products (name, description, price) VALUES
('AirPods Pro 3', 'รุ่นล่าสุดของ AirPods', 9900),
('MacBook Pro 14 M4', 'Laptop ระดับโปรของ Apple', 79900),
('Sony WH-1000XM7', 'Noise Cancelling Headphone', 14900),
('Samsung Galaxy S26 Ultra', 'Android Flagship 2025', 39900),
('Logitech MX Master 4S', 'Professional Work Mouse', 3990);