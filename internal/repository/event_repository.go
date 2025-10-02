package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
)

// EventModel represents the database model for domain events
type EventModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Sequence    int64     `gorm:"autoIncrement;primaryKey;index"`
	AggregateID uuid.UUID `gorm:"type:uuid;not null;index"`
	Type        string    `gorm:"type:varchar(100);not null;index"`
	Version     int64     `gorm:"not null;index"`
	Timestamp   time.Time `gorm:"not null;index"`
	Data        string    `gorm:"type:text;not null"`
}

// TableName returns the table name for EventModel
func (EventModel) TableName() string {
	return "domain_events"
}

// EventRepository implements the domain EventStoreRepository interface using GORM
type EventRepository struct {
	db *gorm.DB
}

// NewEventStoreRepository creates a new EventStoreRepository
func NewEventStoreRepository(db *gorm.DB) domain.EventStoreRepository {
	return &EventRepository{db: db}
}

// SaveEvent saves a domain event to the database
func (r *EventRepository) SaveEvent(ctx context.Context, event domain.DomainEvent) error {
	eventModel := &EventModel{
		ID:          event.GetID(),
		AggregateID: event.GetAggregateID(),
		Type:        event.GetType(),
		Version:     int64(event.GetVersion()),
		Timestamp:   event.GetTimestamp(),
		Data:        r.serializeEventData(event.GetData()),
	}

	return r.db.WithContext(ctx).Create(eventModel).Error
}

// GetEventsByAggregateID retrieves all events for a specific aggregate
func (r *EventRepository) GetEventsByAggregateID(ctx context.Context, aggregateID uuid.UUID) ([]domain.DomainEvent, error) {
	var eventModels []EventModel
	err := r.db.WithContext(ctx).
		Where("aggregate_id = ?", aggregateID).
		Order("timestamp ASC").
		Find(&eventModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events by aggregate ID: %w", err)
	}

	events := make([]domain.DomainEvent, 0, len(eventModels))
	for _, eventModel := range eventModels {
		event, err := r.reconstructEvent(eventModel)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct event %s: %w", eventModel.ID, err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventsByType retrieves all events of a specific type
func (r *EventRepository) GetEventsByType(ctx context.Context, eventType string) ([]domain.DomainEvent, error) {
	var eventModels []EventModel
	err := r.db.WithContext(ctx).
		Where("type = ?", eventType).
		Order("timestamp ASC").
		Find(&eventModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}

	events := make([]domain.DomainEvent, 0, len(eventModels))
	for _, eventModel := range eventModels {
		event, err := r.reconstructEvent(eventModel)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct event %s: %w", eventModel.ID, err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventsByTimeRange retrieves events within a specific time range
func (r *EventRepository) GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]domain.DomainEvent, error) {
	var eventModels []EventModel
	err := r.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp ASC").
		Find(&eventModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events by time range: %w", err)
	}

	events := make([]domain.DomainEvent, 0, len(eventModels))
	for _, eventModel := range eventModels {
		event, err := r.reconstructEvent(eventModel)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct event %s: %w", eventModel.ID, err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventVersion retrieves the current version of an aggregate
func (r *EventRepository) GetEventVersion(ctx context.Context, aggregateID uuid.UUID) (int64, error) {
	var result struct {
		MaxVersion int64
	}

	err := r.db.WithContext(ctx).
		Model(&EventModel{}).
		Select("COALESCE(MAX(version), 0)").
		Where("aggregate_id = ?", aggregateID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get event version: %w", err)
	}

	return result.MaxVersion, nil
}

// UpdateEventVersion updates the version of an aggregate (typically not needed as version is stored with each event)
func (r *EventRepository) UpdateEventVersion(ctx context.Context, aggregateID uuid.UUID, version int64) error {
	// This is typically not needed as we store the version with each event
	// But we can implement it if needed for specific use cases
	return nil
}

// GetAllEvents retrieves all events from the database
func (r *EventRepository) GetAllEvents(ctx context.Context) ([]domain.DomainEvent, error) {
	var eventModels []EventModel
	err := r.db.WithContext(ctx).
		Order("timestamp ASC").
		Find(&eventModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}

	events := make([]domain.DomainEvent, 0, len(eventModels))
	for _, eventModel := range eventModels {
		event, err := r.reconstructEvent(eventModel)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct event %s: %w", eventModel.ID, err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventsAfterSequence retrieves events after a specific sequence number
func (r *EventRepository) GetEventsAfterSequence(ctx context.Context, sequence int64) ([]domain.DomainEvent, error) {
	var eventModels []EventModel
	err := r.db.WithContext(ctx).
		Where("sequence > ?", sequence).
		Order("sequence ASC").
		Limit(1000). // Limit to prevent huge queries
		Find(&eventModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events after sequence: %w", err)
	}

	events := make([]domain.DomainEvent, 0, len(eventModels))
	for _, eventModel := range eventModels {
		event, err := r.reconstructEvent(eventModel)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct event %s: %w", eventModel.ID, err)
		}
		events = append(events, event)
	}

	return events, nil
}

// serializeEventData serializes event data to JSON
func (r *EventRepository) serializeEventData(data map[string]interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// If marshaling fails, store as empty string to avoid database errors
		return ""
	}
	return string(jsonData)
}

// reconstructEvent reconstructs a domain event from the database model
func (r *EventRepository) reconstructEvent(eventModel EventModel) (domain.DomainEvent, error) {
	// This is a simplified reconstruction
	// In a real implementation, you might want to use a registry pattern
	// to reconstruct specific event types with their proper data structures

	data := make(map[string]interface{})
	if eventModel.Data != "" {
		err := json.Unmarshal([]byte(eventModel.Data), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
	}

	// Create a base domain event with the stored data
	baseEvent := domain.NewBaseDomainEvent(
		eventModel.AggregateID,
		eventModel.Type,
		data,
	)

	// Set the actual ID and timestamp from the database
	baseEvent.ID = eventModel.ID
	baseEvent.Timestamp = eventModel.Timestamp

	return &baseEvent, nil
}
