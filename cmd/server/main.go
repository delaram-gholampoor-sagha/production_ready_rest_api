package main

import (
	"fmt"
	"net/http"

	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/database"
	transportHTTP "github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/transport/http"
)

// App - the struct which contains things like pointers
// to database connecitons
type App struct {}

// Run - sets up our application
func (app *App) Run() error{
    fmt.Println("Setting Up our app")
	var err error
	db , err  := database.NewDatabse()
	if err != nil {
		return err
	}
	err = database.MigrationDB(db)
	if err != nil {
		return  err
	}

	commentService := comment.NewService(db)
	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()
	if err := http.ListenAndServe(":8081" , handler.Router) ; err != nil {
        fmt.Println("failed to setup server ")
		return err
	}
	return nil 
}

func main() {
   app := App{}
   if err := app.Run(); err != nil {
	   fmt.Println("Error starting up our rest api")
	   fmt.Println(err)
   }

}