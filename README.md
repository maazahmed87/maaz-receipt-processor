# Receipt Processor Challenge

A simple REST API service that processes receipts and calculates points based on specific rules.

## Technologies Used

- Go
- Gin
- Docker
- In-memory Storage

## Getting Started

### Prerequisites

- Docker
- Go 1.21+

### Running with Docker

1. Build and start the container:
```bash
docker-compose up --build
```

### Running Locally

```bash
go mod download
go run cmd/api/main.go
```
The service runs at `http://localhost:8080`

## API Endpoints

### Process Receipt
- **POST** `/receipts/process`
- Processes a receipt and returns an ID for point retrieval
- Request Body:
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
- Response:
```json
{
  "id": "UUID"
}
```

### Get Points
- **GET** `/receipts/{id}/points`
- Retrieves the points awarded for a receipt
- Response:
```json
{
  "points": 28
}
```

### Health Check
- **GET** `/receipts/health`
- Checks if the service is running
- Response:
```json
{
  "status": "ok",
  "time": "2025-02-03T20:04:56.837656305Z"
}
```

## Project Structure

```
.
├── cmd/api/          # Entry point
├── internal/
│   ├── api/         # HTTP handlers
│   ├── domain/      # Core business logic
│   └── storage/     # Data persistence
├── pkg/             # Shared utilities
└── docker files     # Container configuration

```

## Implementation Details

### Design Principles & Architecture

- **Clean Architecture**
  - Clear separation of concerns with layered architecture
  - Domain logic isolated from external dependencies
  - Dependency injection for better testability
  - Interface-driven design enabling easy storage implementation swapping in future

- **SOLID Principles**
  - Single Responsibility: Each package has a single focused purpose
  - Open/Closed: New receipt rules can be added without modifying existing code
  - Interface Segregation: Small, focused interfaces (eg: ReceiptStorage)
  - Dependency Inversion: Business logic depends on abstractions

### Key Features

- Thread-safe receipt processing
- Comprehensive input validation
- Clear error handling
- Containerized deployment
- Efficient point calculation

## Whats next

### 1. Testing & CI/CD
   - Unit and integration tests
   - Automated deployments
   - Code quality checks

### 2. Persistence & Performance
   - PostgreSQL integration
   - Redis caching
   - Performance monitoring

### 3. Production Features
   - Authentication
   - Rate limiting
   - Structured logging
   - Metrics collection

### 4. Monitoring Setup
   - Prometheus metrics
   - Grafana dashboards
   - Basic alerting

Each improvement can be implemented incrementally based on business needs and deployment priorities.

