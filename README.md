# Go URL Shortener Microservice

A high-performance, containerized URL shortener service built with Go (Golang) and Redis. Designed to demonstrate microservice architecture, REST API design, and container orchestration.

## 🚀 Features

*   **URL Shortening**: Generates unique 6-character short codes for long URLs.
*   **Fast Redirection**: High-performance HTTP redirects using Redis caching.
*   **Analytics**: Tracks click counts for every shortened link.
*   **Containerized**: Fully Dockerized for "write once, run anywhere" deployment.
*   **REST API**: Clean JSON-based API endpoints.

## 🛠️ Tech Stack

*   **Language**: Go (Golang) 1.23+
*   **Database**: Redis (Alpine)
*   **Infrastructure**: Docker & Docker Compose
*   **Testing**: Go testing package & Miniredis

## 📦 Installation & Run

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

## 🔌 API Endpoints

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

## 🧪 Testing

To run the unit and integration tests (mocking Redis):

```bash
go test -v ./...
```
