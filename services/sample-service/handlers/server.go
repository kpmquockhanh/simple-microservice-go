package handlers

import (
	"simple-micro/core/db"
	sample_services "simple-micro/exmsg/services"
)

type Server struct {
	sample_services.UnimplementedSampleServer
	MongoDb *db.MongoDatabase
}
