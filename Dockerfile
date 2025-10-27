
FROM golang:1.23.0-alpine AS builer
# # builer stage
# alpine là phiên bản Linux nhẹ, tối giản
# AS builder đặt tên cho stage này là "builder" để tham chiếu sau

RUN apk add --no-cache git
# Cài đặt git bằng package manager apk của Alpine
# --no-cache không lưu cache package để giảm kích thước image

WORKDIR /app
# Tạo và chuyển đến thư mục làm việc /app trong container

RUN go.mod go.sum ./
# Copy file go.mod và go.sum từ host vào thư mục hiện tại (/app) trong container

RUN go mod download
# Tải về tất cả dependencies được khai báo trong go.mod

COPY . .
# Copy toàn bộ source code từ host vào container

RUN go build -o main .
# Biên dịch ứng dụng Go thành file binary tên main
# Dấu . chỉ định build package trong thư mục hiện tại


# Giai đoạn 2: Production Stage

FROM alpine:latest
# Sử dụng image alpine gốc (nhẹ) cho production, không chứa Go toolchain

RUN apk add --no-cache ca-certificates
# Cài đặt CA certificates để hỗ trợ kết nối SSL/TLS

WORKDIR /app
# Tạo thư mục làm việc /app trong container producti

COPY  --from=builer /app/main .
# Chỉ copy file binary main từ stage "builder" sang stage production
# Đây là kỹ thuật multi-stage build giúp image cuối cùng nhỏ hơn

COPY app.env .env
# Copy file cấu hình môi trường từ host vào container
EXPOSE 8080
