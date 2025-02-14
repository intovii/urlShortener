package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewConfig() (*ConfigModel, error) {
	var cfg ConfigModel
	var found bool
	v := viper.New()
	v.AddConfigPath("config")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		slog.Error("fail to read config", err)
		return &cfg, err
	}
	err = v.Unmarshal(&cfg)
	if err != nil {
		slog.Error("", fmt.Errorf("unable to decode config into struct, %w", err))
		return &cfg, err
	}
	if err := godotenv.Load(); err != nil {
		slog.Error("", fmt.Errorf("unable to get env, %w", err))
		return nil, errors.New("unable to get env")
	}
	cfg.StorageType, found = os.LookupEnv("STORAGE_TYPE")
	if !found {
		slog.Error("", fmt.Errorf("unable to get STORAGE_TYPE from env,"))
		return nil, errors.New("unable to get STORAGE_TYPE from env")
	}
	return &cfg, nil
}
