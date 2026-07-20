package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"metis-scheduler/internal/db" // Aangepast naar internal/db
	"metis-scheduler/internal/server"
)

func main() {
	port := ":8080"
	fmt.Printf("Metis Scheduler start op http://localhost%s...\n", port)

	// 1. Start de database connectie (nu via het 'db' pakket)
	dbConn := db.ConnectDB()
	defer dbConn.Close(context.Background())

	// 2. Start de server
	myServer := server.NewServer(dbConn)
	mux := http.NewServeMux()

	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	})
	mux.HandleFunc("/swagger", swaggerUIHandler)
	mux.HandleFunc("POST /projects/{id}/schedule", myServer.HandleSchedule)
	mux.HandleFunc("GET /projects/{id}/critical-path", myServer.HandleCriticalPath)

	log.Fatal(http.ListenAndServe(port, mux))
}

func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Metis Scheduler API Docs</title>
		<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
		<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
	</head>
	<body>
		<div id="swagger-ui"></div>
		<script>
			window.onload = () => {
				window.ui = SwaggerUIBundle({
					url: '/openapi.yaml',
					dom_id: '#swagger-ui',
				});
			};
		</script>
	</body>
	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
