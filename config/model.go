package config

import (
)

type ConfigModel struct {
	Server   		ServerConfig    `yaml:"Server"`
	Postgres 		PostgresConfig  `yaml:"Postgres"`
	StorageType 	string			`env:"STORAGE_TYPE"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"DBName"`
	SSLMode  string `yaml:"sslMode"`
}

type ServerConfig struct {
	Host       string `yaml:"host" validate:"required"`
	GRPCPort       string `yaml:"gRPCport" validate:"required"`
	HTTPPort       string `yaml:"HTTPport" validate:"required"`
}
