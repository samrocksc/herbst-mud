package repository

import (
	"context"
	"time"
)

// ChannelSubscription represents a character's subscription to a chat channel.
type ChannelSubscription struct {
	ID       int
	CharID   int
	Channel  string
	Enabled  bool
	Color    string
}

// OfflineTell represents a message sent to an offline player.
type OfflineTell struct {
	ID            int
	FromID        int
	RecipientID   int
	RecipientName string
	Message       string
	QueuedAt      time.Time
}

// IgnoreEntry represents an ignore relationship between characters.
type IgnoreEntry struct {
	ID         int
	CharacterID int
	IgnoredID   int
	CreatedAt   time.Time
}

// ChannelSubscriptionRepo defines data access for channel subscriptions.
type ChannelSubscriptionRepo interface {
	ListByCharacter(ctx context.Context, charID int) ([]*ChannelSubscription, error)
	GetByCharacterAndChannel(ctx context.Context, charID int, channel string) (*ChannelSubscription, error)
	Create(ctx context.Context, charID int, channel, color string) (*ChannelSubscription, error)
	UpdateEnabled(ctx context.Context, id int, enabled bool) error
	UpdateColor(ctx context.Context, id int, color string) error
	Delete(ctx context.Context, id int) error
	DeleteByCharacterAndChannel(ctx context.Context, charID int, channel string) error
}

// OfflineTellRepo defines data access for offline tells.
type OfflineTellRepo interface {
	ListByRecipient(ctx context.Context, recipientID int) ([]*OfflineTell, error)
	ListByRecipientName(ctx context.Context, recipientName string) ([]*OfflineTell, error)
	Create(ctx context.Context, fromID int, recipientName string, message string) (*OfflineTell, error)
	Delete(ctx context.Context, id int) error
	DeleteByRecipient(ctx context.Context, recipientID int) error
}

// IgnoreRepo defines data access for ignore lists.
type IgnoreRepo interface {
	ListByCharacter(ctx context.Context, charID int) ([]*IgnoreEntry, error)
	Exists(ctx context.Context, charID, ignoredID int) (bool, error)
	Create(ctx context.Context, charID, ignoredID int) (*IgnoreEntry, error)
	Delete(ctx context.Context, id int) error
	DeleteByCharacterAndIgnored(ctx context.Context, charID, ignoredID int) error
}
