package main

import (
    "Server_GO/routes/pdf"   
    "Server_GO/routes/upload" 
    "fmt"
    "net/http"

    "github.com/rs/cors"
)

func main() {
    mux := http.NewServeMux()

    // Register routes from different packages
    pdf.RegisterRoutes(mux)
    upload.RegisterRoutes(mux)

    handler := cors.Default().Handler(mux)

    fmt.Println("Server is running on port 5000")
    if err := http.ListenAndServe(":5000", handler); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}
