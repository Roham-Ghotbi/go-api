package main

import (
	"fmt"

	"github.com/Roham-Ghotbi/go-api/internal/comment"
	"github.com/Roham-Ghotbi/go-api/internal/db"
	transportHttp "github.com/Roham-Ghotbi/go-api/internal/transport/http"
)

// Run - is going to be responsible for instantiation
// and startup of our go application
func Run() error {
	fmt.Println("Starting up our application")

	db, err := db.NewDataBase()

	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}
	if err := db.MigrateDB(); err != nil {
		fmt.Println("failed to migrate the database")
		return err
	}
	cmtService := comment.NewService(db)

	httpHandler := transportHttp.NewHandler(cmtService)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Go REST API Course")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
