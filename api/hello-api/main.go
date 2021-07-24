package main

import (
	app2 "simple-micro/core/app"
)

func main() {
	app := app2.App{
		Name: "Hello",
		//Code: "hello",
		Type: app2.APIType,
		BasePath: "/hello",
		Port: 50052,
	}

	app.NewServer(&server{})
}