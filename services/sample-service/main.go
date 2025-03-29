package main

import (
	app2 "simple-micro/core/app"
	sample_services "simple-micro/exmsg/services"
	"simple-micro/pkg"
	"simple-micro/services/sample-service/handlers"
)

func main() {
	app := app2.App{
		Name: "Sample service",
		Code: pkg.SampleClient,
		Type: app2.ServiceType,
		Port: 50051,
	}

	app.LoadConfig()
	s := app.NewGrpcSever()

	err := app.InitDatabase("mongo")
	if err != nil {
		panic(err)
	}

	sample_services.RegisterSampleServer(s, &handlers.Server{
		MongoDb: app.Dbs.Mongo,
	})
	app.NewServer(s)
}
