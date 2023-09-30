# Gunakan image Golang sebagai base image
FROM golang:latest

# Set working directory di dalam container
WORKDIR /app

# Copy file Go yang diperlukan ke dalam container
COPY go.mod .
COPY go.sum .

# Download dependencies Go yang diperlukan
RUN go mod download

# Copy seluruh proyek Anda ke dalam container
COPY . .

# Build aplikasi Go
RUN go build -o main .

# Eksekusi aplikasi saat container dijalankan
CMD ["./main"]
