# Go URL Shortener Microservice

A high-performance, containerized URL shortener service built with Go (Golang) and Redis. Designed to demonstrate microservice architecture, REST API design, and container orchestration.

##  Features

*   **URL Shortening**: Generates unique 6-character short codes for long URLs.
*   **Fast Redirection**: High-performance HTTP redirects using Redis caching.
*   **Analytics**: Tracks click counts for every shortened link.
*   **Containerized**: Fully Dockerized for "write once, run anywhere" deployment.
*   **REST API**: Clean JSON-based API endpoints.

## Tech Stack

*   **Language**: Go (Golang) 1.23+
*   **Database**: Redis (Alpine)
*   **Infrastructure**: Docker & Docker Compose
*   **Testing**: Go testing package & Miniredis

## Installation & Run

You need **Docker** installed on your machine.

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/YOUR_USERNAME/url-shortener.git
    cd url-shortener
    ```

2.  **Start the services**:
    ```bash
    docker compose up --build -d
    ```

3.  **Verify it's running**:
    ```bash
    curl http://localhost:8080/health
    # or just shorten a URL
    curl -X POST -d '{"url": "https://google.com"}' http://localhost:8080/shorten
    ```

##  API Endpoints

| Method | Endpoint | Description | Body / Payload |
| :--- | :--- | :--- | :--- |
| `POST` | `/shorten` | Create a short URL | `{"url": "https://example.com"}` |
| `GET` | `/{code}` | Redirect to original | N/A |
| `GET` | `/analytics/{code}` | Get click stats | N/A |

### Example Usage

**Shorten a URL:**
```bash
curl -X POST -d '{"url": "https://github.com"}' http://localhost:8080/shorten
# Output: {"short_url": "http://localhost:8080/AbC123", "short_code": "AbC123"}
```

**Check Stats:**
```bash
curl http://localhost:8080/analytics/AbC123
# Output: {"short_code": "AbC123", "clicks": 5}
```

##  Testing

To run the unit and integration tests (mocking Redis):

```bash
go test -v ./...
```

## ☸️ Kubernetes Deployment

### Prerequisites
- **Docker Desktop** (Make sure Kubernetes is enabled in Settings -> Kubernetes -> Enable Kubernetes)
- `kubectl` (Included with Docker Desktop)

### Steps

1.  **Build the Docker Image**
    Since you are using Docker Desktop, the image built locally is immediately available to the local Kubernetes cluster.
    ```bash
    docker build -t url-shortener:latest .
    ```

2.  **Deploy Redis**
    ```bash
    kubectl apply -f k8s/redis.yaml
    ```

3.  **Deploy Configuration and Application**
    ```bash
    kubectl apply -f k8s/configmap.yaml
    kubectl apply -f k8s/deployment.yaml
    kubectl apply -f k8s/service.yaml
    ```

4.  **Verify Deployment**
    Check if all pods are running:
    ```bash
    kubectl get pods
    ```

5.  **Access the Application**
    Docker Desktop exposes NodePorts on `localhost`.
    
    **Test Health Endpoint:**
    ```bash
    curl http://localhost:30000/healthz
    # Should output: OK
    ```

    **Shorten a URL:**
    ```bash
    curl -X POST -d '{"url": "https://google.com"}' http://localhost:30000/shorten
    ```

### Troubleshooting
- **"no such host" or connection errors?**
  Ensure you are applying the Kubernetes manifests (`kubectl apply ...`) and *not* just running `docker-compose up` or the "Play" button in Docker Desktop's container list. The container list runs the default Docker Compose stack, which is different from the Kubernetes deployment.


