package f1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"race-weekend-bot/internal/racingapi"
	"race-weekend-bot/internal/utils"

	"golang.org/x/net/context/ctxhttp"
)

type F1ApiResponse struct {
	MRData struct {
		RaceTable struct {
			Races []Race `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type F1API struct {
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

type F1Response struct {
	StatusCode int
	Body       []byte
}

func (a F1API) GetData(ctx context.Context, logger racingapi.Logger) (racingapi.RaceWeekendSchedule, error) {
	source := "f1.api.source"
	log := logger.With("source", source)
	output := racingapi.RaceWeekendSchedule{}

	var apiResponse F1ApiResponse

	response, err := a.Call(ctx, "GET", logger)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(response.Body, &apiResponse)
	if err != nil {
		log.Error("unamrshal error: %w")
		return output, err
	}

	fornatted := formatresponse(apiResponse)
	if fornatted.Race.Date == "" {
		return output, racingapi.NoRaceThisWeekend
	}

	eventDateTime := fmt.Sprintf("%sT%s", fornatted.Race.Date, fornatted.Race.Time)

	isUpcoming, err := utils.IsDateOnCurrentWeek(eventDateTime)
	if err != nil {
		return output, err
	}

	if !isUpcoming {
		return output, racingapi.NoRaceThisWeekend
	}
	return fornatted, err
}

func (a F1API) Call(ctx context.Context, method string, logger racingapi.Logger) (F1Response, error) {
	source := "f1.api.source.call"
	log := logger.With("source", source)
	log.Info(fmt.Sprintf("sending request to: %s", a.URL))

	req, err := http.NewRequestWithContext(ctx, method, a.URL, http.NoBody)
	if err != nil {
		log.Error(fmt.Sprintf("error calling %s: %s", a.URL, err.Error()))
		return F1Response{
			StatusCode: 500,
			Body:       nil,
		}, err
	}

	resp, err := ctxhttp.Do(ctx, &http.Client{}, req)
	if err != nil {
		log.Error(fmt.Sprintf("error calling %s: %s", a.URL, err.Error()))
		return F1Response{
			StatusCode: resp.StatusCode,
			Body:       nil,
		}, err
	}
	defer resp.Body.Close()

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("error reading body: %w", err)
		return F1Response{
			StatusCode: 500,
			Body:       nil,
		}, err
	}

	log.Info(fmt.Sprintf("response status code: %d", resp.StatusCode))

	log.Info(fmt.Sprintf("Got response: %s", bodyByte))

	return F1Response{
		StatusCode: resp.StatusCode,
		Body:       bodyByte,
	}, nil
}

func formatresponse(r F1ApiResponse) racingapi.RaceWeekendSchedule {
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
