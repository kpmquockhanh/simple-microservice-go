package configuration

type CicdBuild struct {
	CmdBinDir string `json:"cmd_bin_dir"`
	ImageName string `json:"image_name"`
}

type CicdDeploy struct {
	ServiceName      string   `json:"service_name"`
	DevExposePort    int64    `json:"dev_expose_port"`
	ConfigRemoteKeys []string `json:"config_remote_keys"`
}

type Cicd struct {
	AllowGenerateJenkinsfile bool       `json:"allow_generate_jenkinsfile"`
	AppType                  string     `json:"app_type"`
	Build                    CicdBuild  `json:"build"`
	Deploy                   CicdDeploy `json:"deploy"`
}
