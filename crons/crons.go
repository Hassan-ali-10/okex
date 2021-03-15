package crons

import (
	
	"github.com/mileusna/crontab"
	"log"
	// "fmt"
	 buy_controllers "okex/crons/methods/buy"
	
)

func AllCrons() {
	
	ctab := crontab.New() // create cron tab
	// AddJob and test the errors
    err := ctab.AddJob("0 12 1 * *",buy_controllers.MyFunc) // on 1st day of month
    if err != nil {
        log.Println(err)
        return
    }
    // MustAddJob is like AddJob but panics on wrong syntax or problems with func/args
    // This aproach is similar to regexp.Compile and regexp.MustCompile from go's standard library,  used for easier initialization on startup
    ctab.MustAddJob("* * * * *", buy_controllers.MyFunc) // every minute
    ctab.MustAddJob("0 12 * * *", buy_controllers.MyFunc3) // noon lauch

    // fn with args
    ctab.MustAddJob("0 0 * * 1,2", buy_controllers.MyFunc2, "Monday and Tuesday midnight", 123) 
    ctab.MustAddJob("*/5 * * * *", buy_controllers.MyFunc2, "every five min", 0)

    // all your other app code as usual, or put sleep timer for demo
    // time.Sleep(10 * time.Minute)    
	
}