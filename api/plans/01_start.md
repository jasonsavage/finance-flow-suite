Agent, I would like you to build the REST api backend for my application called "Finance Flow".

The api will support a frontend web application for tracking a user's spending over time and make predictions on future spending month over month.

This project will be written in Go 1.26 and connect to a Postgres DB running in a Docker container.

App routing should be handled using the external module `github.com/go-chi` and connections to postgres will use `github.com/jackc/pgx`.

The scafolding for the project should loosely follow this:
```
api/
├── plans/               # Claude agent planning docs
│── cmd/
│──── web/          
│       └── main.go      # Entry point — starts server, sets up DB connection, creates routes
│── internal/
│   ├── db/
│   │   └── db.go          # DB connection pool (pgx or sqlx)
│   │   └── repository.go  # DB sql commands used as a ORM for DB access from handlers
│   ├── models/
│   │   └── *.go         # Structs mapping to DB tables
│   ├── handlers/
│   │   └── *.go         # HTTP handler functions
│   ├── routes/
│   │   └── routes.go    # Route registration
│   └── middleware/
│       └── *.go         # Logging, CORS, auth, etc.
└── config/
│   └── config.go        # Loads env vars into a config struct
├── migrations/          # Table definitions run on first Postgres start
├── .env                 # Env variable for go aplication and docker compose
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

Once the initial scafolding has been created, each agent will be responsable for implementing each feature in the ./plans/features folder.