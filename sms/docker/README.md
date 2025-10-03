# SMS VCS Backend Service

A Go-based SMS VCS backend service with Elasticsearch, Redis, and PostgreSQL support.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.24.1+ (for development)

### Running with Docker Compose

1. **Clone and navigate to the project:**
   ```bash
   cd /home/vudangducminh/Desktop/go-tutorial/vcs-be/sms/docker
   ```

2. **Start all services:**
   ```bash
   docker compose up -d
   ```

3. **Monitor service health:**
   ```bash
   docker compose ps
   ```

### Service Startup Order

The services start in the following order with health checks:

1. **Elasticsearch** starts first and must be healthy
2. **Redis Health Check** verifies Redis Cloud connectivity
3. **SMS Application** starts only after both Elasticsearch and Redis are healthy

### 🔍 Health Monitoring

#### Service Health Checks
```bash
# Check all services status
docker compose ps

# Check individual service health
docker inspect elasticsearch --format='{{.State.Health.Status}}'
docker inspect redis --format='{{.State.Health.Status}}'
docker inspect sms-app --format='{{.State.Health.Status}}'
```

#### Application Health Endpoint
```bash
# Check application health
curl http://localhost:8080/health

# Expected response:
{
  "status": "healthy",
  "services": {
    "elasticsearch": "healthy",
    "postgresql": "healthy", 
    "redis": "healthy"
  },
  "timestamp": "2025-09-03T10:00:00Z"
}
```

### 📡 Service URLs

- **SMS Application**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Swagger API Docs**: http://localhost:8080/swagger/index.html
- **Elasticsearch**: http://localhost:9200
- **Kibana**: http://localhost:5601

### 🛠 Development

#### Building the Application
```bash
# Build Docker image
docker compose build sms-app

# Rebuild and start
docker compose up --build sms-app
```

#### View Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f sms-app
docker compose logs -f elasticsearch
docker compose logs -f redis
```

### 🔧 Configuration

Environment variables are stored in `.env` file:
- Redis Cloud credentials
- Elasticsearch settings
- Application settings

### 🚦 Stopping Services

```bash
# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v
```

### 📊 Service Dependencies

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│Elasticsearch│    │  Redis Cloud │    │  PostgreSQL │
│   (Local)   │    │  (External)  │    │  (External) │
└──────┬──────┘    └──────┬───────┘    └──────┬──────┘
       │                  │                   │
       └──────────────────┼───────────────────┘
                          │
                    ┌─────▼──────┐
                    │  SMS App   │
                    │ (Port 8080)│
                    └────────────┘
```

### 🔐 Security

- Redis credentials are stored in `.env` file
- `.env` file is excluded from git via `.gitignore`
- Health checks use internal container networking

### 🐛 Troubleshooting

1. **Services not starting**: Check `docker compose logs <service-name>`
2. **Health checks failing**: Verify network connectivity and credentials
3. **Build failures**: Ensure Go modules are properly downloaded
4. **Port conflicts**: Make sure ports 8080, 9200, 5601 are available
