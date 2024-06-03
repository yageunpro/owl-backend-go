package location

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/internal/query"
)

type Store interface {
	GetQuery(ctx context.Context, queryString string) (*resGetQuery, error)
	SaveQuery(ctx context.Context, queryString string, locations []Location) error
	UpdateQuery(ctx context.Context, queryId uuid.UUID, locations []Location) error
	GetLocation(ctx context.Context, locationId uuid.UUID) (*Location, error)
}

type store struct {
	pool *db.Pool
}

func New(pool *db.Pool) Store {
	return &store{pool: pool}
}

func (s *store) GetQuery(ctx context.Context, queryString string) (*resGetQuery, error) {
	qry := query.New(s.pool)

	row, err := qry.FindQueryId(ctx, queryString)
	if err != nil {
		return nil, err
	}

	rows, err := qry.GetLocationWithQueryId(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	res := resGetQuery{
		QueryId:   row.ID,
		QueryTime: row.UpdatedAt.Time,
		Locations: make([]Location, len(rows)),
	}
	for i := range rows {
		res.Locations[i].Id = rows[i].ID
		res.Locations[i].Title = rows[i].Title
		res.Locations[i].Address = rows[i].Address
		res.Locations[i].Category = rows[i].Category
		res.Locations[i].MapX = int(rows[i].MapX)
		res.Locations[i].MapY = int(rows[i].MapY)
	}

	return &res, nil
}

func (s *store) SaveQuery(ctx context.Context, queryString string, locations []Location) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to begin transaction"), err)
	}
	defer tx.Rollback(ctx)

	qry := query.New(s.pool).WithTx(tx)

	queryId := uuid.Must(uuid.NewV7())
	err = qry.CreateQuery(ctx, query.CreateQueryParams{
		ID:   queryId,
		Data: queryString,
	})
	if err != nil {
		return errors.Join(errors.New("failed to create query"), err)
	}

	for i := range locations {
		err = qry.CreateLocation(ctx, query.CreateLocationParams{
			ID:       locations[i].Id,
			QueryID:  queryId,
			Title:    locations[i].Title,
			Address:  locations[i].Address,
			Category: locations[i].Category,
			MapX:     int32(locations[i].MapX),
			MapY:     int32(locations[i].MapY),
		})
		if err != nil {
			return errors.Join(errors.New("failed to create location"), err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}

	return nil
}

func (s *store) UpdateQuery(ctx context.Context, queryId uuid.UUID, locations []Location) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to begin transaction"), err)
	}
	defer tx.Rollback(ctx)

	qry := query.New(s.pool).WithTx(tx)

	err = qry.DeprecateQuery(ctx, queryId)
	if err != nil {
		return errors.Join(errors.New("failed to deprecate query"), err)
	}

	for i := range locations {
		err = qry.CreateLocation(ctx, query.CreateLocationParams{
			ID:       uuid.Must(uuid.NewV7()),
			QueryID:  queryId,
			Title:    locations[i].Title,
			Address:  locations[i].Address,
			Category: locations[i].Category,
			MapX:     int32(locations[i].MapX),
			MapY:     int32(locations[i].MapY),
		})
		if err != nil {
			return errors.Join(errors.New("failed to create location"), err)
		}
	}

	err = qry.UpdateQueryTime(ctx, queryId)
	if err != nil {
		return errors.Join(errors.New("failed to update query"), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}

	return nil
}

func (s *store) GetLocation(ctx context.Context, locationId uuid.UUID) (*Location, error) {
	qry := query.New(s.pool)

	row, err := qry.GetLocation(ctx, locationId)
	if err != nil {
		return nil, err
	}

	return &Location{
		Id:       row.ID,
		Title:    row.Title,
		Address:  row.Address,
		Category: row.Category,
		MapX:     int(row.MapX),
		MapY:     int(row.MapY),
	}, nil
}
