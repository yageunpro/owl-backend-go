package location

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yageunpro/owl-backend-go/internal/naver/search"
	"github.com/yageunpro/owl-backend-go/store"
	"github.com/yageunpro/owl-backend-go/store/location"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	Query(ctx context.Context, q []string) (*resQuery, error)
}

type service struct {
	store *store.Store
}

func New(sto *store.Store) Service {
	return &service{store: sto}
}

func (s *service) Query(ctx context.Context, q []string) (*resQuery, error) {
	if q == nil || len(q) == 0 {
		return &resQuery{Locations: make([]Location, 0)}, nil
	}
	queryString := strings.Join(q, " ")

	isCache := true
	out, err := s.store.Location.GetQuery(ctx, queryString)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			isCache = false
		} else {
			return nil, err
		}
	}

	if isCache && !out.QueryTime.Before(time.Now().Add(-24*time.Hour)) {
		res := resQuery{
			Locations: make([]Location, 0, len(out.Locations)),
		}
		for i := range len(out.Locations) {
			res.Locations = append(res.Locations, Location{
				Id:       out.Locations[i].Id,
				Title:    out.Locations[i].Title,
				Address:  out.Locations[i].Address,
				Category: out.Locations[i].Category,
				MapX:     out.Locations[i].MapX,
				MapY:     out.Locations[i].MapY,
			})
		}
		return &res, nil
	}

	searchResult, err := search.Query(ctx, queryString)
	var locations []location.Location
	if searchResult.Items == nil {
		locations = make([]location.Location, 0)
	} else {
		locations = make([]location.Location, 0, len(searchResult.Items))
		for i := range searchResult.Items {
			x, err := strconv.ParseInt(searchResult.Items[i].Mapx, 10, 32)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseInt(searchResult.Items[i].Mapy, 10, 32)
			if err != nil {
				return nil, err
			}

			locations = append(locations, location.Location{
				Id:       uuid.Must(uuid.NewV7()),
				Title:    searchResult.Items[i].Title,
				Address:  searchResult.Items[i].RoadAddress,
				Category: searchResult.Items[i].Category,
				MapX:     int(x),
				MapY:     int(y),
			})
		}
	}

	if isCache {
		err = s.store.Location.UpdateQuery(ctx, out.QueryId, locations)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.store.Location.SaveQuery(ctx, queryString, locations)
		if err != nil {
			return nil, err
		}
	}

	res := resQuery{
		Locations: make([]Location, 0, len(locations)),
	}
	for i := range locations {
		res.Locations = append(res.Locations, Location{
			Id:       locations[i].Id,
			Title:    locations[i].Title,
			Address:  locations[i].Address,
			Category: locations[i].Category,
			MapX:     locations[i].MapX,
			MapY:     locations[i].MapY,
		})
	}

	return &res, nil
}
