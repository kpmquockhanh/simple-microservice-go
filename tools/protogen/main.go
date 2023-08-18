package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {
	//baseDirByte, _ := exec.Command("pwd").Output()
	//baseDir := strings.ReplaceAll(string(baseDirByte), "\n", "")
	listPath := []string{"proto/models", "proto/services"}
	//listPath := []string{"proto/models"}

	for _, path := range listPath {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			// If is
			currentPath := fmt.Sprintf("%s/%s", path, file.Name())
			log.Printf("Generating file %s", currentPath)
			out, err := exec.Command(
				"protoc",
				"--go-grpc_out=exmsg",
				"--proto_path=proto",
				"--go_out=exmsg",
				"--go_opt=paths=source_relative",
				"--go-grpc_opt=paths=source_relative",
				currentPath).CombinedOutput()

			if err != nil {
				log.Printf("Errors %s", out)
			}
			log.Printf("Generating done %s", currentPath)
		}
	}
}
