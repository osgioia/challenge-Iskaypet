-- Crear la tabla 'clients' (si aún no existe)
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    birth_day DATE NOT NULL,
    age INTEGER NOT NULL,
    telephone TEXT
);

-- DDL: Membuat tabel users dan groups
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE  IF NOT EXISTS groups (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE  IF NOT EXISTS user_groups (
    user_id INT,
    group_id INT,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- DML: Mengisi data awal
INSERT INTO users (username, email, password, first_name, last_name) VALUES 
('john_doe', 'john@example.com', 'hashed_password', 'John', 'Doe'),
('jane_doe', 'jane@example.com', 'hashed_password', 'Jane', 'Doe');

INSERT INTO groups (name, description) VALUES 
('Admin', 'Administrative Group'),
('User', 'Regular User Group');

INSERT INTO user_groups (user_id, group_id) VALUES 
(1, 1), -- John Doe assigned to Admin
(2, 2); -- Jane Doe assigned to User

-- Insertar datos iniciales de clientes
INSERT INTO clients (name, last_name, email, birth_day) VALUES
('John', 'Doe', 'johndoe@example.com', '1985-05-15'),
('Jane', 'Smith', 'janesmith@example.com', '1990-06-20'),
('Alice', 'Johnson', 'alicej@example.com', '1978-12-05');