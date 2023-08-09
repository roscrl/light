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

type Server struct {
	Port    string
	Env     Environment
	Mocking bool

	SqliteDBPath string
}

func (cfg *Server) MustValid() {
	var issues []string

	if cfg.Env == 0 {
		issues = append(issues, "ENV is required")
	}

	if cfg.Port == "" {
		issues = append(issues, "PORT is required")
	}

	if cfg.SqliteDBPath == "" {
		issues = append(issues, "SQLITE_DB_PATH is required")
	}

	if len(issues) > 0 {
		log.Fatalf("invalid server config: %s", strings.Join(issues, ", "))
	}
}

func FromEnv() *Server {
	env, err := parseEnvironment(os.Getenv(Env))
	if err != nil {
		log.Fatalf("unknown environment in config dotfile file: %v", err)
	}

	return &Server{
		Port:         os.Getenv(Port),
		Env:          env,
		Mocking:      os.Getenv(Mocking) == "true",
		SqliteDBPath: os.Getenv(SqliteDBPath),
	}
}

func TestConfig() *Server {
	cfgPath := ".dev.test"

	cfg, err := loadConfig(os.DirFS("./config"), cfgPath)
	if err != nil {
		log.Fatalf("error loading test server config file: %v", err)
	}

	return cfg
}

func FromCustomConfig(cfgPath string) *Server {
	absPath, err := filepath.Abs(cfgPath)
	if err != nil {
		log.Fatalf("error resolving server config path: %v", err)
	}

	cfg, err := loadConfig(os.DirFS(filepath.Dir(absPath)), filepath.Base(absPath))
	if err != nil {
		log.Fatalf("error loading custom server config file: %v", err)
	}

	return cfg
}

func loadConfig(fsys fs.FS, filePath string) (*Server, error) {
	file, err := fsys.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Server

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
		return nil, err
	}

	return &config, nil
}
