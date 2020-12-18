package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/PedPet/breeder/config"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	grpcAddr := os.Getenv("PORT")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "breeder",
			"time:", log.DefaultTimestampUTC,
			"called", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	settings, err := config.LoadSettings()
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit((-1))
	}

	// Instantiate database connection
	var db *sql.DB
	{
		var err error

		dbSource := settings.DB.User + ":" + settings.DB.Password +
			"@tcp(" + settings.DB.Host + ")/" + settings.DB.Database
		db, err = sql.Open("mysql", dbSource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
        }
        
        err = db.Ping()
        if err != nil {
            level.Error(logger).Log("exit", err)
        }
    }
    
    ctx := context.Background()

    // Instantiate service 
    var srv service.Breeder 
    {

    }
}
