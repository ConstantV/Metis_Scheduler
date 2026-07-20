package server

import "github.com/jackc/pgx/v5"

// Server zorgt ervoor dat al onze handlers dadelijk bij de database kunnen
type Server struct {
	db *pgx.Conn
}

// NewServer maakt een nieuwe server instantie aan met de database connectie
func NewServer(db *pgx.Conn) *Server {
	return &Server{
		db: db,
	}
}
