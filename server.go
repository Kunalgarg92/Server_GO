package main

import (
    "fmt"
    "net/http"
    "Server_GO/routes" 
    "github.com/rs/cors"
)

func main() {
    mux := http.NewServeMux()
    routes.RegisterRoutes(mux) 

    handler := cors.Default().Handler(mux)

    fmt.Println("Server is running on port 5000")
    if err := http.ListenAndServe(":5000", handler); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}
