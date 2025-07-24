FROM golang:1.25rc2-alpine AS builder
#docker-compose down --volumes  --remove-orphans для отмены
RUN apk add --no-cache git

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN ls -l /go/bin/migrate

WORKDIR /under_app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./marketplace/cmd/server

RUN chmod +x /under_app/app

FROM alpine:latest

RUN apk add --no-cache postgresql-client

WORKDIR /Conditional_Web_Store

COPY --from=builder /under_app/app /Conditional_Web_Store/app
COPY --from=builder /under_app/migrations /Conditional_Web_Store/migrations
COPY --from=builder /go/bin/migrate /Conditional_Web_Store/migrate

EXPOSE ${API_PORT}

CMD ["sh", "-c", "export PGPASSWORD=\"${POSTGRES_PASSWORD}\" && \
    while ! pg_isready -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do sleep 2; done && \
    ./migrate -path migrations -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/marketplace?sslmode=disable\" up && \
    ./app"]