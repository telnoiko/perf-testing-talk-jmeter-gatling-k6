package main

import (
	"task-service/app/rest"
	"task-service/app/store"
)

func main() {

	store, err := store.New()
	defer store.Stop()
	if err != nil {
		panic(err)
	}

	auth := rest.NewAuth(store)

	rest := rest.New(auth, store)

	rest.Start()
}
