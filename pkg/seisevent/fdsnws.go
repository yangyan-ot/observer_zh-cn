package seisevent

import (
	"encoding/csv"
	"errors"
	"strings"
	"time"
)

func ParseFdsnwsEvent(dataText, timeLayout string) ([]Event, error) {
	// Convert to CSV format
	csvDataStr := strings.ReplaceAll(dataText, ",", " - ")
	csvDataStr = strings.ReplaceAll(csvDataStr, "|", ",")
	csvRecords, err := csv.NewReader(strings.NewReader(csvDataStr)).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(csvRecords) <= 1 {
		return nil, errors.New("no seismic event found")
	}

	var resultArr []Event
	for _, record := range csvRecords[1:] {
		var (
			seisEvent Event
			magType   string
		)
		for idx, val := range record {
			switch idx {
			case 0:
				seisEvent.Event = val
			case 1:
				seisEvent.Verfied = true
				if len(val) > len(timeLayout) {
					val = val[:len(timeLayout)]
				}
				t, err := time.Parse(timeLayout, val)
				if err != nil {
					return nil, err
				}
				seisEvent.Timestamp = t.UnixMilli()
			case 2:
				seisEvent.Latitude = string2Float(val)
			case 3:
				seisEvent.Longitude = string2Float(val)
			case 4:
				seisEvent.Depth = string2Float(val)
			case 9:
				magType = val
			case 10:
				seisEvent.Magnitude = []Magnitude{
					{Type: ParseMagnitude(magType), Value: string2Float(val)},
				}
			case 12:
				seisEvent.Region = val
			}
		}
		resultArr = append(resultArr, seisEvent)
	}

	return resultArr, nil
}
