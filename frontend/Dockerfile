# Use the official Bun image
FROM oven/bun:latest as builder

# Set working directory
WORKDIR /app

# Copy package.json and bun.lockb
COPY package.json bun.lockb ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy the rest of the code
COPY . .

# Build the app
RUN bun run build

# Start a new stage for a smaller final image
FROM oven/bun:1

WORKDIR /app

# Copy built assets from builder stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Copy src directory
COPY --from=builder /app/src ./src

# Expose the port Next.js runs on
EXPOSE 3000

# Run the app
CMD ["bun", "start"]