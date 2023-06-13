FROM golang:alpine
USER root
WORKDIR /app

# 安装 Chrome 和 Chromedriver
RUN apk add --no-cache chromium chromium-chromedriver

# 定義 CHROMEDRIVER_PATH 的預設值
ARG CHROMEDRIVER_PATH=/usr/bin/chromedriver
# 將 CHROMEDRIVER_PATH 複製到環境變數中
ENV CHROMEDRIVER_PATH=${CHROMEDRIVER_PATH}

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .
CMD go run main.go --APP_ENV=prod