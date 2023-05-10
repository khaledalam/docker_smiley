package entity

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"os"
	"strings"
)

type LevelType int

// (highest (level1) to lowest (level7))
const (
	Level1 LevelType = iota
	Level2
	Level3
	Level4
	Level5
	Level6
	Level7
)

var LevelToString = map[LevelType]string{
	Level1: "Level1",
	Level2: "Level2",
	Level3: "Level3",
	Level4: "Level4",
	Level5: "Level5",
	Level6: "Level6",
	Level7: "Level7",
}

func (lT LevelType) String() string {
	return LevelToString[lT]
}

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
}

type Env struct {
	types.PluginEnv
	Levels []Level `json:"Levels"`
}

type Level struct {
	LevelType   LevelType `json:"LevelType"`
	LevelString string    `json:"LevelString"`
	Descriptor  string    `json:"Descriptor"`
	IsSet       bool      `json:"IsSet"`
	Value       string    `json:"Value"`
}

type EnvFilters struct {
	Keyword   string
	LevelType LevelType
}

func (dCL DockerContainersList) GroupByApp() map[string][]DockerContainer {

	var ret = map[string][]DockerContainer{}

	for _, d := range dCL {
		ret[d.LabelApp] = append(ret[d.LabelApp], d)
	}

	return ret
}

type DockerContainersList []DockerContainer
type EnvsList []Env
type LevelsList []Level

func (env Env) FilterByKeyword(keyword string) EnvsList {
	var envs []Env
	for _, l := range envs {
		if strings.Contains(l.Name, keyword) || strings.Contains(*l.Value, keyword) {
			envs = append(envs, l)
		}
	}
	return envs
}

func (env Env) FilterByLevel(level LevelType, isSet bool) EnvsList {
	var envs []Env
	for _, e := range envs {
		for _, l := range e.Levels {
			if l.LevelType == level && l.IsSet == isSet {
				envs = append(envs, e)
				break
			}
		}
	}
	return envs
}
