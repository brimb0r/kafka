package app_contain

import (
	"awesomeProject/internal/configuration"
	"awesomeProject/internal/dalmatian"
	terrierRepo "awesomeProject/internal/repo/terrier"
	internalTrans "awesomeProject/internal/translator/terrier_translator"
	"awesomeProject/pkg/bark"
	"awesomeProject/pkg/sleep"
	"awesomeProject/pkg/translator"
	"context"
	"github.com/robfig/cron/v3"
	"log"
)

var (
	d = &dalmatian.Dalmatian{}
)

func StartEngine(config *configuration.Configuration) cron.Job {
	r := &runner{
		bark: &bark.Bark{
			DogName: "Cedar",
			Runner:  d.DalmatianRunner(),
		},
		sleep:   &sleep.Sleep{},
		terrier: &terrierRepo.Repo{Database: config.Mongo.MongoDatabase()},
	}
	return r
}

type runner struct {
	bark    bark.IBark
	sleep   sleep.ISleep
	terrier terrierRepo.ITerrierRepo
}

func (app runner) Run() {
	terrierQuery, err := app.terrier.QueryTerriers(context.Background())
	terrierChan := make(chan translator.ITranslator, len(terrierQuery))
	for _, t := range terrierQuery {
		terrierChan <- &internalTrans.TerrierTranslator{
			Terrier: t,
			Repo:    app.terrier,
		}
	}

	if err != nil {
		log.Printf("Error %v", err)
	}

	err = app.bark.Bark()
	if err != nil {
		return
	}
	close(terrierChan)

	err = app.sleep.Sleep(terrierChan)
	if err != nil {
		return
	}
}
