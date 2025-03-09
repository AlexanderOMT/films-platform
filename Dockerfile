# Use the go version from go.mod to build the go image 
# Step 1: create a go envioroment builder
FROM golang:1.23.6

# Create and set the working dictory
RUN mkdir -p /app
WORKDIR /app

# Copy the project to working dictory
COPY . /app

# Need for initialization of database
RUN apt-get update && apt-get install -y jq postgresql-client && rm -rf /var/lib/apt/lists/*

# Download and clean depedencies
RUN go mod tidy
# Build the go binary
RUN go build -o ./out/films-api ./cmd/api/main.go

# TODO Step 2: craete a cleaned builder with go binary but clean go enviroment and a light linux image, in order to get a lighter container
# We will not be able to use go tools inside the container as we will only keep the go API binary

# TODO: consider if this is a good idea: see https://www.docker.com/blog/developing-go-apps-docker/

# Run golang API binary
CMD ["./out/films-api"]
