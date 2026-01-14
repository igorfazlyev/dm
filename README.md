# dm
# 1. Clone or create project
mkdir dental-marketplace && cd dental-marketplace

# 2. Initialize go module
go mod init github.com/yourorg/dental-marketplace

# 3. Copy all files above into respective directories

# 4. Install dependencies
go mod tidy

# 5. Start PostgreSQL
docker-compose up -d postgres

# 6. Run the API (will auto-migrate)
make run

# Or run directly
go run cmd/api/main.go
