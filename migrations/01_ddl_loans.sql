CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE 
);

CREATE TABLE loans (
    loan_id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(customer_id),
    principal_amount NUMERIC(15, 2) NOT NULL,
    interest_rate NUMERIC(5, 2) NOT NULL DEFAULT 10.00,
    term_in_weeks INT NOT NULL,
    outstanding_amount NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES loans(loan_id),
    payment_amount NUMERIC(15, 2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    week INT NOT NULL,
    status VARCHAR(20) DEFAULT 'Paid'
);

CREATE TABLE payment_schedule (
    schedule_id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES loans(loan_id),
    week INT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'Unpaid'
);