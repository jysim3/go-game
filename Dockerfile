FROM node:17-alpine AS front
WORKDIR /www
COPY www/ .
RUN npm install
RUN npm run build

FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main main.go 
COPY --from=0 /www/build ./www/build
CMD ["/app/main"]
