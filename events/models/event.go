package models

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/tech-nimble/go-tools/events"
)

type Event struct {
	ID         int
	Type       int
	Payload    pgtype.JSON
	Headers    pgtype.JSONB
	Status     int
	ModelID    int
	CreatedAt  pgtype.Timestamp
	UpdatedAt  sql.NullTime
	Exchange   string
	RoutingKey string
}

func NewEvent(entity *events.Event) (*Event, error) {
	var payloadPgJson pgtype.JSON
	if err := payloadPgJson.Set(entity.Payload); err != nil {
		return nil, err
	}

	var headersPgJsonb pgtype.JSONB
	if err := headersPgJsonb.Set(entity.Headers); err != nil {
		return nil, err
	}

	createdAt := pgtype.Timestamp{}
	if err := createdAt.Set(entity.CreatedAt); err != nil {
		return nil, err
	}

	return &Event{
		ID:         entity.ID,
		Type:       1,
		Payload:    payloadPgJson,
		Headers:    headersPgJsonb,
		Status:     entity.Status,
		ModelID:    entity.ModelID,
		CreatedAt:  createdAt,
		UpdatedAt:  sql.NullTime{Time: entity.UpdatedAt, Valid: !entity.UpdatedAt.IsZero()},
		Exchange:   entity.Exchange,
		RoutingKey: entity.RoutingKey,
	}, nil
}

func (e *Event) GetEntity() (*events.Event, error) {
	var headers map[string]any
	if err := json.Unmarshal(e.Headers.Bytes, &headers); err != nil {
		return nil, err
	}

	return &events.Event{
		ID:         e.ID,
		Payload:    e.Payload,
		Headers:    headers,
		Status:     e.Status,
		ModelID:    e.ModelID,
		RoutingKey: e.RoutingKey,
		Exchange:   e.Exchange,
		CreatedAt:  e.CreatedAt.Time,
		UpdatedAt:  e.UpdatedAt.Time,
	}, nil
}
