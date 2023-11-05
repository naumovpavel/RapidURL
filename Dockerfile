FROM golang:1.21

# Set destination for COPY
WORKDIR /RapidURL
ENV JWT_SECRET_KEY=secret
ENV POSTGRES_PASSWORD=admin
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

CMD migrate -database postgres://postgres:admin@localhost:5432/RapidURL?sslmode=disable -path deployments/migrations up

# Run
CMD CGO_ENABLED=0 GOOS=linux go run ./cmd/RapidURL/main.go -config="./config/local.yaml"