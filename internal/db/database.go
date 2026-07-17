package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

// ConnectDB zet de verbinding op met jouw lokale Postgres
func ConnectDB() *pgx.Conn {
	// PAS DIT AAN: Vul hier jouw eigen Postgres gebruikersnaam en wachtwoord in!
	// Structuur: postgres://GEBRUIKER:WACHTWOORD@localhost:5432/metis_scheduler
	connStr := "postgres://postgres:wachtwoord@localhost:5432/metis_scheduler"

	// Als je een omgevingsvariabele gebruikt, pakken we die
	if envUrl := os.Getenv("DATABASE_URL"); envUrl != "" {
		connStr = envUrl
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Kan geen verbinding maken met de database: %v\n", err)
	}

	// Test of de verbinding echt werkt
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database reageert niet (Ping-fout): %v\n", err)
	}

	fmt.Println("Successfully connected to your local PostgreSQL database!")
	return conn
}
