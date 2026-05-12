package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"herbst-server/repository"
)

// chatService implements ChatService using repository interfaces.
type chatService struct {
	charRepo        repository.CharacterRepo
	subRepo         repository.ChannelSubscriptionRepo
	tellRepo        repository.OfflineTellRepo
	ignoreRepo      repository.IgnoreRepo
}

// NewChatService creates a new ChatService.
func NewChatService(
	charRepo repository.CharacterRepo,
	subRepo repository.ChannelSubscriptionRepo,
	tellRepo repository.OfflineTellRepo,
	ignoreRepo repository.IgnoreRepo,
) ChatService {
	return &chatService{
		charRepo:   charRepo,
		subRepo:    subRepo,
		tellRepo:   tellRepo,
		ignoreRepo: ignoreRepo,
	}
}

// DefaultChannels returns the list of available chat channels.
func DefaultChannels() []string {
	return []string{"ooc", "guild", "group", "trade", "help", "race", "clan"}
}

func (s *chatService) SendSay(ctx context.Context, charID, roomID int, message string) (*MessageResult, error) {
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	return &MessageResult{
		FromCharacterID:   charID,
		FromCharacterName: char.Name,
		ToCharacterIDs:    []int{},
		Channel:           "say",
		Message:           message,
		Type:              "say",
		RoomID:            &roomID,
	}, nil
}

func (s *chatService) SendYell(ctx context.Context, charID, roomID int, message string) (*MessageResult, error) {
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	return &MessageResult{
		FromCharacterID:   charID,
		FromCharacterName: char.Name,
		ToCharacterIDs:    []int{},
		Channel:           "yell",
		Message:           message,
		Type:              "yell",
		RoomID:            &roomID,
	}, nil
}

func (s *chatService) SendShout(ctx context.Context, charID int, message string) (*MessageResult, error) {
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	return &MessageResult{
		FromCharacterID:   charID,
		FromCharacterName: char.Name,
		ToCharacterIDs:    []int{},
		Channel:           "shout",
		Message:           message,
		Type:              "shout",
	}, nil
}

func (s *chatService) SendTell(ctx context.Context, fromID, toID int, message string) (*MessageResult, error) {
	fromChar, err := s.charRepo.Get(ctx, fromID)
	if err != nil {
		return nil, fmt.Errorf("sender not found: %w", err)
	}
	_, err = s.charRepo.Get(ctx, toID)
	if err != nil {
		return nil, fmt.Errorf("recipient not found: %w", err)
	}
	ignored, err := s.ignoreRepo.Exists(ctx, toID, fromID)
	if err != nil {
		return nil, fmt.Errorf("checking ignore: %w", err)
	}
	if ignored {
		return nil, fmt.Errorf("player is ignoring you")
	}
	return &MessageResult{
		FromCharacterID:   fromID,
		FromCharacterName: fromChar.Name,
		ToCharacterIDs:    []int{toID},
		Channel:           "tell",
		Message:           message,
		Type:              "tell",
	}, nil
}

func (s *chatService) SendWhisper(ctx context.Context, fromID, toID int, message string) (*MessageResult, error) {
	fromChar, err := s.charRepo.Get(ctx, fromID)
	if err != nil {
		return nil, fmt.Errorf("sender not found: %w", err)
	}
	_, err = s.charRepo.Get(ctx, toID)
	if err != nil {
		return nil, fmt.Errorf("recipient not found: %w", err)
	}
	ignored, err := s.ignoreRepo.Exists(ctx, toID, fromID)
	if err != nil {
		return nil, fmt.Errorf("checking ignore: %w", err)
	}
	if ignored {
		return nil, fmt.Errorf("player is ignoring you")
	}
	return &MessageResult{
		FromCharacterID:   fromID,
		FromCharacterName: fromChar.Name,
		ToCharacterIDs:    []int{toID},
		Channel:           "whisper",
		Message:           message,
		Type:              "whisper",
	}, nil
}

func (s *chatService) SendEmote(ctx context.Context, charID int, action string) (*MessageResult, error) {
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	return &MessageResult{
		FromCharacterID:   charID,
		FromCharacterName: char.Name,
		ToCharacterIDs:    []int{},
		Channel:           "emote",
		Message:           action,
		Type:              "emote",
	}, nil
}

func (s *chatService) SendChannel(ctx context.Context, channel, message string, charID int) (*MessageResult, error) {
	channel = strings.ToLower(channel)
	validChannels := DefaultChannels()
	isValid := false
	for _, c := range validChannels {
		if c == channel {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("invalid channel: %s", channel)
	}
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	sub, err := s.subRepo.GetByCharacterAndChannel(ctx, charID, channel)
	if err != nil {
		// subscription doesn't exist, create it
		_, err = s.subRepo.Create(ctx, charID, channel, "white")
		if err != nil {
			return nil, fmt.Errorf("creating subscription: %w", err)
		}
	} else if !sub.Enabled {
		return nil, fmt.Errorf("channel %s is disabled", channel)
	}
	return &MessageResult{
		FromCharacterID:   charID,
		FromCharacterName: char.Name,
		ToCharacterIDs:    []int{},
		Channel:           channel,
		Message:           message,
		Type:              "channel",
	}, nil
}

func (s *chatService) GetChannels(charID int) ([]ChannelState, error) {
	subs, err := s.subRepo.ListByCharacter(context.Background(), charID)
	if err != nil {
		return nil, err
	}
	states := make([]ChannelState, 0, len(subs))
	for _, sub := range subs {
		states = append(states, ChannelState{
			Name:    sub.Channel,
			Enabled: sub.Enabled,
			Color:   sub.Color,
		})
	}
	// fill in defaults for channels not subscribed
	defaults := DefaultChannels()
	for _, ch := range defaults {
		found := false
		for _, sub := range subs {
			if sub.Channel == ch {
				found = true
				break
			}
		}
		if !found {
			states = append(states, ChannelState{
				Name:    ch,
				Enabled: true,
				Color:   "white",
			})
		}
	}
	return states, nil
}

func (s *chatService) SetChannelEnabled(ctx context.Context, charID int, channel string, enabled bool) error {
	sub, err := s.subRepo.GetByCharacterAndChannel(ctx, charID, channel)
	if err != nil {
		// create subscription
		color := "white"
		_, err = s.subRepo.Create(ctx, charID, channel, color)
		return err
	}
	return s.subRepo.UpdateEnabled(ctx, sub.ID, enabled)
}

func (s *chatService) SetChannelColor(ctx context.Context, charID int, channel string, color string) error {
	sub, err := s.subRepo.GetByCharacterAndChannel(ctx, charID, channel)
	if err != nil {
		// create subscription
		_, err = s.subRepo.Create(ctx, charID, channel, color)
		return err
	}
	return s.subRepo.UpdateColor(ctx, sub.ID, color)
}

func (s *chatService) IgnorePlayer(ctx context.Context, charID, ignoredID int) error {
	if charID == ignoredID {
		return fmt.Errorf("cannot ignore yourself")
	}
	exists, err := s.ignoreRepo.Exists(ctx, charID, ignoredID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("player already ignored")
	}
	_, err = s.ignoreRepo.Create(ctx, charID, ignoredID)
	return err
}

func (s *chatService) UnignorePlayer(ctx context.Context, charID, ignoredID int) error {
	return s.ignoreRepo.DeleteByCharacterAndIgnored(ctx, charID, ignoredID)
}

func (s *chatService) GetIgnoredPlayers(ctx context.Context, charID int) ([]int, error) {
	entries, err := s.ignoreRepo.ListByCharacter(ctx, charID)
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(entries))
	for i, e := range entries {
		ids[i] = e.IgnoredID
	}
	return ids, nil
}

func (s *chatService) QueueOfflineTell(ctx context.Context, fromID int, recipientName string, message string) error {
	fromChar, err := s.charRepo.Get(ctx, fromID)
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}
	// check if recipient exists and is online
	toChar, err := s.charRepo.GetByName(ctx, recipientName)
	if err == nil && toChar != nil {
		// character exists - check if online via lastSeenAt
		// for now, we'll queue regardless; the client can check
		_ = toChar
	}
	_, err = s.tellRepo.Create(ctx, fromID, recipientName, message)
	if err != nil {
		return fmt.Errorf("queueing tell: %w", err)
	}
	_ = fromChar // used for validation
	return nil
}

func (s *chatService) DeliverQueuedTells(ctx context.Context, charID int) ([]QueuedTell, error) {
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	tells, err := s.tellRepo.ListByRecipientName(ctx, char.Name)
	if err != nil {
		return nil, fmt.Errorf("fetching queued tells: %w", err)
	}
	result := make([]QueuedTell, 0, len(tells))
	for _, t := range tells {
		fromChar, _ := s.charRepo.Get(ctx, t.FromID)
		fromName := "unknown"
		if fromChar != nil {
			fromName = fromChar.Name
		}
		result = append(result, QueuedTell{
			ID:            t.ID,
			FromID:        t.FromID,
			FromName:      fromName,
			RecipientName: t.RecipientName,
			Message:       t.Message,
			QueuedAt:      t.QueuedAt.Format(time.RFC3339),
		})
	}
	// delete delivered tells
	for _, t := range tells {
		_ = s.tellRepo.Delete(ctx, t.ID)
	}
	return result, nil
}
