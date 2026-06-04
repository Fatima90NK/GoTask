
# GoTask - Task Management API

A simple task management API built with Go, MongoDB, and Kubernetes.

## Features

- ✅ Create tasks
- ✅ View all tasks
- ✅ Mark tasks as done
- ✅ Runs on Kubernetes
- ✅ MongoDB Atlas integration
- ✅ Auto-scaling and recovery

## Prerequisites

- Go 1.21+
- Docker
- Kubernetes (Docker Desktop or Minikube)
- MongoDB Atlas account
- Docker Hub account

## Local Setup

### 1. Clone the repository
```bash
git clone https://github.com/YOUR_USERNAME/gotask.git
cd gotask
```

### 2. Set environment variables
```bash
export MONGO_URI="mongodb+srv://USERNAME:PASSWORD@cluster0.zaxqypx.mongodb.net/gotask?retryWrites=true&w=majority"
```

### 3. Run the app
```bash
go run main.go
```

The server runs on `http://localhost:8080`

## API Endpoints

### Get all tasks
```bash
curl http://localhost:8080/tasks
```

### Create a task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"name":"Buy milk","done":false}'
```

### Mark task as done
```bash
curl -X POST http://localhost:8080/tasks/done \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

## Docker Setup

### 1. Build Docker image
```bash
docker build -t gotask:v1 .
```

### 2. Run Docker container
```bash
docker run -e MONGO_URI="your-connection-string" -p 8080:8080 gotask:v1
```

## Kubernetes Deployment

### 1. Create secrets
Copy `k8s/secret.yaml.example` to `k8s/secret.yaml` and add your MongoDB URI.

### 2. Deploy
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### 3. Check status
```bash
kubectl get all -n gotask
```

### 4. Access the service
```bash
curl http://localhost/tasks
```

## CI/CD Pipeline

This project uses GitHub Actions for automated deployment:
1. Tests the code
2. Builds Docker image
3. Pushes to Docker Hub
4. Deploys to Kubernetes

### Setup GitHub Secrets

Add these to your GitHub repo settings:
- `DOCKER_USERNAME` - Your Docker Hub username
- `DOCKER_PASSWORD` - Your Docker Hub access token
- `KUBE_CONFIG` - Your Kubernetes config (base64 encoded)

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────┐      ┌──────────────┐
│  GoTask API │─────▶│  MongoDB     │
│ (2 replicas)│      │  Atlas       │
└─────────────┘      └──────────────┘
   (Kubernetes)
```

## Monitoring

### View logs
```bash
kubectl logs -n gotask -l app=gotask --tail=50
```

### Watch pods
```bash
kubectl get pods -n gotask --watch
```

### Check service
```bash
kubectl get svc -n gotask
```

## Environment Variables

- `MONGO_URI` - MongoDB Atlas connection string (required)

## Author

Fatima90NK

## License

MIT
```