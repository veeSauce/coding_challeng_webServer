package weather

import (
	"errors"
	. "fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

const (
	NOT_IMPLEMENTED = "Not Implemented"
	NOT_FOUND       = "timestamp not found"
	BAD_REQUEST     = "bad request"
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
func (s *WeatherService) GetMeasurement(timestamp time.Time) (*Measurement, error) {

	if j, ok := s.measurements.Load(timestamp); !ok {
		Println("Timestamp info doesn't exist, returning empty struct")
		// app will return Measurement object without any metrics data
		emptyMatrix := make(map[string]float32)
		return &Measurement{Timestamp:timestamp, Metrics:emptyMatrix}, errors.New("Timestamp info doesn't exist")
	} else {
		Println("Timestamp map found")
		m, ok := j.(map[string]float32)
		if !ok {
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

	if _, ok := s.measurements.LoadOrStore(newMeasurement.Timestamp, newMeasurement.Metrics); !ok {

		Println("Measurement created in memory !")
		return nil
	} else {
		return errors.New("Timestamp info exists")
	}

	return nil
}

// GetStats obtains a list of statistics from the system based on stats, metrics, and a time range
func (s *WeatherService) GetStats(stats []string, metrics []string, from time.Time, to time.Time) ([]StatisticRow, error) {

	//initialize a variable to return
	sttSlice := []StatisticRow{}
	stt := StatisticRow{}

	// keeps track of results from
	tempDataVar := []*Measurement{}

	// increment this guy until you hit "to" timeStamp
	timeVar := from

	// all the data is in tempDataVar. Putting it in a slice will make min, max, average arithmetic easy
	metricMap := make(map[string][]float32)

	//	//	//	//
	// timestamps given are precise and incerementing by a factor of 10 minutes - as per feature files/test cases + use case story, data is pushed at strict intervals
	//	//	//	//

	// input validation to make sure time stamps are not invalid
	if to.Sub(from) < 0 {
		Println("please enter from time older than to time")
		return sttSlice, errors.New(BAD_REQUEST)

	} else if from == to {
		Print("from time stamp equals to")
		tempVar, _ := s.GetMeasurement(from)
			// Special Case !

			// for every metric
			for _, metric := range metrics {
				Print("ranging over: metric")

				// get every stat
				for _, stat := range stats {

					// since special case, value is the same
					switch stat {
					case "min", "max", "average":
						stt := StatisticRow{
							Metric: metric,
							Stat:   stat,
							Value:  tempVar.Metrics[metric],
						}
						sttSlice = append(sttSlice, stt)

					default:
						return []StatisticRow{}, errors.New("Invalid stat")
					}
				}
			}
			return sttSlice, nil

	} else {

		// get each measurement value for from timestamp incrementing by 10 minutes until "to" timestamp
		for to.Sub(timeVar) > 0 {
			data, err := s.GetMeasurement(timeVar)
			if err != nil {
				return sttSlice, errors.New("Problem with data")
			}
			tempDataVar = append(tempDataVar, data)
			timeVar = timeVar.Add(10 * time.Minute)
		}

		floatSlice := []float32{}

		// for every metric
		for _, metric := range metrics {

			// flush out the slice for new metric data to be loaded
			floatSlice = nil

			// if metric is in measurements data, collect it
			for _, mapVal := range tempDataVar {

				if val, ok := mapVal.Metrics[metric]; ok {
					Println("value found for ", metric)
					floatSlice = append(floatSlice, val)
					metricMap[metric] = floatSlice
				}
			}
		}


	}

	// we have the data, now check for what stats are asked for

	for _, metric := range metrics {
		if slc, ok := metricMap[metric]; ok {
			for _, stat := range stats {
				switch stat {

				case "min":

					// to get the min, sort the slice in ascending order and pick the first indexed value
					sort.SliceStable(slc, func(i, j int) bool { return slc[i] < slc[j] })
					stt.Stat = stat
					stt.Value = slc[0]
					stt.Metric = metric
					sttSlice = append(sttSlice, stt)

				case "max":

					// to get the max, sort the slice in ascending order and pick the last indexed value
					sort.SliceStable(slc, func(i, j int) bool { return slc[i] < slc[j] })

					stt.Stat = stat
					stt.Value = slc[len(slc)-1]
					stt.Metric = metric
					sttSlice = append(sttSlice, stt)

				case "average":

					var sum float32

					for _, val := range slc {
						sum += val

					}
					stt.Stat = stat
					stt.Value = sum / float32(len(slc))
					stt.Metric = metric
					sttSlice = append(sttSlice, stt)
				}
			}
		}
	}

	return sttSlice, nil
}
