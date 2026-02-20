package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// InMemoryDB is an in-memory implementation of the DB interface.
// For TEST purposes only! Do not use in production!
type InMemoryDB struct {
	mu                      sync.RWMutex
	alerts                  map[string]json.RawMessage
	issues                  map[string]*inMemoryIssueRecord
	moveMappings            map[string]json.RawMessage
	channelProcessingStates map[string]*ChannelProcessingState
}

type inMemoryIssueRecord struct {
	channelID     string
	correlationID string
	postID        string
	isOpen        bool
	body          json.RawMessage
}

// NewInMemoryDB creates a new InMemoryDB instance.
// For TEST purposes only! Do not use in production!
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		alerts:                  make(map[string]json.RawMessage),
		issues:                  make(map[string]*inMemoryIssueRecord),
		moveMappings:            make(map[string]json.RawMessage),
		channelProcessingStates: make(map[string]*ChannelProcessingState),
	}
}

// Init is a no-op for the in-memory implementation.
func (db *InMemoryDB) Init(_ context.Context, _ bool) error {
	return nil
}

// SaveAlert saves an alert to the in-memory store.
func (db *InMemoryDB) SaveAlert(_ context.Context, alert *Alert) error {
	if alert == nil {
		return errors.New("alert is nil")
	}

	body, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.alerts[alert.UniqueID()] = body

	return nil
}

// SaveIssue creates or updates a single issue.
func (db *InMemoryDB) SaveIssue(_ context.Context, issue Issue) error {
	if issue == nil {
		return errors.New("issue is nil")
	}

	body, err := issue.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal issue: %w", err)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.issues[issue.UniqueID()] = &inMemoryIssueRecord{
		channelID:     issue.ChannelID(),
		correlationID: issue.GetCorrelationID(),
		postID:        issue.CurrentPostID(),
		isOpen:        issue.IsOpen(),
		body:          body,
	}

	return nil
}

// SaveIssues creates or updates multiple issues.
func (db *InMemoryDB) SaveIssues(ctx context.Context, issues ...Issue) error {
	for _, issue := range issues {
		if err := db.SaveIssue(ctx, issue); err != nil {
			return err
		}
	}

	return nil
}

// MoveIssue moves an issue from one channel to another.
// Returns an error if sourceChannelID and targetChannelID are the same.
// If the issue does not exist in the store, this is a no-op.
func (db *InMemoryDB) MoveIssue(_ context.Context, issue Issue, sourceChannelID, targetChannelID string) error {
	if sourceChannelID == targetChannelID {
		return errors.New("source and target channel IDs are the same")
	}

	if issue == nil {
		return errors.New("issue is nil")
	}

	body, err := issue.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal issue: %w", err)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	record, ok := db.issues[issue.UniqueID()]
	if !ok {
		return nil
	}

	record.channelID = targetChannelID
	record.postID = issue.CurrentPostID()
	record.isOpen = issue.IsOpen()
	record.body = body

	return nil
}

// FindOpenIssueByCorrelationID finds a single open issue by channel ID and correlation ID.
// Returns an error if channelID or correlationID are empty, or if multiple open issues match.
func (db *InMemoryDB) FindOpenIssueByCorrelationID(_ context.Context, channelID, correlationID string) (string, json.RawMessage, error) {
	if channelID == "" {
		return "", nil, errors.New("channelID is required")
	}

	if correlationID == "" {
		return "", nil, errors.New("correlationID is required")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	var foundID string
	var foundRecord *inMemoryIssueRecord

	for id, record := range db.issues {
		if record.channelID == channelID && record.correlationID == correlationID && record.isOpen {
			if foundRecord != nil {
				return "", nil, fmt.Errorf("multiple open issues found for channel %q and correlationID %q", channelID, correlationID)
			}

			foundID = id
			foundRecord = record
		}
	}

	if foundRecord == nil {
		return "", nil, nil
	}

	return foundID, foundRecord.body, nil
}

// FindIssueBySlackPostID finds a single issue by channel ID and Slack post ID.
// Returns an error if channelID or postID are empty.
func (db *InMemoryDB) FindIssueBySlackPostID(_ context.Context, channelID, postID string) (string, json.RawMessage, error) {
	if channelID == "" {
		return "", nil, errors.New("channelID is required")
	}

	if postID == "" {
		return "", nil, errors.New("postID is required")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	for id, record := range db.issues {
		if record.channelID == channelID && record.postID == postID {
			return id, record.body, nil
		}
	}

	return "", nil, nil
}

// FindActiveChannels returns a list of all channels that have at least one open issue.
func (db *InMemoryDB) FindActiveChannels(_ context.Context) ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	seen := make(map[string]struct{})

	for _, record := range db.issues {
		if record.isOpen {
			seen[record.channelID] = struct{}{}
		}
	}

	channels := make([]string, 0, len(seen))
	for channelID := range seen {
		channels = append(channels, channelID)
	}

	return channels, nil
}

// LoadOpenIssuesInChannel loads all open issues for the specified channel.
func (db *InMemoryDB) LoadOpenIssuesInChannel(_ context.Context, channelID string) (map[string]json.RawMessage, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make(map[string]json.RawMessage)

	for id, record := range db.issues {
		if record.channelID == channelID && record.isOpen {
			result[id] = record.body
		}
	}

	return result, nil
}

// SaveMoveMapping creates or updates a move mapping.
func (db *InMemoryDB) SaveMoveMapping(_ context.Context, moveMapping MoveMapping) error {
	if moveMapping == nil {
		return errors.New("moveMapping is nil")
	}

	body, err := moveMapping.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal move mapping: %w", err)
	}

	key := moveMappingKey(moveMapping.ChannelID(), moveMapping.GetCorrelationID())

	db.mu.Lock()
	defer db.mu.Unlock()

	db.moveMappings[key] = body

	return nil
}

// FindMoveMapping finds a move mapping by channel ID and correlation ID.
// Returns an error if channelID or correlationID are empty.
func (db *InMemoryDB) FindMoveMapping(_ context.Context, channelID, correlationID string) (json.RawMessage, error) {
	if channelID == "" {
		return nil, errors.New("channelID is required")
	}

	if correlationID == "" {
		return nil, errors.New("correlationID is required")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.moveMappings[moveMappingKey(channelID, correlationID)], nil
}

// DeleteMoveMapping deletes a move mapping. No error is returned if the mapping does not exist.
func (db *InMemoryDB) DeleteMoveMapping(_ context.Context, channelID, correlationID string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.moveMappings, moveMappingKey(channelID, correlationID))

	return nil
}

// SaveChannelProcessingState creates or updates a channel processing state.
func (db *InMemoryDB) SaveChannelProcessingState(_ context.Context, state *ChannelProcessingState) error {
	if state == nil {
		return errors.New("state is nil")
	}

	stateCopy := *state

	db.mu.Lock()
	defer db.mu.Unlock()

	db.channelProcessingStates[state.ChannelID] = &stateCopy

	return nil
}

// FindChannelProcessingState finds a channel processing state by channel ID.
// Returns nil without an error if no state is found.
func (db *InMemoryDB) FindChannelProcessingState(_ context.Context, channelID string) (*ChannelProcessingState, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	state, ok := db.channelProcessingStates[channelID]
	if !ok {
		return nil, nil //nolint:nilnil // DB interface contract: return nil, nil when not found
	}

	stateCopy := *state

	return &stateCopy, nil
}

// DropAllData clears all data from the in-memory store.
func (db *InMemoryDB) DropAllData(_ context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.alerts = make(map[string]json.RawMessage)
	db.issues = make(map[string]*inMemoryIssueRecord)
	db.moveMappings = make(map[string]json.RawMessage)
	db.channelProcessingStates = make(map[string]*ChannelProcessingState)

	return nil
}

func moveMappingKey(channelID, correlationID string) string {
	return channelID + "\x00" + correlationID
}
