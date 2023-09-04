package config

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	Env     = "ENV"
	Port    = "PORT"
	Mocking = "MOCKING"

	SqliteDBPath = "SQLITE_DB_PATH"
)

type App struct {
	Port    string
	Env     Environment
	Mocking bool

	SqliteDBPath string

	FrontendDistFS fs.FS
}

func (cfg *App) MustValidate() {
	var issues []string

	if cfg.Port == "" {
		issues = append(issues, "PORT is required")
	}

	if cfg.SqliteDBPath == "" {
		issues = append(issues, "SQLITE_DB_PATH is required")
	}

	if cfg.FrontendDistFS == nil {
		issues = append(issues, "FrontendDistFS is required")
	}

	if len(issues) > 0 {
		log.Fatalf("invalid app config: %s", strings.Join(issues, ", "))
	}
}

func NewFromEnv() *App {
	env, err := parseEnvironment(os.Getenv(Env))
	if err != nil {
		log.Fatalf("unknown environment set: %v", err)
	}

	return &App{
		Port:    os.Getenv(Port),
		Env:     env,
		Mocking: os.Getenv(Mocking) == "true",

		SqliteDBPath: os.Getenv(SqliteDBPath),
	}
}

func NewTestConfig() *App {
	cfgPath := ".test"

	cfg, err := loadConfig(os.DirFS("./config"), cfgPath)
	if err != nil {
		log.Fatalf("error loading test app config file: %v", err)
	}

	return cfg
}

func NewFromCustomConfig(cfgPath string) *App {
	absPath, err := filepath.Abs(cfgPath)
	if err != nil {
		log.Fatalf("error resolving app config path: %v", err)
	}

	cfg, err := loadConfig(os.DirFS(filepath.Dir(absPath)), filepath.Base(absPath))
	if err != nil {
		log.Fatalf("error loading custom app config file: %v", err)
	}

	return cfg
}

func loadConfig(fsys fs.FS, filePath string) (*App, error) {
	file, err := fsys.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening config dotfile file: %w", err)
	}
	defer file.Close()

	var config App

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			// Skip comments and empty lines
			continue
		}

		parts := strings.SplitN(line, "=", 2) //nolint:gomnd
		if len(parts) != 2 {                  //nolint:gomnd
			return nil, fmt.Errorf("invalid line in .env file: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case Env:
			env, err := parseEnvironment(value)
			if err != nil {
				return nil, fmt.Errorf("unknown environment in config dotfile file: %s = %s", key, value)
			}

			config.Env = env
		case Port:
			config.Port = value
		case Mocking:
			config.Mocking = value == "true"
		case SqliteDBPath:
			config.SqliteDBPath = value
		default:
			return nil, fmt.Errorf("unknown key in config dotfile file: %s", key)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading config dotfile file: %w", err)
	}

	return &config, nil
}
