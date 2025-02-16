package integration

import (
	"avito-shop/config"
	"avito-shop/internal/app"
)

func setupTestApp() (*app.App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	if err := cfg.InitDB(); err != nil {
		return nil, err
	}

	return app.New(cfg)
}
