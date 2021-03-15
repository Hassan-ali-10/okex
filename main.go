package main

import (
	
	routes "okex/routes"
	"sync"
	 "fmt"
	 "net/http"
     crons "okex/crons"
	
)

func main() {
	// links := []string{
	// 	"https://github.com/fabpot",
	// 	"https://github.com/andrew",
	// 	"https://github.com/taylorotwell",
	// 	"https://github.com/egoist",
	// 	"https://github.com/HugoGiraudel",
	// }

	checkUrls("https://github.com/HugoGiraudel")
	crons.AllCrons();
	routes.SetRouters()
	
}

func checkUrls(link string) {
	fmt.Println("creating channel")
	c := make(chan bool)
	var wg sync.WaitGroup

	
		wg.Add(1)   // This tells the waitgroup, that there is now 1 pending operation here
		go checkUrl(link, c, &wg)
	

    // this function literal (also called 'anonymous function' or 'lambda expression' in other languages)
    // is useful because 'go' needs to prefix a function and we can save some space by not declaring a whole new function for this
	go func() {
		fmt.Println("channel is closed")
		wg.Wait()	// this blocks the goroutine until WaitGroup counter is zero
		close(c)    // Channels need to be closed,
	}()    // This calls itself
	fmt.Println("heelo")
    // waits for results to come in through the 'c' channel
	msg := <- c
	fmt.Println("value is assigned to msg variable")
		fmt.Println(msg)
	
}

func checkUrl(url string, c chan bool, wg *sync.WaitGroup) {
	defer (*wg).Done()
	_, err := http.Get(url)

	if err != nil {
		c <- false    // pump the result into the channel
	} else {
		c <- true    // pump the result into the channel
	}
}