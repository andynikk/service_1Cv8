package main

import (
	"Service_1Cv8/internal/cli/forms_exchange"
	"log"
)

func main() {

	c := forms_exchange.NewExchanger()
	if err := c.Run(); err != nil {
		log.Fatal(err.Error())
	}

	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	//<-stop
	//err := c.Shutdown()
	//if err != nil {
	//	log.Println(err.Error())
	//}
}
