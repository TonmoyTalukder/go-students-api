# Student API : GO REST API SERVER
A simple, modular, RESTful API built in Go to manage student records using SQLite.

## 🚀 Features
✅ Create, Read, Update, Delete (CRUD) student records

✅ Pagination support for listing students

✅ Structured logging using log/slog

✅ Graceful shutdown on SIGINT/SIGTERM

✅ YAML-based configuration

✅ Idiomatic, modular Go project structure

## 📁 Project Structure
```
.
├── cmd/
│   └── student-api/           # Main application entrypoint
├── config/                    # YAML-based config files
├── internal/
│   ├── config/                # Config loader
│   ├── http/
│   │   └── handlers/
│   │       └── student/       # All student-related handlers
│   └── storage/
│       └── sqlite/            # SQLite implementation
├── types/                     # Shared types (e.g., Student)
├── utils/
│       └── response/          # Response Structure implemented
├── storage/                   # storage.db
├── .gitignore
├── go.mod
├── go.sum
├── student-api.postman_collection.json
└── README.md
```

## ⚙️ Configuration
Create a config file at config/local.yaml:
```yml
env: "dev"
storage_path: "storage/storage.db"
http_server:
  address: "localhost:8082"
``` 

## 🏁 Getting Started
#### Prerequisites
- Go 1.24.2
- SQLite (included with Go stdlib via database/sql)

#### Run the Server
```bash
go run ./cmd/student-api/main.go  -config ./config/local.yaml 
```
**Server starts on:**
```arduino
http://localhost:8082
```

## 📡 API Endpoints
| Method | Endpoint           | Description                         |
| ------ | ------------------ | ----------------------------------- |
| GET    | /api/              | Welcome route                       |
| GET    | /api/students      | List students (supports pagination) |
| GET    | /api/students/{id} | Get student by ID                   |
| POST   | /api/students      | Create student                      |
| PUT    | /api/students/{id} | Update student by ID                |
| DELETE | /api/students/{id} | Delete student by ID                |

## 📬 Postman Collection

You can import the API into Postman using the pre-built collection:

🔗 [Download Postman Collection](./student-api.postman_collection.json)

To use:

1. Open Postman.
2. Click on "Import" → "Upload Files".
3. Select the downloaded JSON file.
4. Run requests against: `http://localhost:8082`

💡 Tip: Set a Postman environment variable base_url = http://localhost:8082 for convenience.


## 📝 Sample Payload
Student JSON for POST/PUT:
```json
{
  "name": "Tonmoy Talukder",
  "email": "tonmoy@mail.com",
  "age": 25
}
```

## 🧪 Sample CURL Request
Create a student:
```bash
curl -X POST http://localhost:8082/api/students \
  -H "Content-Type: application/json" \
  -d '{"name":"Tonmoy", "email":"tonmoy@mail.com", "age":25}'
```

## 📦 Response Format
All API responses follow a consistent format:
```json
{
  "status": "OK",
  "code": 200,
  "message": "Student updated successfully",
  "data": {
    "id": 2,
    "name": "Tonmoy Updated",
    "email": "tonmoy.updated@mail.com",
    "age": 26
  },
  "meta": null
}
```

Error:
```json
{
  "status": "Error",
  "code": 400,
  "message": "Invalid JSON body",
  "error": "unexpected EOF"
}
```

## 🛑 Graceful Shutdown
The server shuts down cleanly on Ctrl+C or kill signal using a 5-second timeout context.

## 👨‍💻 Author
**Tonmoy Talukder** <br/>
GitHub: https://github.com/TonmoyTalukder

## 📄 License
MIT License © 2025 Tonmoy Talukder