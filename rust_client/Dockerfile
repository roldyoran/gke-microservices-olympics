# Build stage
FROM rust:1.72.0-bullseye AS builder

# Install protobuf compiler
RUN apt-get update && apt-get install -y protobuf-compiler build-essential && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Copy the entire project
COPY . .

# Clean the build directory
RUN cargo clean --release

# Build the application
RUN cargo build --release

# Runtime stage
FROM debian:bullseye-slim

# Install necessary runtime libraries
RUN apt-get update && apt-get install -y libssl1.1 ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /usr/src/app/target/release/grpc-client .

# Set the startup command
CMD ["./grpc-client"]
