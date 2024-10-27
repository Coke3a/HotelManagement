CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role INT NOT NULL DEFAULT 0,
    rank VARCHAR(50),
    hire_date DATE DEFAULT CURRENT_DATE,
    last_login TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert the first user with hashed password
INSERT INTO users (username, password, role, rank, status)
VALUES ('admin', '$2a$10$k.HrLQjLb/ymtTaYG4gwfeXIgGzUrTYxUZF0.qxk0STgupDwg/Lwm', 1, 'admin', 'active');