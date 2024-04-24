package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"scheduler_api/api"
	"scheduler_api/logger"
)

type Holiday struct {
	HolidayName     string
	ActualDate      string
	ObservedDate    string
	Type            string
	DayOfWeek       string
	DayOfWeekNumber string
	WeekOfYear      string
	Day             string
	Month           string
	Year            string
}

func fetchHolidays() ([]Holiday, error) {
	request := api.NewRequest(
		"https://api.public-holidays.nz/v1",
		"/all",
		api.GET)
	request.Params["apikey"] = "49ef83deb13f4c6fab015b1df1ebd1a5"
	client := api.Client{
		Request: request,
	}
	res, err := client.SendRequest()
	if err != nil {
		return nil, fmt.Errorf("Holiday fetching failed: %s", err)
	}
	defer res.Body.Close()
	reader := io.Reader(res.Body)
	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, reader)
	if err != nil {
		return nil, fmt.Errorf("Holiday fetching failed: %s", err)
	}
	var holidays []Holiday
	err = json.Unmarshal(buffer.Bytes(), &holidays)
	if err != nil {
		logger.E("json.Unmarshal failed, err", err)
		return nil, fmt.Errorf("Holiday json unmarshal filed: %s", err)
	}
	var publicHolidays []Holiday
	for _, holiday := range holidays {
		if holiday.Type == "National" {
			publicHolidays = append(publicHolidays, holiday)
		}
	}
	return publicHolidays, nil
}
