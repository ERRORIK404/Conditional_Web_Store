FROM golang:1.25rc2-alpine AS builder

RUN apk add --no-cache git

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o marketplace ./cmd/server

FROM alpine:latest

RUN apk add --no-cache postgresql-client

COPY --from=builder /app/marketplace /app/marketplace
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /go/bin/migrate /app/migrate

EXPOSE ${API_PORT}

CMD ["sh", "-c", "while ! PGPASSWORD=${POSTGRES_PASSWORD} pg_isready -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do sleep 2; done; \
     ./migrate -path migrations -database \"$DB_URL\" up; \
     ./marketplace"]