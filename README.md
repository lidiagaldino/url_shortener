# 🔗 URL Shortener in Go

A simple URL shortener built in Go with DDD + Clean Architecture and MongoDB.

### Phase 1 Features
- Shorten URLs
- Redirect to the original URL

### Phase 2 Features
- Register and Login
- Protected routes

---

## Technologies

- Go 1.25  
- MongoDB  
- Chi Router  
- DDD + Clean Architecture  

---

## Project Structure

```bash
.
├── cmd/server           # Ponto de entrada da aplicação  
├── internal  
│   ├── services         # Services and DTOs 
│   ├── config           # Env vars 
│   ├── domain           # Entities, interfaces and custom exceptions  
│   ├── infrastructure   # MongoDB and JWT Implementations 
│   └── interfaces       # Handlers HTTP and Middlewares
├── pkg                  # ID generator and hasher
├── docker-compose.yml   # MongoDB + APP
├── Dockerfile           # Multistage Build
└── README.md  
```

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

### Create User

**POST /users**

Request Body:

```
{
  "name": "John Doe",
  "password": "123",
  "email": "john@example.com"
}
```

Response:

```
{
  "id": "abc123",
  "name": "John Doe",
  "email": "john@example.com"
}
```

---

### Sign In User

**POST /users/signin**

Request Body:

```
{
  "password": "123",
  "email": "john@example.com"
}
```

Response:

```
{
  "token": "myjwttoken"
}
```

---

### Shorten URL

**POST /shorten (protected)**

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
  "owner_id": "123"
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

## Next Steps (Phase 3)

- Custom domains per user  
- Click tracking and statistics  
- Simple web interface to manage URLs  
