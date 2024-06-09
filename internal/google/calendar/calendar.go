package calendar

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log/slog"
	"time"
)

type Event struct {
	CalcId    string
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type Param struct {
	Source    oauth2.TokenSource
	SyncToken string
}

type Result struct {
	SyncToken string
	Events    []Event
}

func GetAllEvents(ctx context.Context, arg Param) (*Result, error) {
	srv, err := calendar.NewService(ctx, option.WithTokenSource(arg.Source))
	if err != nil {
		return nil, errors.Join(errors.New("failed to init calendar client"), err)
	}

	krTimeZone := time.FixedZone("KR", 9*60*60)
	qry := srv.Events.List("primary").TimeZone("Asia/Seoul").SingleEvents(true)

	res := Result{
		SyncToken: "",
		Events:    make([]Event, 0),
	}

	useSyncToken := false
	pageToken := ""
	if arg.SyncToken != "" {
		useSyncToken = true
	}

	for {
		q := qry
		if useSyncToken {
			q = qry.SyncToken(arg.SyncToken)
		}
		if pageToken != "" {
			q = qry.PageToken(pageToken)
		}
		events, err := q.Do()
		if err != nil {
			if useSyncToken {
				useSyncToken = false
				continue
			}
			return nil, errors.Join(errors.New("failed to list events"), err)
		}

		pageToken = events.NextPageToken

		for i := range events.Items {
			var startTime, endTime time.Time

			if events.Items[i].Start == nil || events.Items[i].End == nil {
				slog.Error("failed to parse time, nil", slog.String("title", events.Items[i].Summary))
				continue
			}
			startTime, err = time.Parse(time.RFC3339, events.Items[i].Start.DateTime)
			if err != nil {
				startTime, err = time.ParseInLocation(time.DateOnly, events.Items[i].Start.Date, krTimeZone)
				if err != nil {
					slog.Error("failed to parse start time", slog.String("event", events.Items[i].Start.DateTime))
					continue
				}
			}
			endTime, err = time.Parse(time.RFC3339, events.Items[i].End.DateTime)
			if err != nil {
				endTime, err = time.ParseInLocation(time.DateOnly, events.Items[i].End.Date, krTimeZone)
				if err != nil {
					slog.Error("failed to parse end time", slog.String("event", events.Items[i].Start.DateTime))
					continue
				}
			}

			res.Events = append(res.Events, Event{
				CalcId:    events.Items[i].Id,
				Title:     events.Items[i].Summary,
				StartTime: startTime,
				EndTime:   endTime,
			})
		}

		if pageToken == "" || events.NextSyncToken != "" {
			res.SyncToken = events.NextSyncToken
			break
		}
	}

	return &res, nil
}
