package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/tech-nimble/go-tools/events"
	"github.com/tech-nimble/go-tools/events/models"
	"github.com/tech-nimble/go-tools/helpers/jaeger"
	"github.com/tech-nimble/go-tools/repositories"
)

type Events struct {
	*repositories.Repositories
	tableName string
}

func NewEvents(rep *repositories.Repositories, tableName string) *Events {
	return &Events{
		Repositories: rep,
		tableName:    tableName,
	}
}

func (e *Events) GetByID(ctx context.Context, id int) (*events.Event, error) {
	span, ctx := jaeger.StartSpanFromContext(ctx, "DB 'events':GetByID")
	defer span.Finish()

	db := e.GetConnect(ctx)

	m := &models.Event{}

	if err := db.QueryRow(ctx, fmt.Sprintf(
		"select id, type, payload, headers, status, model_id, exchange, routing_key, created_at, updated_at"+
			" from %s where id=$1",
		e.tableName,
	), id).Scan(
		&m.ID,
		&m.Type,
		&m.Payload,
		&m.Headers,
		&m.Status,
		&m.ModelID,
		&m.Exchange,
		&m.RoutingKey,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	entity, err := m.GetEntity()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (e *Events) Create(ctx context.Context, entity *events.Event) error {
	span, ctx := jaeger.StartSpanFromContext(ctx, "DB 'events':Create")
	defer span.Finish()

	db := e.GetConnect(ctx)

	if entity.ID != 0 {
		return errors.New("нельзя создать существующую модель")
	}

	query := fmt.Sprintf(
		"insert into %s(type, payload, headers, status, model_id, created_at, updated_at, exchange, routing_key)"+
			" values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id",
		e.tableName,
	)

	m, err := models.NewEvent(entity)
	if err != nil {
		return err
	}

	rows, err := db.Query(
		ctx,
		query,
		m.Type,
		m.Payload,
		m.Headers,
		m.Status,
		m.ModelID,
		m.CreatedAt,
		m.UpdatedAt,
		m.Exchange,
		m.RoutingKey,
	)
	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&m.ID)

		if err != nil {
			return err
		}

		entity.ID = m.ID
	}

	return nil
}

func (e *Events) Update(ctx context.Context, entity *events.Event) error {
	span, ctx := jaeger.StartSpanFromContext(ctx, "DB 'events':Update")
	defer span.Finish()

	db := e.DB

	if entity.ID == 0 {
		return errors.New("нельзя обновить не существующую модель ")
	}

	query := fmt.Sprintf(
		"update %s set status = $1, updated_at = $2 where id = $3 returning id",
		e.tableName,
	)

	m, err := models.NewEvent(entity)
	if err != nil {
		return err
	}

	rows, err := db.Query(
		ctx,
		query,
		m.Status,
		m.UpdatedAt,
		m.ID,
	)
	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	defer rows.Close()

	var id int
	for rows.Next() {
		err = rows.Scan(&id)

		if err != nil {
			return err
		}

		if id != m.ID {
			return errors.New("не удалось обновить модель")
		}
	}

	return nil
}
