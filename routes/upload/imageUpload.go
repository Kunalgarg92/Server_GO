package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/api/upload", uploadImageHandler)
}

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Upload handler invoked")
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        fmt.Println("Error parsing form:", err)
        http.Error(w, "Error parsing form", http.StatusInternalServerError)
        return
    }
    file, handler, err := r.FormFile("image")
    if err != nil {
        fmt.Println("Error retrieving the file:", err)
        http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
        return
    }
    defer file.Close()
    fmt.Println("File retrieved:", handler.Filename)
    if err := ensureUploadDir(); err != nil {
        fmt.Println("Error creating upload directory:", err)
        http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
        return
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

func ensureUploadDir() error {
    if _, err := os.Stat("upload"); os.IsNotExist(err) {
        return os.Mkdir("upload", os.ModePerm)
    }
    return nil
}
