package app

import (
	"context"
)

type App struct {
	logger  Logger
	storage Storage
	limiter Limiter
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateBlacklistSubnet(ctx context.Context, subnet string) error
	DeleteBkacklistSubnet(ctx context.Context, subnet string) error
	CreateWhitelistSubnet(ctx context.Context, subnet string) error
	DeleteWhitelistSubnet(ctx context.Context, subnet string) error
	IsIPInBlacklist(ctx context.Context, subnet string) (bool, error)
	IsIPInWhitelist(ctx context.Context, subnet string) (bool, error)
}

type Limiter interface {
	Add(login, password, ip string) bool
	Reset(login, ip string)
}

func New(logger Logger, storage Storage, limiter Limiter) *App {
	return &App{
		logger:  logger,
		storage: storage,
		limiter: limiter,
	}
}

func (a *App) Login(ctx context.Context, login, password, ip string) (bool, error) {
	isInBlacklist, err := a.storage.IsIPInBlacklist(ctx, ip)
	if err != nil {
		return false, err
	}
	if isInBlacklist {
		return false, nil
	}

	isInWhitelist, err := a.storage.IsIPInWhitelist(ctx, ip)
	if err != nil {
		return false, err
	}
	if isInWhitelist {
		return true, nil
	}

	res := a.limiter.Add(login, password, ip)

	return res, nil
}

func (a *App) ResetBuket(ctx context.Context, login, ip string) error {
	a.limiter.Reset(login, ip)

	return nil
}

func (a *App) AddToBlacklist(ctx context.Context, subnet string) error {
	err := a.storage.CreateBlacklistSubnet(ctx, subnet)

	return err
}

func (a *App) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	err := a.storage.DeleteBkacklistSubnet(ctx, subnet)

	return err
}

func (a *App) AddToWhitelist(ctx context.Context, subnet string) error {
	err := a.storage.CreateWhitelistSubnet(ctx, subnet)

	return err
}

func (a *App) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	err := a.storage.DeleteWhitelistSubnet(ctx, subnet)

	return err
}
