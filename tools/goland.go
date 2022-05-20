package cmd

//
//import (
//	"encoding/json"
//	"encoding/xml"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"os"
//	"path/filepath"
//)
//
//package cmd
//
//import (
//"gitlab.shopbase.dev/brodev/beetoolkit/devez/configuration"
//"gitlab.shopbase.dev/brodev/beetoolkit/devez/log"
//"encoding/json"
//"encoding/xml"
//"fmt"
//"github.com/spf13/cobra"
//"io/ioutil"
//"os"
//"path/filepath"
//)
//
//func main() {
//	currentDir, _ := os.Getwd()
//
//	pathProject := fmt.Sprintf("%v/.idea", currentDir)
//
//	if _, err := os.Stat(pathProject); os.IsNotExist(err) {
//		log.Fatal("Directory a has never been opened by goland")
//		return
//	}
//
//	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
//		if info.Name() == "cicd.json" {
//			RunCommand(currentDir, filepath.Dir(path))
//		}
//		return nil
//	})
//
//	if err != nil {
//		log.Fatal(fmt.Sprintf("Error %v", err))
//	}
//}
//
//func RunCommand(currentDir, dir string) {
//	println(dir)
//
//	file, err := ioutil.ReadFile(dir + "/cicd.json")
//
//	cicd := configuration.Cicd{}
//
//	err = json.Unmarshal(file, &cicd)
//	path := currentDir + "/.idea/runConfigurations"
//
//	if _, err := os.Stat(path); os.IsNotExist(err) {
//		os.Mkdir(path, os.ModePerm)
//	}
//
//	nameFile := fmt.Sprintf("%v/%v.xml", path, cicd.Deploy.ServiceName)
//
//	writer, err := os.OpenFile(nameFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
//
//	output, err := xml.MarshalIndent(configuration.New(cicd, filepath.Base(currentDir)), "  ", "    ")
//	if err != nil {
//		println(fmt.Sprintf("Error 1 %v", err))
//	}
//
//	_, err = writer.Write(output)
//
//	if err != nil {
//		println(fmt.Sprintf("Error 2 %v", err))
//	}
//
//}
