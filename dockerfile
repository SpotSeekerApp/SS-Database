# Start from golang base image
FROM golang:latest

WORKDIR /app

COPY . .
# Install the package
RUN go mod tidy

# Build the Go app
RUN cd ./src && go build -o main && cd ..

EXPOSE 8080
CMD [ "./src/main" ]