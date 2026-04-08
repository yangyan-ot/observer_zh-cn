package seisevent

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const SCEA_ID = "scea"

type SCEA struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
	timeSource      *timesource.Source
}

func (s *SCEA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      SCEA_ID,
		Country: "CN",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Sichuan Earthquake Administration (Early Warning)",
			"zh-TW": "四川地震局預警",
		},
	}
}

func (s *SCEA) buildUrlSign(params map[string]string) string {
	const (
		SIGN_TOKEN = "OMEoiAuaExMuTjpovqKrhYDkZMkUaoCE"
		SLAT       = "earthquake_app"
	)

	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	raw := ""
	for i, k := range keys {
		if i > 0 {
			raw += "&"
		}
		raw += k + "=" + params[k]
	}
	raw += "&token=" + SIGN_TOKEN

	sum := md5.Sum([]byte(raw + SLAT))
	return hex.EncodeToString(sum[:])
}

func (s *SCEA) getRequestUrl(lat, lon float64) string {
	params := map[string]string{
		"pageNo":    "1",
		"pageSize":  "50",
		"orderType": "1",
		"userLng":   strconv.FormatFloat(lon, 'f', -1, 64),
		"userLat":   strconv.FormatFloat(lat, 'f', -1, 64),
		"timeStamp": strconv.FormatInt(s.timeSource.Now().UnixMilli(), 10),
	}

	params["sign"] = s.buildUrlSign(params)

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return "http://118.113.105.29:8002/api/earlywarning/jsonPageList?" + values.Encode()
}

func (s *SCEA) GetEvents(latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if s.cache.Valid() {
		baseEvents = s.cache.Get()
	} else {
		v, err, _ := s.sf.Do(SCEA_ID, func() (any, error) {
			if s.cache.Valid() {
				return s.cache.Get(), nil
			}

			res, err := request.GET(
				s.getRequestUrl(latitude, longitude),
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse SCEA_B JSON response
			var dataMap map[string]any
			err = json.Unmarshal(res, &dataMap)
			if err != nil {
				return nil, err
			}

			// Check server response
			if dataMap["code"].(float64) != 0 {
				return nil, fmt.Errorf("server error: %s", dataMap["msg"])
			}

			dataMapEvents, ok := dataMap["data"].([]any)
			if !ok {
				return nil, errors.New("seismic event data is not available")
			}

			// Ensure the response has the expected keys and values
			expectedKeys := []string{"eventId", "shockTime", "longitude", "latitude", "placeName", "magnitude", "depth", "infoTypeName"}

			var resultArr []Event
			for _, event := range dataMapEvents {
				if !isMapHasKeys(event.(map[string]any), expectedKeys) || !isMapKeysEmpty(event.(map[string]any), expectedKeys) {
					continue
				}

				resultArr = append(resultArr, Event{
					Depth:     -1,
					Verfied:   event.(map[string]any)["infoTypeName"].(string) == "[正式]",
					Event:     event.(map[string]any)["eventId"].(string),
					Region:    event.(map[string]any)["placeName"].(string),
					Latitude:  event.(map[string]any)["latitude"].(float64),
					Longitude: event.(map[string]any)["longitude"].(float64),
					Magnitude: s.getMagnitude(event.(map[string]any)["magnitude"].(float64)),
					Timestamp: time.UnixMilli(int64(event.(map[string]any)["shockTime"].(float64))).UnixMilli(),
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			s.cache.Set(sortedEvents)
			return sortedEvents, nil
		})
		if err != nil {
			return nil, err
		}

		baseEvents = v.([]Event)
	}

	for i := range baseEvents {
		baseEvents[i].Distance = getDistance(latitude, baseEvents[i].Latitude, longitude, baseEvents[i].Longitude)
		baseEvents[i].Estimation = getSeismicEstimation(
			s.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (s *SCEA) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
