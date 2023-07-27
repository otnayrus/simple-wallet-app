package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/otnayrus/simple-wallet-app/delivery/rest"
	"github.com/otnayrus/simple-wallet-app/repository"
	"github.com/otnayrus/simple-wallet-app/service"
)

func main() {
	db, err := sql.Open("sqlite3", "wallet.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	migrate(db)

	walletRepo := repository.NewWalletRepositiory(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := rest.NewWalletHandler(walletService)

	router := gin.Default()
	v1 := router.Group("/api/v1")

	v1.POST("/init", walletHandler.Initialize)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("Server started at 127.0.0.1:8000")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http.ListenAndServe: %v\n", err)
	}

}

func migrate(db *sql.DB) {
	createWalletsStmt := `
		CREATE TABLE IF NOT EXISTS wallets (
			id string primary key,
			owned_by string not null unique,
			token string not null unique,
			status int not null,
			updated_at timestamp,
			balance int not null
		);
	`
	_, err := db.Exec(createWalletsStmt)
	if err != nil {
		log.Fatal(err)
	}

	createMutationsStmt := `
		CREATE TABLE IF NOT EXISTS mutations (
			id string primary key,
			reference_id string unique,
			sender string not null,
			recipient string not null,
			action int not null,
			amount int not null
		);
	`
	_, err = db.Exec(createMutationsStmt)
	if err != nil {
		log.Fatal(err)
	}
}
