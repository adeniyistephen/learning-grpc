# Generate proto files from .proto files
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto

# Create a pgsql instance with docker
docker pull postgres
docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres

# Check your db instance running on docker
docker exec -it some-postgres bash
su postgres
psql
\conninfo
\dt
SELECT * FROM users;