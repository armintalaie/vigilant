// cmd/desktop/main.go
package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"log"
	"sync"
	"time"
	pb "vigilant/internal/logger"
)

var assets embed.FS

// App struct holds all our application components
type App struct {
	ctx       context.Context
	wg        sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
	db        *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db, err := sql.Open("sqlite3", "logs.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	} else {
		a.db = db

	}
}

// StartServices starts both gRPC and Kafka services
func (a *App) StartServices() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isRunning {
		return "Services are already running"
	}

	a.isRunning = true
	return "Services started successfully"
}

// StopServices gracefully stops all services
func (a *App) StopServices() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isRunning {
		return "Services are not running"
	}

	a.isRunning = false
	return "Services stopped successfully"
}

// GetServiceStatus returns the current status of services
func (a *App) GetServiceStatus() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isRunning {
		return "Services are running"
	}
	return "Services are stopped"
}

// shutdown is called when the app closes
func (a *App) shutdown(ctx context.Context) {
	a.StopServices()
	a.wg.Wait()
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Vigilant",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) GetLogs() ([]*pb.Log, error) {
	rows, err := a.db.Query(`SELECT id, message, timestamp, level, severity, source, "group", tags, type, origin, data FROM logs`)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %v", err)
	}
	defer rows.Close()

	var logs []*pb.Log
	for rows.Next() {
		var logMessage pb.Log
		var dataJSON string

		var timestamp time.Time // Temporary variable to hold the timestamp
		var level string        // Temporary variable to hold the level

		if err := rows.Scan(
			&logMessage.Id,
			&logMessage.Message,
			&timestamp, // Scan into timestamp variable first
			&level,     // Scan into level variable first
			&logMessage.Severity,
			&logMessage.Source,
			&logMessage.Group,
			&logMessage.Tags,
			&logMessage.Type,
			&logMessage.Origin,
			&dataJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Convert timestamp to Unix timestamp (seconds since epoch)
		logMessage.Timestamp = timestamp.Unix()
		logMessage.Level = logTextToLevel(level)

		// Unmarshal the Data JSON string to a map
		if err := json.Unmarshal([]byte(dataJSON), &logMessage.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %v", err)
		}

		logs = append(logs, &logMessage)
	}

	return logs, err
}

func logTextToLevel(level string) pb.LogLevel {
	switch level {
	case "DEBUG":
		return pb.LogLevel_DEBUG
	case "INFO":
		return pb.LogLevel_INFO
	case "WARN":
		return pb.LogLevel_WARN
	case "ERROR":
		return pb.LogLevel_ERROR
	case "NONE":
		return pb.LogLevel_NONE
	default:
		return pb.LogLevel_ALL
	}
}
