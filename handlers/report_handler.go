package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetReportToday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	today := time.Now().Format("2006-01-02")

	report, err := h.service.GetReport(today, today)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	// Redirect if old parameters are used
	if from != "" || to != "" {
		if startDate == "" {
			startDate = from
		}
		if endDate == "" {
			endDate = to
		}
		newQuery := r.URL.Query()
		newQuery.Del("from")
		newQuery.Del("to")
		newQuery.Set("start_date", startDate)
		newQuery.Set("end_date", endDate)
		
		url := r.URL.Path + "?" + newQuery.Encode()
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}

	if startDate == "" || endDate == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Remove GetReportByRange as it's replaced by GetReport with redirects

