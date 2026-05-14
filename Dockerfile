FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache git

# Copy everything first
COPY . .

# Get all deps, tidy, then build in one layer
RUN go get github.com/gin-gonic/gin@v1.9.1 \
    && go get github.com/joho/godotenv@v1.5.1 \
    && go get gorm.io/driver/mysql@v1.5.7 \
    && go get gorm.io/gorm@v1.25.12 \
    && go mod tidy \
    && go build -o main .

EXPOSE 7000

CMD ["./main"]
