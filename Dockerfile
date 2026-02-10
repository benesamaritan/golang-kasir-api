FROM jetpackio/devbox:latest

# WORKDIR sets the destination folder INSIDE the container.
# Your code from the project root will be copied here.
WORKDIR /app

# Copy devbox configuration first for better caching
COPY devbox.json devbox.lock ./

# Install the environment
RUN devbox run -- echo "Environment initialized"

# Copy all files from your project root into /app inside the container
COPY . .

# Start the application
CMD ["devbox", "run", "kasir-up"]
