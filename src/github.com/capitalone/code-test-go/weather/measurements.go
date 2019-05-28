/*
 * This file contains the HTTP handlers as well as JSON serialization,
 * de-serialization, and validation code to verify payloads. You should
 * not need to modify any of this code unless you want to modify the interface
 * used by the service
 */
package weather

import (
	"encoding/json"
	. "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type measurementsService interface {
	GetMeasurement1(timestamp time.Time) (*Measurement, error)
	CreateMeasurement(newMeasurement Measurement) error
}

// MeasurementsHandler is a RESTful HTTP endpoint for dealing with measurements
type MeasurementsHandler struct {
	svc measurementsService
}

// NewMeasurementsHandler creates a new measurement handler instance
func NewMeasurementsHandler(svc measurementsService) *MeasurementsHandler {
	return &MeasurementsHandler{
		svc,
	}
}

// CreateMeasurement is the HTTP POST handler for /measurements
func (h *MeasurementsHandler) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	Info.Println("Handler invoked")
	decoder := json.NewDecoder(r.Body)

	var measurement Measurement
	err := decoder.Decode(&measurement)
	if err != nil {
		Error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	Info.Printf("Attempting to create measurement: %+v\n", measurement)
	err = h.svc.CreateMeasurement(measurement)
	if err != nil {
		Error.Println(err)
		if err.Error() == NOT_IMPLEMENTED {
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	measureTime := measurement.Timestamp.Format(time.RFC3339Nano)
	loc := Sprintf("/measurements/%s", measureTime)
	w.Header().Set("Location", strconv.Quote(loc))
	w.WriteHeader(http.StatusCreated)
}

// GetMeasurement is the HTTP GET handler for /measurements/{timestamp}
func (h *MeasurementsHandler) GetMeasurement(w http.ResponseWriter, r *http.Request) {
	Info.Println("Handler invoked")
	vars := mux.Vars(r)
	ts, exists := vars["timestamp"]
	if !exists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	timestamp, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	Info.Printf("Attempting to get measurement for timestamp %v", timestamp)
	measurement, err := h.svc.GetMeasurement1(timestamp)
	if err != nil {
		Error.Println(err)
		if err.Error() == NOT_FOUND {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	bytes, err := json.Marshal(measurement)
	if err != nil {
		Error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	Println("writing response")

	w.Write(bytes)

	// this is sending multiple headwrites, w.write is already writing a header for status 200
	//w.WriteHeader(http.StatusOK)
}
