# Loan Management API
This API is used to manage customer data, loans, payment schedules, and payments within a loan management system. Each endpoint includes examples of requests and successful responses. Note that responses may vary in the event of an error.

## Running the Application Locally
To run this application on your local machine, follow these steps:
1. Clone the repository:

```sh
git clone git@github.com:jacky-htg/loans.git
```

2. Execute the SQL queries in the file ./migrations/01_ddl_loans.sql to create the required database tables.

3. Copy the .env.example file to a new .env file and fill in the necessary configuration details.

4. Install the dependencies with the command:

```sh
go mod tidy
```

5. Run the application with the command:

```sh
go run main.go
```

## API List
### 1. Create Customer
- Endpoint: POST /customers
- Deskripsi: Adds a new customer to the system.
- Request:
```bash
curl --location 'localhost:8080/customers' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-secret-token' \
--data-raw '{
    "name": "Jacky Chan",
    "email": "jacky.chan@gmail.com",
    "phone": "08123456789"
}'
```
- Response:
```json
{
    "customer_id": 1,
    "message": "Customer created successfully"
}
```

### 2. Get Customer by ID
- Endpoint: GET /customers/:customerId
- Deskripsi: Retrieves customer information by ID.
- Request:
```bash
curl --location 'localhost:8080/customers/1' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "customer_id": 1,
    "name": "Jacky Chan",
    "email": "jacky.chan@gmail.com",
    "phone": "08123456789"
}
```

### 3. Create Loan
- Endpoint: POST /loans
- Deskripsi: Creates a new loan for a customer.
- Request:
```bash
curl --location 'localhost:8080/loans' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-secret-token' \
--data '{
    "customer_id": 1,
    "principal_amount": 5000000,
    "term_in_weeks": 50,
    "interest_rate": 10
}'
```
- Response:
```json
{
    "loan_id": 1,
    "message": "Loan created successfully"
}
```

### 4. Get Loan by ID
- Endpoint: GET /loans/:loanId
- Deskripsi: Retrieves loan information by ID.
- Request:
```bash
curl --location 'localhost:8080/loans/1' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "loan_id": 1,
    "customer_id": 1,
    "principal_amount": 5000000,
    "term_in_weeks": 50,
    "interest_rate": 10
}
```

### 5. Get Payment Schedule
- Endpoint: GET /loans/:loanId/schedule
- Deskripsi: Retrieves the payment schedule for a specific loan.
- Request:
```bash
curl --location 'localhost:8080/loans/1/schedule' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
[
    {"installment_number":1,"due_date":"2024-08-20","amount":110000,"status":"Unpaid"},
    {"installment_number":2,"due_date":"2024-08-27","amount":110000,"status":"Unpaid"},
    {"installment_number":3,"due_date":"2024-09-03","amount":110000,"status":"Unpaid"},
    ...
    {"installment_number":50,"due_date":"2025-07-29","amount":110000,"status":"Unpaid"}
]
```

### 6. Make Payment
- Endpoint: POST /loans/:loanId/payment
- Deskripsi: Makes a payment for a specific loan installment.
- Request:
```bash
curl --location 'localhost:8080/loans/1/payment' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-secret-token' \
--data '{
    "installment_number": 1,
    "date": "2024-08-13",
    "amount": 110000 
}'
```
- Response:
```json
{
    "payment_id": 1,
    "message": "Payment successfully recorded"
}
```

### 7. Get Outstanding Amount
- Endpoint: GET /loans/:loanId/outstanding
- Deskripsi: Retrieves the outstanding amount for a specific loan.
- Request:
```bash
curl --location 'localhost:8080/loans/1/outstanding' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "loan_id": 1,
    "outstanding_amount": 4900000
}
```

### 8. Check Delinquent Status
- Endpoint: GET /loans/:loanId/delinquent
- Deskripsi: Checks whether the borrower is delinquent.
- Request:
```bash
curl --location 'localhost:8080/loans/1/delinquent' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "is_delinquent": false,
    "loan_id": "1"
}
```

### Notes
The responses shown above are examples of successful responses. In case of an error, the responses will vary depending on the type of error, such as validation errors, authentication errors, or server errors.