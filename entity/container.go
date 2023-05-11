package entity

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"os"
)

func (dC DockerContainer) ToJson() []byte {
	jsonRet, _ := json.Marshal(dC)
	return jsonRet
}

type DockerContainer struct {
	types.Container
	Envs          map[string]Env `json:"Envs"`
	Processes     []os.Process
	LabelApp      string `json:"LabelApp"`
	ContainerJson types.ContainerJSON
	Logs          []Log `json:"Logs"`
}

func (dCL DockerContainersList) GroupByApp() map[string][]DockerContainer {

	var ret = map[string][]DockerContainer{}

	for _, d := range dCL {
		ret[d.LabelApp] = append(ret[d.LabelApp], d)
	}

	return ret
}

type DockerContainersList []DockerContainer
