package seisevent

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const BCSF_ID = "bcsf"

type BCSF struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (c *BCSF) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      BCSF_ID,
		Country: "FR",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "French Central Seismological Office",
			"zh-TW": "法國中央地震局",
		},
	}
}

func (c *BCSF) GetEvents(latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err, _ := c.sf.Do(BCSF_ID, func() (any, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			res, err := request.GET(
				"https://api.franceseisme.fr/fdsnws/event/1/query?minmag=-1&format=text&timezone=UTC&limit=300&orderby=time",
				30*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			resultArr, err := ParseFdsnwsEvent(string(res), "2006-01-02T15:04:05")
			if err != nil {
				return nil, err
			}
			sorted := sortSeismicEvents(resultArr)
			c.cache.Set(sorted)
			return sorted, nil
		})
		if err != nil {
			return nil, err
		}

		baseEvents = v.([]Event)
	}

	for i := range baseEvents {
		baseEvents[i].Distance = getDistance(latitude, baseEvents[i].Latitude, longitude, baseEvents[i].Longitude)
		baseEvents[i].Estimation = getSeismicEstimation(
			c.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}
