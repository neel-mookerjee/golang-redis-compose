### Builder Image ###

FROM golang:latest as builder

LABEL maintainer="Arghanil <arghanil@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Test
RUN go test -v

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


### Actual Image ###
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY . .

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 