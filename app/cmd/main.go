package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jasonsavage/financeflow/config"
	"github.com/jasonsavage/financeflow/internal/db"
	"github.com/jasonsavage/financeflow/internal/repository"
	"github.com/jasonsavage/financeflow/internal/routes"
)

func main() {
	cfg := config.Load()

	pool := db.Connect(cfg)
	defer pool.Close()

	repo := repository.NewPostgresRepo(pool)

	r := routes.Register(repo)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server is running on port %s\n", cfg.AppPort)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
