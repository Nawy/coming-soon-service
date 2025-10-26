# Quick Start Guide

Get the Coming Soon Service up and running in 5 minutes!

## Prerequisites

- Docker installed on your system
- Terminal or command prompt access

## Step 1: Clone or Download the Project

```bash
cd comint-soon-service
```

## Step 2: Build the Docker Image

```bash
docker build -t coming-soon-service:latest .
```

Expected output: `Successfully tagged coming-soon-service:latest`

## Step 3: Create Data Directory

This directory will store your emails outside the container.

**Windows (PowerShell or Command Prompt):**
```cmd
mkdir data
```

**Linux/Mac:**
```bash
mkdir -p data
```

## Step 4: Run the Container

### Windows (PowerShell):
```powershell
docker run -d `
  --name coming-soon-service `
  -p 8080:8080 `
  -e SECRET_TOKEN="my-secret-token-123" `
  -v ${PWD}/data:/app/data `
  coming-soon-service:latest
```

### Windows (Command Prompt):
```cmd
docker run -d ^
  --name coming-soon-service ^
  -p 8080:8080 ^
  -e SECRET_TOKEN="my-secret-token-123" ^
  -v %cd%/data:/app/data ^
  coming-soon-service:latest
```

### Linux/Mac:
```bash
docker run -d \
  --name coming-soon-service \
  -p 8080:8080 \
  -e SECRET_TOKEN="my-secret-token-123" \
  -v $(pwd)/data:/app/data \
  coming-soon-service:latest
```

## Step 5: Verify It's Running

Check the container status:
```bash
docker ps
```

View the logs:
```bash
docker logs coming-soon-service
```

You should see: `Service started, loaded 0 emails.`

## Step 6: Test the API

### Submit an Email (POST)

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Uri http://localhost:8080/coming-soon `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"test@example.com"}'
```

**Linux/Mac (curl):**
```bash
curl -X POST http://localhost:8080/coming-soon \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

**Expected Response:**
```json
{"message":"Email registered successfully."}
```

### Retrieve All Emails (GET)

**Windows (PowerShell):**
```powershell
Invoke-RestMethod -Uri http://localhost:8080/coming-soon `
  -Method GET `
  -Headers @{"X-Secret-Token"="my-secret-token-123"}
```

**Linux/Mac (curl):**
```bash
curl -X GET http://localhost:8080/coming-soon \
  -H "X-Secret-Token: my-secret-token-123"
```

**Expected Response:**
```json
["test@example.com"]
```

## Step 7: Check Your Data

The emails are stored in `data/emails.txt` on your host machine:

**Windows:**
```cmd
type data\emails.txt
```

**Linux/Mac:**
```bash
cat data/emails.txt
```

## Alternative: Using Docker Compose

For easier management, use Docker Compose:

### Step 1: Update the SECRET_TOKEN in docker-compose.yml

Edit `docker-compose.yml` and change:
```yaml
- SECRET_TOKEN=${SECRET_TOKEN:-change-this-to-a-secure-token}
```

### Step 2: Run with Docker Compose

```bash
docker-compose up -d
```

### Step 3: View Logs

```bash
docker-compose logs -f
```

### Step 4: Stop the Service

```bash
docker-compose down
```

## Common Commands

| Task | Command |
|------|---------|
| Stop container | `docker stop coming-soon-service` |
| Start container | `docker start coming-soon-service` |
| Restart container | `docker restart coming-soon-service` |
| View logs | `docker logs -f coming-soon-service` |
| Remove container | `docker rm -f coming-soon-service` |

## What's Next?

1. **Change the SECRET_TOKEN:** Use a secure random token for production
2. **Read the full README:** Check `README.md` for detailed documentation
3. **Integrate with your frontend:** Use the POST endpoint from your landing page
4. **Set up monitoring:** Track your container's health and performance
5. **Deploy to production:** Use a cloud provider or your own server

## Troubleshooting

### Port Already in Use
If port 8080 is already taken, change it:
```bash
docker run -d -p 3000:8080 ...  # Use port 3000 instead
```

### Container Won't Start
Check the logs for errors:
```bash
docker logs coming-soon-service
```

### Cannot Access API
1. Verify container is running: `docker ps`
2. Check if you can access from inside the container:
```bash
docker exec -it coming-soon-service wget -O- localhost:8080/coming-soon
```

### Permission Errors with Data Directory
On Linux/Mac, ensure the directory has proper permissions:
```bash
chmod 755 data
```

## Security Warning

⚠️ **IMPORTANT:** The default `SECRET_TOKEN` used in this guide is for demonstration only!

For production, generate a secure token:

**Linux/Mac:**
```bash
openssl rand -hex 32
```

**PowerShell:**
```powershell
-join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | % {[char]$_})
```

---

That's it! Your Coming Soon Service is now running. 🚀

For more details, see the full [README.md](README.md)