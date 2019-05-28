package weather

import (
	"errors"
	. "fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

const (
	NOT_IMPLEMENTED = "Not Implemented"
	NOT_FOUND = "timestamp not found"
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// WeatherService is the backing service invoked by HTTP/REST handlers. Add
// whatever stateful fields you may need here.
type WeatherService struct {
	measurements sync.Map
}

// NewWeatherService creates an instance of the weather service struct
func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

// GetMeasurement retrieves a single measurement based on timestamp
func (s *WeatherService) GetMeasurement1(timestamp time.Time) (*Measurement, error) {

	if j, ok := s.measurements.Load(timestamp); !ok{
		Println("Timestamp info doesn't exist !")
		return &Measurement{}, errors.New(NOT_FOUND)
	}else {
		Println("Timestamp map found")
		m,ok := j.(map[string]float32)
		if !ok{
			Println("there was an error in type assertion")
		}
		mObject := Measurement{}
		mObject.Timestamp = timestamp
		mObject.Metrics = m

		return &mObject, nil
		}

	// for code fulfillment sake
	return &Measurement{}, nil
}

// CreateMeasurement creates and stores a new measurement
func (s *WeatherService) CreateMeasurement(newMeasurement Measurement) error {

	if _, ok := s.measurements.LoadOrStore(newMeasurement.Timestamp, newMeasurement.Metrics); !ok{

		Println("Measurement created in memory !")
		return nil
		}else {
		return errors.New("Timestamp info exists")
		}

	return nil
}

// GetStats obtains a list of statistics from the system based on stats, metrics, and a time range
func (s *WeatherService) GetStats(stats []string, metrics []string, from time.Time, to time.Time) ([]StatisticRow, error) {

	return nil, errors.New(NOT_IMPLEMENTED)
}
