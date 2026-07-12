package routes

import (
	"context"
	"fmt"
	"log/slog"

	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// applyDialogEffects runs a list of effect IDs attached to a dialog response or node entry.
func applyDialogEffects(ctx context.Context, wsc *WSConn, repos *repository.Container, client *db.Client, effectIDs []int) {
	for _, id := range effectIDs {
		if err := applyDialogEffect(ctx, wsc, repos, client, id); err != nil {
			dblog.Error("dialog effect failed", err, slog.Int("effect_id", id), slog.Int("character_id", wsc.CharacterID))
		}
	}
}

// applyDialogEffect applies a single data-driven effect to the current player.
func applyDialogEffect(ctx context.Context, wsc *WSConn, repos *repository.Container, client *db.Client, effectID int) error {
	eff, err := repos.Effect.Get(ctx, effectID)
	if err != nil {
		return err
	}

	switch eff.EffectType {
	case "change_race":
		raceName, _ := eff.Parameters["race_name"].(string)
		if raceName == "" {
			return fmt.Errorf("change_race effect %d missing race_name", effectID)
		}
		_, err := repos.Character.Update(ctx, wsc.CharacterID, repository.CharacterUpdates{Race: &raceName})
		if err != nil {
			return err
		}
		sendNotification(wsc, fmt.Sprintf("Your race has changed to %s.", raceName))
		slog.Info("race changed via dialog", slog.Int("character_id", wsc.CharacterID), slog.String("race", raceName), slog.String("effect", eff.Name))
	case "change_class":
		className, _ := eff.Parameters["class_name"].(string)
		if className == "" {
			return fmt.Errorf("change_class effect %d missing class_name", effectID)
		}
		// Validate class against DB: faction in "class" category for the character's world.
		classFaction, err := getClassFactionByName(ctx, client, className, wsc.World)
		if err != nil || classFaction == nil {
			return fmt.Errorf("change_class effect %d has invalid class %q", effectID, className)
		}
		// Get first specialty ID from the faction's specialties JSON field.
		newSpecialty := "generalist"
		if len(classFaction.Specialties) > 0 {
			newSpecialty = classFaction.Specialties[0].ID
		}
		_, err = repos.Character.Update(ctx, wsc.CharacterID, repository.CharacterUpdates{Class: &className, Specialty: &newSpecialty})
		if err != nil {
			return err
		}
		sendNotification(wsc, fmt.Sprintf("Your class has changed to %s.", className))
		slog.Info("class changed via dialog", slog.Int("character_id", wsc.CharacterID), slog.String("class", className), slog.String("effect", eff.Name))
	case "message":
		text, _ := eff.Parameters["text"].(string)
		if text != "" {
			sendNotification(wsc, text)
		}
	}
	return nil
}