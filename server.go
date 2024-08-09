package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", uploadImageHandler)

	handler := cors.Default().Handler(mux)

	fmt.Println("Server is running on port 5000")
	if err := http.ListenAndServe(":5000", handler); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload handler invoked")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Println("File retrieved:", handler.Filename)

	if _, err := os.Stat("upload"); os.IsNotExist(err) {
		err := os.Mkdir("upload", os.ModePerm)
		if err != nil {
			fmt.Println("Error creating upload directory:", err)
			http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
			return
		}
	}

	filePath := filepath.Join("upload", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		fmt.Println("Error saving the file:", err)
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	fmt.Println("File saved successfully:", filePath)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "message": "File uploaded successfully!"}`))
}
