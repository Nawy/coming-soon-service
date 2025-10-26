# Coming Soon Service

A lightweight Go service for collecting email addresses for "coming soon" landing pages. Built with Gin framework and designed to run in Docker containers.

## Features

- ✅ REST API for email collection
- ✅ Email validation and deduplication
- ✅ Thread-safe operations
- ✅ Protected endpoint for retrieving emails
- ✅ Persistent storage with external volume support
- ✅ Docker-ready with multi-stage builds

## Prerequisites

- **For Local Development:**
  - Go 1.25.1 or higher
  - Git

- **For Docker:**
  - Docker installed and running
  - Docker Compose (optional, for easier management)

## Project Structure

```
comint-soon-service/
├── main.go           # Main application code
├── go.mod            # Go module dependencies
├── go.sum            # Go module checksums
├── Dockerfile        # Docker build configuration
└── README.md         # This file
```

## Configuration

The service uses environment variables for configuration:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SECRET_TOKEN` | Yes | - | Secret token for protecting the GET endpoint |
| `EMAIL_FILE_PATH` | No | `emails.txt` | Path where emails are stored |

## API Endpoints

### POST /coming-soon
Submit a new email address (public endpoint).

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response (201 Created):**
```json
{
  "message": "Email registered successfully."
}
```

**Response (409 Conflict):**
```json
{
  "message": "Email already registered."
}
```

### GET /coming-soon
Retrieve all registered emails (protected endpoint).

**Headers:**
```
X-Secret-Token: your-secret-token
```

**Response (200 OK):**
```json
[
  "user1@example.com",
  "user2@example.com"
]
```

## Building the Application

### Option 1: Build Docker Image

Build the Docker image for the coming soon service:

```bash
docker build -t coming-soon-service:latest .
```

This creates a lightweight Docker image (~15MB) using a multi-stage build process.

### Option 2: Build Locally (Without Docker)

```bash
# Navigate to the project directory
cd comint-soon-service

# Download dependencies
go mod download

# Build the binary
go build -o coming-soon-service main.go
```

## Running the Application

### Option 1: Run with Docker (Recommended)

1. **Create a directory on your host machine for storing emails:**

```bash
# Windows (PowerShell)
mkdir C:\data\coming-soon

# Windows (Command Prompt)
mkdir C:\data\coming-soon

# Linux/Mac
mkdir -p /path/to/data/coming-soon
```

2. **Run the container:**

```bash
# Windows (PowerShell)
docker run -d `
  --name coming-soon-service `
  -p 8080:8080 `
  -e SECRET_TOKEN="your-super-secret-token-here" `
  -v C:\data\coming-soon:/app/data `
  coming-soon-service:latest

# Windows (Command Prompt)
docker run -d ^
  --name coming-soon-service ^
  -p 8080:8080 ^
  -e SECRET_TOKEN="your-super-secret-token-here" ^
  -v C:\data\coming-soon:/app/data ^
  coming-soon-service:latest

# Linux/Mac
docker run -d \
  --name coming-soon-service \
  -p 8080:8080 \
  -e SECRET_TOKEN="your-super-secret-token-here" \
  -v /path/to/data/coming-soon:/app/data \
  coming-soon-service:latest
```

**Important Notes:**
- Replace `your-super-secret-token-here` with a strong, random token
- The `-v` flag maps a local directory to `/app/data` in the container
- The `emails.txt` file will be created in your local directory (outside the container)
- Data persists even if the container is stopped or removed

3. **Check if the container is running:**

```bash
docker ps
```

4. **View logs:**

```bash
docker logs coming-soon-service
```

### Option 2: Run with Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  coming-soon:
    build: .
    container_name: coming-soon-service
    ports:
      - "8080:8080"
    environment:
      - SECRET_TOKEN=${SECRET_TOKEN:-your-super-secret-token-here}
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

Then run:

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### Option 3: Run Locally (Without Docker)

```bash
# Set environment variables
# Windows (PowerShell)
$env:SECRET_TOKEN="your-super-secret-token-here"
$env:EMAIL_FILE_PATH="emails.txt"

# Windows (Command Prompt)
set SECRET_TOKEN=your-super-secret-token-here
set EMAIL_FILE_PATH=emails.txt

# Linux/Mac
export SECRET_TOKEN="your-super-secret-token-here"
export EMAIL_FILE_PATH="emails.txt"

# Run the application
go run main.go

# Or run the built binary
./coming-soon-service
```

## Testing the API

### Submit an email (POST)

```bash
# Using curl
curl -X POST http://localhost:8080/coming-soon \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Using PowerShell
Invoke-RestMethod -Uri http://localhost:8080/coming-soon `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"test@example.com"}'
```

### Retrieve all emails (GET)

```bash
# Using curl
curl -X GET http://localhost:8080/coming-soon \
  -H "X-Secret-Token: your-super-secret-token-here"

# Using PowerShell
Invoke-RestMethod -Uri http://localhost:8080/coming-soon `
  -Method GET `
  -Headers @{"X-Secret-Token"="your-super-secret-token-here"}
```

## Managing the Docker Container

### Stop the container
```bash
docker stop coming-soon-service
```

### Start the container
```bash
docker start coming-soon-service
```

### Remove the container
```bash
docker rm -f coming-soon-service
```

### View container logs
```bash
docker logs -f coming-soon-service
```

### Access the container shell
```bash
docker exec -it coming-soon-service sh
```

## Data Persistence

The `emails.txt` file is stored **outside** the container using Docker volumes:

- **Container path:** `/app/data/emails.txt`
- **Host path:** The directory you specified with the `-v` flag

This ensures:
- ✅ Data persists when containers are stopped/removed
- ✅ Easy backup and access to email data
- ✅ Can be shared across multiple container instances

## Security Best Practices

1. **Never commit secrets:** Don't hardcode the `SECRET_TOKEN` in your code or Dockerfile
2. **Use strong tokens:** Generate a cryptographically secure random string
3. **Restrict access:** Use firewall rules to limit access to the service
4. **HTTPS in production:** Use a reverse proxy (nginx, Traefik) with SSL/TLS certificates
5. **Regular updates:** Keep dependencies and base images up to date

### Generating a Secure Token

```bash
# Linux/Mac
openssl rand -hex 32

# PowerShell
-join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | % {[char]$_})
```

## Troubleshooting

### Container won't start
- Check logs: `docker logs coming-soon-service`
- Verify `SECRET_TOKEN` is set
- Ensure port 8080 is not already in use

### Can't access the API
- Verify container is running: `docker ps`
- Check port mapping is correct: `-p 8080:8080`
- Test from inside container: `docker exec -it coming-soon-service wget -O- localhost:8080/coming-soon`

### Emails not persisting
- Verify volume mount is correct: `docker inspect coming-soon-service`
- Check directory permissions on the host
- Ensure the data directory exists before running the container

### Permission denied errors
- The container runs as non-root user (UID 1000)
- Ensure the host directory has appropriate permissions

## Production Deployment

For production deployments, consider:

1. **Use orchestration:** Kubernetes, Docker Swarm, or similar
2. **Add monitoring:** Prometheus, Grafana, or cloud monitoring services
3. **Set up logging:** Centralized logging with ELK stack or similar
4. **Use secrets management:** Kubernetes secrets, AWS Secrets Manager, etc.
5. **Implement rate limiting:** Protect against abuse
6. **Add CORS:** Configure appropriate CORS headers for frontend integration
7. **Use HTTPS:** Always use TLS/SSL in production

## License

This project is provided as-is for educational and commercial purposes.

## Support

For issues, questions, or contributions, please contact the development team.