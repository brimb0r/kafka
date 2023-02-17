package main

import (
	"awesomeProject/cmd/project/app_contain"
	"awesomeProject/internal/configuration"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"
)

func main() {
	log.Printf("Setting Config")
	os.Setenv("CONFIG_PATH", "C:\\Users\\gagno\\GolandProjects\\awesomeProject\\config_files")
	c := configuration.Configure()
	log.Printf("Starting")
	dogRunner := app_contain.StartEngine(c)
	runScheduler(c.Schedule, dogRunner)
}

func runScheduler(chronSchedule string, job cron.Job) {
	newScheduler := gocron.NewScheduler(time.UTC)
	_, err := newScheduler.Cron(chronSchedule).Do(job.Run)
	if err != nil {
		fmt.Println(err.Error())
	}
	newScheduler.StartBlocking()
}
