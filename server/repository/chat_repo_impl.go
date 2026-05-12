package repository

import (
	"context"
	"time"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/characterchannel"
	"herbst-server/db/characterignore"
	"herbst-server/db/tellqueue"
)

// ChannelToRepo maps CharacterChannel fields to ChannelSubscription.
func channelToSub(id int, channel string, enabled bool, color string) *ChannelSubscription {
	return &ChannelSubscription{ID: id, Channel: channel, Enabled: enabled, Color: color}
}

// entChannelSubscriptionRepo implements ChannelSubscriptionRepo using ent.
type entChannelSubscriptionRepo struct {
	client *db.Client
}

func NewEntChannelSubscriptionRepo(client *db.Client) ChannelSubscriptionRepo {
	return &entChannelSubscriptionRepo{client: client}
}

func (r *entChannelSubscriptionRepo) ListByCharacter(ctx context.Context, charID int) ([]*ChannelSubscription, error) {
	cc, err := r.client.CharacterChannel.Query().
		Where(characterchannel.HasCharacterWith(character.ID(charID))).
		Only(ctx)
	if err != nil {
		return []*ChannelSubscription{}, nil
	}
	subs := make([]*ChannelSubscription, 0, 5)
	subs = append(subs, channelToSub(cc.ID, "chat", cc.ChatEnabled, cc.ChatColor))
	subs = append(subs, channelToSub(cc.ID, "newbie", cc.NewbieEnabled, cc.NewbieColor))
	subs = append(subs, channelToSub(cc.ID, "trade", cc.TradeEnabled, cc.TradeColor))
	subs = append(subs, channelToSub(cc.ID, "clan", cc.ClanEnabled, cc.ClanColor))
	subs = append(subs, channelToSub(cc.ID, "auction", cc.AuctionEnabled, ""))
	return subs, nil
}

func (r *entChannelSubscriptionRepo) GetByCharacterAndChannel(ctx context.Context, charID int, channel string) (*ChannelSubscription, error) {
	cc, err := r.client.CharacterChannel.Query().
		Where(characterchannel.HasCharacterWith(character.ID(charID))).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	switch channel {
	case "chat":
		return channelToSub(cc.ID, "chat", cc.ChatEnabled, cc.ChatColor), nil
	case "newbie":
		return channelToSub(cc.ID, "newbie", cc.NewbieEnabled, cc.NewbieColor), nil
	case "trade":
		return channelToSub(cc.ID, "trade", cc.TradeEnabled, cc.TradeColor), nil
	case "clan":
		return channelToSub(cc.ID, "clan", cc.ClanEnabled, cc.ClanColor), nil
	case "auction":
		return channelToSub(cc.ID, "auction", cc.AuctionEnabled, ""), nil
	default:
		return nil, nil
	}
}

func (r *entChannelSubscriptionRepo) Create(ctx context.Context, charID int, channel, color string) (*ChannelSubscription, error) {
	char, err := r.client.Character.Get(ctx, charID)
	if err != nil {
		return nil, err
	}
	cc, err := r.client.CharacterChannel.Create().
		SetCharacter(char).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return channelToSub(cc.ID, channel, true, color), nil
}

func (r *entChannelSubscriptionRepo) UpdateEnabled(ctx context.Context, id int, enabled bool) error {
	cc, err := r.client.CharacterChannel.Get(ctx, id)
	if err != nil {
		return err
	}
	cc.ChatEnabled = enabled
	_, err = r.client.CharacterChannel.UpdateOne(cc).Save(ctx)
	return err
}

func (r *entChannelSubscriptionRepo) UpdateColor(ctx context.Context, id int, color string) error {
	cc, err := r.client.CharacterChannel.Get(ctx, id)
	if err != nil {
		return err
	}
	cc.ChatColor = color
	_, err = r.client.CharacterChannel.UpdateOne(cc).Save(ctx)
	return err
}

func (r *entChannelSubscriptionRepo) Delete(ctx context.Context, id int) error {
	return r.client.CharacterChannel.DeleteOneID(id).Exec(ctx)
}

func (r *entChannelSubscriptionRepo) DeleteByCharacterAndChannel(ctx context.Context, charID int, channel string) error {
	cc, err := r.client.CharacterChannel.Query().
		Where(characterchannel.HasCharacterWith(character.ID(charID))).
		Only(ctx)
	if err != nil {
		return nil
	}
	switch channel {
	case "chat":
		cc.ChatEnabled = false
	case "newbie":
		cc.NewbieEnabled = false
	case "trade":
		cc.TradeEnabled = false
	case "clan":
		cc.ClanEnabled = false
	case "auction":
		cc.AuctionEnabled = false
	}
	_, err = r.client.CharacterChannel.UpdateOne(cc).Save(ctx)
	return err
}

// entOfflineTellRepo implements OfflineTellRepo using ent.
type entOfflineTellRepo struct {
	client *db.Client
}

func NewEntOfflineTellRepo(client *db.Client) OfflineTellRepo {
	return &entOfflineTellRepo{client: client}
}

func (r *entOfflineTellRepo) ListByRecipient(ctx context.Context, recipientID int) ([]*OfflineTell, error) {
	tells, err := r.client.TellQueue.Query().
		Where(tellqueue.HasRecipientWith(character.ID(recipientID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*OfflineTell, len(tells))
	for i, t := range tells {
		result[i] = &OfflineTell{
			ID:            t.ID,
			FromID:        t.SenderId,
			RecipientID:   recipientID,
			RecipientName: t.SenderName,
			Message:       t.Message,
			QueuedAt:       t.SentAt,
		}
	}
	return result, nil
}

func (r *entOfflineTellRepo) ListByRecipientName(ctx context.Context, recipientName string) ([]*OfflineTell, error) {
	char, err := r.client.Character.Query().
		Where(character.NameEQ(recipientName)).
		Only(ctx)
	if err != nil {
		return []*OfflineTell{}, nil
	}
	return r.ListByRecipient(ctx, char.ID)
}

func (r *entOfflineTellRepo) Create(ctx context.Context, fromID int, recipientName string, message string) (*OfflineTell, error) {
	fromChar, err := r.client.Character.Get(ctx, fromID)
	if err != nil {
		return nil, err
	}
	recipient, err := r.client.Character.Query().
		Where(character.NameEQ(recipientName)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	tell, err := r.client.TellQueue.Create().
		SetSenderId(fromID).
		SetSenderName(fromChar.Name).
		SetMessage(message).
		SetExpiresAt(expiresAt).
		SetRecipient(recipient).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &OfflineTell{
		ID:            tell.ID,
		FromID:        tell.SenderId,
		RecipientID:   recipient.ID,
		RecipientName: recipientName,
		Message:       tell.Message,
		QueuedAt:      tell.SentAt,
	}, nil
}

func (r *entOfflineTellRepo) Delete(ctx context.Context, id int) error {
	return r.client.TellQueue.DeleteOneID(id).Exec(ctx)
}

func (r *entOfflineTellRepo) DeleteByRecipient(ctx context.Context, recipientID int) error {
	tells, err := r.client.TellQueue.Query().
		Where(tellqueue.HasRecipientWith(character.ID(recipientID))).
		All(ctx)
	if err != nil {
		return err
	}
	for _, t := range tells {
		_ = r.client.TellQueue.DeleteOne(t).Exec(ctx)
	}
	return nil
}

// entIgnoreRepo implements IgnoreRepo using ent.
type entIgnoreRepo struct {
	client *db.Client
}

func NewEntIgnoreRepo(client *db.Client) IgnoreRepo {
	return &entIgnoreRepo{client: client}
}

func (r *entIgnoreRepo) ListByCharacter(ctx context.Context, charID int) ([]*IgnoreEntry, error) {
	ignores, err := r.client.CharacterIgnore.Query().
		Where(characterignore.HasIgnorerWith(character.ID(charID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*IgnoreEntry, len(ignores))
	for i, ig := range ignores {
		result[i] = &IgnoreEntry{
			ID:          ig.ID,
			CharacterID: charID,
			IgnoredID:   ig.IgnoredCharacterId,
			CreatedAt:   ig.IgnoredAt,
		}
	}
	return result, nil
}

func (r *entIgnoreRepo) Exists(ctx context.Context, charID, ignoredID int) (bool, error) {
	count, err := r.client.CharacterIgnore.Query().
		Where(
			characterignore.HasIgnorerWith(character.ID(charID)),
			characterignore.IgnoredCharacterId(ignoredID),
		).
		Count(ctx)
	return count > 0, err
}

func (r *entIgnoreRepo) Create(ctx context.Context, charID, ignoredID int) (*IgnoreEntry, error) {
	char, err := r.client.Character.Get(ctx, charID)
	if err != nil {
		return nil, err
	}
	ignore, err := r.client.CharacterIgnore.Create().
		SetIgnoredCharacterId(ignoredID).
		SetIgnorer(char).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &IgnoreEntry{
		ID:          ignore.ID,
		CharacterID: charID,
		IgnoredID:   ignoredID,
		CreatedAt:   ignore.IgnoredAt,
	}, nil
}

func (r *entIgnoreRepo) Delete(ctx context.Context, id int) error {
	return r.client.CharacterIgnore.DeleteOneID(id).Exec(ctx)
}

func (r *entIgnoreRepo) DeleteByCharacterAndIgnored(ctx context.Context, charID, ignoredID int) error {
	ignores, err := r.client.CharacterIgnore.Query().
		Where(
			characterignore.HasIgnorerWith(character.ID(charID)),
			characterignore.IgnoredCharacterId(ignoredID),
		).
		All(ctx)
	if err != nil {
		return err
	}
	for _, ig := range ignores {
		_ = r.client.CharacterIgnore.DeleteOne(ig).Exec(ctx)
	}
	return nil
}
