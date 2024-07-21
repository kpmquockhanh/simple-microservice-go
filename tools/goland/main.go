package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simple-micro/tools/goland/configuration"
)

func main() {
	currentDir, _ := os.Getwd()

	pathProject := fmt.Sprintf("%v/.idea", currentDir)

	if _, err := os.Stat(pathProject); os.IsNotExist(err) {
		log.Fatal("Directory a has never been opened by goland")
		return
	}

	listGenTypes := []string{"apis", "services"}
	for _, genType := range listGenTypes {
		err := filepath.Walk(fmt.Sprintf("%v/%s", currentDir, genType), func(path string, info os.FileInfo, err error) error {
			if info.Name() == "main.go" {
				RunCommand(currentDir, filepath.Dir(path))
			}
			return nil
		})

		if err != nil {
			log.Fatal(fmt.Sprintf("Error %v", err))
		}
	}
}

func RunCommand(currentDir, dir string) {
	log.Printf("Processing: %v", dir)

	file, err := os.ReadFile(dir + "/cicd.json")
	if err != nil {
		log.Printf(fmt.Sprintf("Error when get file, ignore %v", err))
		return
	}

	cicdConfig := configuration.Cicd{}

	err = json.Unmarshal(file, &cicdConfig)
	path := currentDir + "/.idea/runConfigurations"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("Error when create dir %v", err)
			return
		}
	}
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error when Stat %v", err)
	}

	nameFile := fmt.Sprintf("%v/%v.xml", path, cicdConfig.Deploy.ServiceName)

	writer, err := os.OpenFile(nameFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		println(fmt.Sprintf("Error OpenFile %v", err))
		return
	}
	rawOutput := configuration.New(cicdConfig, filepath.Base(currentDir))
	output, err := xml.MarshalIndent(rawOutput, "", "    ")
	if err != nil {
		println(fmt.Sprintf("Error MarshalIndent %v", err))
		return
	}

	_, err = writer.Write(output)
	if err != nil {
		println(fmt.Sprintf("Error Write %v", err))
		return
	}
	log.Printf("Process done: %v", dir)
}
