package handlers

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/m4hi2/MeterAlertBot/internal/database/models"
	"github.com/m4hi2/MeterAlertBot/internal/database/repo"
	"github.com/m4hi2/MeterAlertBot/internal/tgbot/state"
	tele "gopkg.in/telebot.v3"
)

type Handlers struct {
	state        *state.Store
	userRepo     repo.UserRepository
	meterRepo    repo.MeterRepository
	providerRepo repo.ProviderRepository
}

func New(
	st *state.Store,
	userRepo repo.UserRepository,
	meterRepo repo.MeterRepository,
	providerRepo repo.ProviderRepository,
) *Handlers {
	return &Handlers{
		state:        st,
		userRepo:     userRepo,
		meterRepo:    meterRepo,
		providerRepo: providerRepo,
	}
}

func (h *Handlers) getOrCreateUser(ctx context.Context, sender *tele.User) (*models.User, error) {
	platformID := strconv.FormatInt(sender.ID, 10)
	user, err := h.userRepo.GetByPlatformID(ctx, models.PlatformTelegram, platformID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	user = &models.User{
		Platform:   models.PlatformTelegram,
		PlatformID: platformID,
	}
	return user, h.userRepo.Create(ctx, user)
}
