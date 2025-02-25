# Voucher Management API

A RESTful API service for managing vouchers, brands, and redemptions built with Go.

## Features

- Brand management
- Voucher creation and management
- Voucher redemption system
- Point-based transaction system
- PostgreSQL database
- RESTful API endpoints

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- [goose](https://github.com/pressly/goose) for database migrations

## Installation & Setup

1. Clone the repository:
```bash
git clone https://github.com/mahendraiqbal/be_otto.git
cd be_otto
```

2. Install dependencies:
```bash
make install
```

3. Start PostgreSQL using Docker:
```bash
make docker-up
```

4. Run database migrations (this will create tables and add a test customer with 1,000,000 points):
```bash
make migrate-up
```

## Running the Application

1. Start the server:
```bash
make run
```

The API will be available at `http://localhost:8080`

## Testing the API

### 1. Create a Brand

```bash
curl -X POST http://localhost:8080/brand \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Indomaret",
    "description": "Indonesian retail company"
  }'
```

Expected response:
```json
{
    "id": 1,
    "name": "Indomaret",
    "description": "Indonesian retail company",
    "created_at": "2024-03-14T...",
    "updated_at": "2024-03-14T..."
}
```

### 2. Create a Voucher

```bash
curl -X POST http://localhost:8080/voucher \
  -H "Content-Type: application/json" \
  -d '{
    "brand_id": 1,
    "code": "INDO50K",
    "name": "Indomaret 50K Voucher",
    "description": "Voucher worth 50,000 IDR",
    "point_cost": 50000,
    "stock": 100,
    "valid_until": "2024-12-31T23:59:59Z"
  }'
```

### 3. Get Voucher by ID

```bash
curl http://localhost:8080/voucher?id=1
```

### 4. Get All Vouchers by Brand

```bash
curl http://localhost:8080/voucher/brand?id=1
```

### 5. Create a Redemption (Buy Vouchers)

```bash
curl -X POST http://localhost:8080/transaction/redemption \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "items": [
      {
        "voucher_id": 1,
        "quantity": 2
      }
    ]
  }'
```

This will:
- Check if customer has enough points (1,000,000 points for test customer)
- Verify voucher stock availability
- Create redemption record
- Deduct points from customer
- Update voucher stock
- Return transaction details

### 6. Get Transaction Details

Using the ID from the redemption response:
```bash
curl http://localhost:8080/transaction/redemption\?transactionId\=1
```

## Database Schema

The application uses the following tables:
- `brands` - Stores brand information
- `vouchers` - Stores voucher information
- `customers` - Stores customer information and points
- `redemptions` - Stores redemption transactions
- `redemption_items` - Stores items in each redemption

## Development Commands

```bash
# Start PostgreSQL container
make docker-up

# Stop PostgreSQL container
make docker-down

# View PostgreSQL logs
make docker-logs

# Run database migrations
make migrate-up

# Rollback database migrations
make migrate-down

# Run the application
make run
```

## Database Connection

PostgreSQL instance details:
- Host: localhost
- Port: 5433
- Database: voucher_db
- Username: postgres
- Password: postgres

To connect to database:
```bash
docker exec -it voucher_db psql -U postgres -d voucher_db
```

## Troubleshooting

1. If port 5432 is in use, the application uses 5433 instead
2. Default test customer (ID: 1) has 1,000,000 points
3. Check logs using `make docker-logs` if issues occur
4. Ensure Docker is running before starting the application

## License

This project is licensed under the MIT License - see the LICENSE file for details
# be_otto
