package main

import (
	"net/http"

	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/database"
	transportHTTP "github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/transport/http"
	log "github.com/sirupsen/logrus"
)

// App - contains applicaiton information
type App struct {
	Name    string
	Version string
}

// Run - sets up our application
func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
		}).Info("Setting Up Our APP")

	var err error
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	err = database.MigrateDB(db)
	if err != nil {
		log.Error("failed to setup database")
		return err
	}

	commentService := comment.NewService(db)
	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		log.Error("Failed to set up server")
		return err
	}

	log.Info("App startup successful")
	return nil
}

func main() {
	app := App{
		Name:    "commenting service",
		Version: "1.0.0",
	}
	if err := app.Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}

}
