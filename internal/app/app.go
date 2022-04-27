package app

import (
	"context"
)

type App struct {
	logger  Logger
	storage Storage
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

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
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

	return false, nil
}

func (a *App) ResetBuket(ctx context.Context, login, ip string) error {
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
