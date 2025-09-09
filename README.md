# 🔗 URL Shortener in Go

A simple URL shortener built in Go with DDD + Clean Architecture and MongoDB.

### Phase 1 Features
- Shorten URLs
- Redirect to the original URL

---

## Technologies

- Go 1.25  
- MongoDB  
- Chi Router  
- DDD + Clean Architecture  

---

## Project Structure

- `cmd/server` - Application entry point  
- `internal/services` - Use cases / services  
- `internal/config` - Configuration (env vars)  
- `internal/domain` - Entities and interfaces  
- `internal/infra` - Implementations (MongoDB)  
- `internal/handlers` - HTTP handlers / controllers  
- `pkg` - Utilities (e.g., ID generation)  
- `docker-compose.yml` - MongoDB + App setup  
- `Dockerfile` - Multi-stage build  

---

## How to Run

**Prerequisites:** Docker + Docker Compose

Run the application:

```
docker compose up --build
```

API available at: `http://localhost:8080`  
MongoDB at: `mongodb://mongo:27017`

---

## Endpoints

### Shorten URL

**POST /shorten**

Request Body:

```
{
  "url": "https://mysite.com/very-long-article"
}
```

Response:

```
{
  "id": "abc123",
  "original_url": "https://mysite.com/very-long-article",
  "short_url": "http://localhost:8080/abc123"
}
```

---

### Redirect

**GET /{id}**

Example:

```
GET /abc123
```

Redirects to `https://mysite.com/very-long-article`.

---

## Tests

Run unit tests:

```
go test ./...
```

---

## Next Steps (Phase 2)

- User accounts (login/register)  
- Custom domains per user  
- Click tracking and statistics  
- Simple web interface to manage URLs  
