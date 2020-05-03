package appcron

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/alrusov/log"
	"github.com/alrusov/misc"
	"github.com/alrusov/panic"
)

//----------------------------------------------------------------------------------------------------------------------------//

var (
	core   *cron.Cron
	loc    *time.Location
	parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

//----------------------------------------------------------------------------------------------------------------------------//

// Init --
func Init(tzFile string, location string) (err error) {
	cron.DefaultLogger = &log.CronLog{}

	if tzFile != "" {
		zi, err := misc.AbsPath(tzFile)
		if err == nil {
			os.Setenv("ZONEINFO", zi)
		}
	}

	loc, err = time.LoadLocation(location)
	if err != nil {
		log.Message(log.ERR, `Cron: load location error (%s). Try to use UTC`, err.Error())
		location = "UTC"
		loc, err = time.LoadLocation(location)
		if err != nil {
			return
		}
	}
	core = cron.New(
		cron.WithLocation(loc),
	)

	go func() {
		defer panic.SaveStackToLog()
		core.Run()
	}()

	log.Message(log.INFO, `Cron started`)

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Parse --
func Parse(spec string) (cron.Schedule, error) {
	return parser.Parse("TZ=" + loc.String() + " " + spec)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Add --
func Add(sched string, handler cron.Job) (cron.EntryID, error) {
	return core.AddJob(sched, handler)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Remove --
func Remove(id cron.EntryID) {
	if id > 0 {
		core.Remove(id)
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

// Entry --
func Entry(id cron.EntryID) cron.Entry {
	return core.Entry(id)
}

//----------------------------------------------------------------------------------------------------------------------------//
