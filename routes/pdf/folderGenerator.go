package pdf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/api/generate-pdf", GeneratePDF)
}

type InspectionData struct {
    TruckSerialNumber string `json:"Truck Serial Number"`
    TruckModel        string `json:"Truck Model"`
    InspectionID      int    `json:"Inspection ID"`
    InspectorName     string `json:"Inspector Name"`
    // Add more fields as needed
}

func GeneratePDF(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        InspectionData []InspectionData `json:"inspectionData"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    data := requestData.InspectionData

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(10, 10, 10)
    pdf.AddPage()

    pdf.SetFont("Times", "B", 20)
    pdf.SetTextColor(0, 51, 102)
    pdf.Cell(0, 10, "Inspection Report")
    pdf.Ln(20)

    for _, inspection := range data {
        pdf.SetFont("Times", "B", 16)
        pdf.SetTextColor(0, 0, 128)
        pdf.CellFormat(0, 10, inspection.TruckModel, "", 1, "C", false, 0, "")
        pdf.Ln(4)

        val := reflect.ValueOf(inspection)
        typ := reflect.TypeOf(inspection)

        pdf.SetFont("Times", "", 12)
        pdf.SetTextColor(0, 0, 0)
        pdf.SetFillColor(240, 240, 240) 

        for i := 0; i < val.NumField(); i++ {
            fieldName := typ.Field(i).Tag.Get("json")
            fieldValue := val.Field(i).Interface()

            pdf.CellFormat(50, 8, fieldName+":", "1", 0, "L", true, 0, "")
            pdf.CellFormat(0, 8, formatFieldValue(fieldValue), "1", 1, "L", false, 0, "")
        }

        pdf.Ln(10)
        pdf.SetFont("Times", "I", 12)
        pdf.Cell(0, 10, "Images:")
        pdf.Ln(12)
        pdf.SetFont("Times", "", 12)
        pdf.Cell(0, 10, "[Insert images here]")
        pdf.Ln(10)
        pdf.SetDrawColor(128, 128, 128)
        pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
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

func formatFieldValue(value interface{}) string {
    switch v := value.(type) {
    case int:
        return strconv.Itoa(v)
    case float64:
        return fmt.Sprintf("%.2f", v)
    default:
        return fmt.Sprintf("%v", v)
    }
}
