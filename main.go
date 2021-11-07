package main

import (
	"fmt"
	"net/http"

	"github.com/scevallos/WheelOfFortuneEveryday/pkg/roku"
	"github.com/scevallos/WheelOfFortuneEveryday/pkg/web"
)

func main() {
	client, err := roku.NewClient()
	panicIfErr(err)

	service, err := web.NewService(client)
	panicIfErr(err)

	fmt.Println("Starting HTTP server")
	if err := http.ListenAndServe(":8787", service.GetRouter()); err != nil {
		panic(err)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
