package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"simple-micro/core/logger"
	sample_services "simple-micro/exmsg/services"
	"strings"
)

func (s *Server) GetNumber(ctx context.Context, req *sample_services.SampleRequest) (*sample_services.SampleResponse, error) {
	lg := logger.NewLogger()
	lg.Infof("Received id: %d", req.GetId())
	dbNames := make([]string, 0)
	var err error
	if s.MongoDb != nil {
		client := s.MongoDb.GetClient()
		dbNames, err = client.ListDatabaseNames(ctx, bson.D{{}})
		if err != nil {
			lg.Errorw("Failed to list databases", "error", err)
			return nil, err
		}

		coll := client.Database(s.MongoDb.DbName).Collection("attachments")
		var result bson.M
		err = coll.FindOne(context.TODO(), bson.D{{"name", "ae32ddf4c4b848569bfdf562c3b768c3.webp"}}).
			Decode(&result)
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Printf("No document was found with the title %s\n", "title")
			return nil, err
		}
		if err != nil {
			panic(err)
		}
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
	}

	return &sample_services.SampleResponse{
		Status: true,
		Data: map[string]string{
			"number": cast.ToString(req.GetId()),
			"dbs":    strings.Join(dbNames, ", "),
		},
	}, nil
}
