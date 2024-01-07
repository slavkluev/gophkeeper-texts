package texts

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"texts/internal/domain/models"
)

type Texts struct {
	log          *zap.Logger
	textSaver    TextSaver
	textUpdater  TextUpdater
	textProvider TextProvider
	tokenTTL     time.Duration
}

type TextProvider interface {
	GetAll(ctx context.Context, userUID uint64) ([]models.Text, error)
}

type TextSaver interface {
	SaveText(ctx context.Context, text string, info string, userUID uint64) (uid uint64, err error)
}

type TextUpdater interface {
	UpdateText(ctx context.Context, id uint64, text string, info string, userUID uint64) error
}

func New(
	log *zap.Logger,
	textSaver TextSaver,
	textUpdater TextUpdater,
	textProvider TextProvider,
	tokenTTL time.Duration,
) *Texts {
	return &Texts{
		textSaver:    textSaver,
		textUpdater:  textUpdater,
		textProvider: textProvider,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

func (a *Texts) GetAll(ctx context.Context) ([]models.Text, error) {
	const op = "Texts.GetAll"

	log := a.log.With(
		zap.String("op", op),
	)

	log.Info("attempting to get all texts")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return nil, fmt.Errorf("%s: failed to find user uid", op)
	}

	texts, err := a.textProvider.GetAll(ctx, userUID)
	if err != nil {
		a.log.Error("failed to get all texts", zap.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("texts are got successfully")

	return texts, nil
}

func (a *Texts) SaveText(ctx context.Context, text string, info string) (uint64, error) {
	const op = "Texts.SaveText"

	log := a.log.With(
		zap.String("op", op),
	)

	log.Info("attempting to save text")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return 0, fmt.Errorf("%s: failed to find user uid", op)
	}

	id, err := a.textSaver.SaveText(ctx, text, info, userUID)
	if err != nil {
		log.Error("failed to save text", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("text saved successfully")

	return id, nil
}

func (a *Texts) UpdateText(ctx context.Context, id uint64, text string, info string) error {
	const op = "Texts.UpdateText"

	log := a.log.With(
		zap.String("op", op),
		zap.Uint64("id", id),
	)

	log.Info("attempting to update text")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return fmt.Errorf("%s: failed to find user uid", op)
	}

	err := a.textUpdater.UpdateText(ctx, id, text, info, userUID)
	if err != nil {
		log.Error("failed to update text", zap.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("text updated successfully")

	return nil
}
