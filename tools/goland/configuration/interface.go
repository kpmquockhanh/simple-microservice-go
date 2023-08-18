package configuration

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type moduleConfig struct {
	Name string `xml:"name,attr"`
}

type parameter struct {
	Value string `xml:"value,attr"`
}

type method struct {
	Value string `xml:"v,attr"`
}

type Configuration struct {
	Default          bool         `xml:"default,attr"`
	Name             string       `xml:"name,attr"`
	Type             string       `xml:"type,attr"`
	FactoryName      string       `xml:"factoryName,attr"`
	Module           moduleConfig `xml:"module"`
	WorkingDirectory parameter    `xml:"working_directory"`
	GoParameters     parameter    `xml:"go_parameters"`
	Parameters       parameter    `xml:"parameters"`
	Kind             parameter    `xml:"kind"`
	FilePath         parameter    `xml:"filePath"`
	Package          parameter    `xml:"package"`
	Directory        parameter    `xml:"directory"`
	Method           method       `xml:"method"`
}

type Component struct {
	XMLName       xml.Name      `xml:"component"`
	Name          string        `xml:"name,attr"`
	Configuration Configuration `xml:"configuration"`
}

type IdeConfig struct {
	Component Component `xml:"component"`
}

func New(cicd Cicd, namePackage string) Component {
	AdditionArgs := []string{""}
	if cicd.AppType == "api" {
		AdditionArgs = append(AdditionArgs, fmt.Sprintf("-http-port=%v", cicd.Deploy.DevExposePort), "-api-timeout=10000000")
		cicd.Deploy.ConfigRemoteKeys = append(cicd.Deploy.ConfigRemoteKeys, "/api/common.toml")
	}

	if cicd.AppType == "service" {
		AdditionArgs = append(AdditionArgs, fmt.Sprintf("-grpc-port=%v", cicd.Deploy.DevExposePort), "-dms-timeout=10000000")
		cicd.Deploy.ConfigRemoteKeys = append(cicd.Deploy.ConfigRemoteKeys, "/service/common.toml")
	}

	AdditionArgs = append(AdditionArgs, "-config-remote-address=127.0.0.1:8500 -logger-type=default -logger-dev=true -allow-loopback=true")
	if len(cicd.Deploy.ConfigRemoteKeys) > 0 {
		fullConfigsArgs := fmt.Sprintf("-config-remote-keys=%v", strings.Join(cicd.Deploy.ConfigRemoteKeys, ","))
		AdditionArgs = append(AdditionArgs, fullConfigsArgs)
		AdditionArgs = append(AdditionArgs, "-config-type=remote")
	}

	return Component{
		Name: "ProjectRunConfigurationManager",
		Configuration: Configuration{
			Default:     false,
			Name:        cicd.Deploy.ServiceName,
			Type:        "GoApplicationRunConfiguration",
			FactoryName: "Go Application",
			Module: moduleConfig{
				Name: namePackage,
			},
			WorkingDirectory: parameter{
				Value: "$PROJECT_DIR$/",
			},
			GoParameters: parameter{
				Value: "",
			},
			Parameters: parameter{
				Value: strings.Join(AdditionArgs, " "),
			},
			Kind: parameter{
				Value: "FILE",
			},
			FilePath: parameter{
				Value: fmt.Sprintf("$PROJECT_DIR$/%v/main.go", cicd.Build.CmdBinDir),
			},
			Package: parameter{
				Value: namePackage,
			},
			Directory: parameter{
				Value: "$PROJECT_DIR$/",
			},
			Method: method{
				Value: "2",
			},
		},
	}
}
