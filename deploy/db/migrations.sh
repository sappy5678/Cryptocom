echo "Creating database Cryptocom"
psql "postgres://postgres:password@sql:5432/postgres?sslmode=disable" -c "create database cryptocom;"

echo "Running migrations"
./migrate -source file://deploy/db/migrations -database "postgres://postgres:password@sql:5432/cryptocom?sslmode=disable" up

echo "Migrations completed"