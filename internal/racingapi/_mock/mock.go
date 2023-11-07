package mock

import (
	"context"
	"encoding/json"
	"race-weekend-bot/internal/racingapi"
)

var testData = `
{
    "MRData": {
        "xmlns": "http://ergast.com/mrd/1.5",
        "series": "f1",
        "url": "http://ergast.com/api/f1/current/next.json",
        "limit": "30",
        "offset": "0",
        "total": "1",
        "RaceTable": {
            "season": "2023",
            "round": "20",
            "Races": [
                {
                    "season": "2023",
                    "round": "20",
                    "url": "https://en.wikipedia.org/wiki/2023_S%C3%A3o_Paulo_Grand_Prix",
                    "raceName": "São Paulo Grand Prix",
                    "Circuit": {
                        "circuitId": "interlagos",
                        "url": "http://en.wikipedia.org/wiki/Aut%C3%B3dromo_Jos%C3%A9_Carlos_Pace",
                        "circuitName": "Autódromo José Carlos Pace",
                        "Location": {
                            "lat": "-23.7036",
                            "long": "-46.6997",
                            "locality": "São Paulo",
                            "country": "Brazil"
                        }
                    },
                    "date": "2023-11-05",
                    "time": "17:00:00Z",
                    "FirstPractice": {
                        "date": "2023-11-03",
                        "time": "14:30:00Z"
                    },
                    "Qualifying": {
                        "date": "2023-11-03",
                        "time": "18:00:00Z"
                    },
                    "SecondPractice": {
                        "date": "2023-11-04",
                        "time": "14:30:00Z"
                    },
                    "Sprint": {
                        "date": "2023-11-04",
                        "time": "18:30:00Z"
                    }
                }
            ]
        }
    }
}
`

type MockApiResponse struct {
	MRData struct {
		RaceTable struct {
			Races []Race `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type MockAPI struct {
	URL string
}

type Event struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Race struct {
	RaceName       string `json:"raceName"`
	Date           string `json:"date"`
	Time           string `json:"time"`
	FirstPractice  Event  `json:"FirstPractice"`
	SecondPractice Event  `json:"SecondPractice"`
	ThirdPractice  Event  `json:"ThirdPractice"`
	Qualifying     Event  `json:"Qualifying"`
	Sprint         Event  `json:"Sprint"`
}

type MockResponse struct {
	StatusCode int
	Body       []byte
}

func (a MockAPI) GetData(ctx context.Context, logger racingapi.Logger) (racingapi.RaceWeekendSchedule, error) {

	output := racingapi.RaceWeekendSchedule{}

	var apiResponse MockApiResponse

	err := json.Unmarshal([]byte(testData), &apiResponse)
	if err != nil {
		return output, err
	}

	fornatted := formatresponse(apiResponse)
	if fornatted.Race.Date == "" {
		return output, racingapi.ErrNoRaceThisWeekend
	}

	return fornatted, err
}

func formatresponse(r MockApiResponse) racingapi.RaceWeekendSchedule {
	race := r.MRData.RaceTable.Races[0]
	return racingapi.RaceWeekendSchedule{
		RaceType: racingapi.SERIES_f1,
		RaceName: race.RaceName,
		Race: racingapi.Event{
			Date: race.Date,
			Time: race.Time,
		},
		FirstPractice: racingapi.Event{
			Date: race.FirstPractice.Date,
			Time: race.FirstPractice.Time,
		},
		SecondPractice: racingapi.Event{
			Date: race.SecondPractice.Date,
			Time: race.SecondPractice.Time,
		},
		ThirdPractice: racingapi.Event{
			Date: race.ThirdPractice.Date,
			Time: race.ThirdPractice.Time,
		},
		Qualifying: racingapi.Event{
			Date: race.Qualifying.Date,
			Time: race.Qualifying.Time,
		},
		Sprint: racingapi.Event{
			Date: race.Sprint.Date,
			Time: race.Sprint.Time,
		},
	}
}
