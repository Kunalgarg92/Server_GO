package routes

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

type InspectionData struct {
	TruckSerialNumber string `json:"Truck Serial Number"`
	TruckModel        string `json:"Truck Model"`
	InspectionID      int    `json:"Inspection ID"`
	InspectorName     string `json:"Inspector Name"`
}

func GeneratePDF(w http.ResponseWriter, r *http.Request) {
	var data []InspectionData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	for _, inspection := range data {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, inspection.TruckModel)
		pdf.Ln(10)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		http.Error(w, "Could not generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(buf.Bytes())
}
