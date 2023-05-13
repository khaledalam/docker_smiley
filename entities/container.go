package entities

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (dC Container) ToJson() []byte {
	jsonRet, _ := json.Marshal(dC)
	return jsonRet
}

type Container struct {
	types.Container
	Envs          map[string]Env `json:"Envs"`
	Processes     []os.Process
	LabelApp      string `json:"LabelApp"`
	ContainerJson types.ContainerJSON
	Logs          []Log `json:"Logs"`
}

func (dCL ContainersList) GroupByApp() map[string][]Container {

	var ret = map[string][]Container{}

	for _, d := range dCL {
		ret[d.LabelApp] = append(ret[d.LabelApp], d)
	}

	return ret
}

type ContainersList []Container

func GetEnvsOfContainer(cli *client.Client, container types.Container) map[string]Env {
	reader, err := cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		return nil
	}

	envs := map[string]Env{}
	for _, env := range reader.Config.Env {

		var key, val = getEnvStringParts(env)

		pair := types.PluginEnv{
			Name:  key,
			Value: &val,
		}

		envs[key] = Env{
			PluginEnv: pair,
			Levels:    getEnvLevels(container, pair),
		}
	}
	return envs
}

func getEnvStringParts(env string) (string, string) {
	parts := strings.Split(env, "=")
	key := parts[0]
	val := strings.Join(parts[1:], "=")
	return key, val
}

func GetContainerLogs(cli *client.Client, container types.Container) []Log {
	logsReader, err := cli.ContainerLogs(context.Background(), container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: false,
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	defer func(logsReader io.ReadCloser) {
		_ = logsReader.Close()
	}(logsReader)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, logsReader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
		return nil
	}

	var logs []Log

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		logs = append(logs, Log{Line: scanner.Text()})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return logs
}

func getEnvLevels(container types.Container, pair types.PluginEnv) LevelsList {

	nilValue := "nil_value"

	level1Value := nilValue
	level2Value := nilValue
	level3Value := nilValue
	level4Value := nilValue
	level5Value := nilValue
	level6Value := nilValue
	level7Value := nilValue

	/*
		Reverse engineering [Level 1] :
		(Set using docker compose run -e in the CLI)
		''''''''''''''''''''''''''''''''''''''''''''''

		1. Filter all running machine processes that contains container names or ID in its command using "ps" and "grep" commands.
		2. Filter result of step(1) that contain substring "PluginEnv.Name=" in its command.
	*/
	// remove prefix "/"
	grepArg := lo.Map(container.Names, func(name string, index int) string {
		return strings.TrimPrefix(name, "/")
	})
	grepArg = append(grepArg, container.ID)
	grepArg = append(grepArg, container.Labels["com.docker.compose.project"])
	for _, grepArg := range grepArg {
		grepCmd := exec.Command("/bin/sh", "-c", "ps -a | grep "+grepArg)
		out, err := grepCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		commandsContainsContainerIdOrNameOutput := strings.Split(string(out), "\n")
		for _, command := range commandsContainsContainerIdOrNameOutput {

			// exclude grepCmd from processes list
			// @TODO: remove "grep" command too.
			if !strings.Contains(command, strconv.Itoa(grepCmd.Process.Pid)) {
				if strings.Contains(command, pair.Name+"=") {

					split := strings.Split(command, pair.Name+"=")
					if len(split) >= 1 {
						// @TODO: handle values that containers ("\"", "'", " ")
						envVal := strings.Split(split[1], " ")
						if len(envVal) >= 1 {
							level1Value = envVal[0]
							break
						}
					}
				}
			}
		}
	}

	/*
		Reverse engineering [Level 2] :
		(Substituted from your shell)
		'''''''''''''''''''''''''''''''

		1. Parse docker-compose.yml file if existed.
		2. Check if there is any line contains ${PluginEnv.Name}.
		3. Get ${PluginEnv.Name} value from container shell. (empty string if var not set)
	*/
	file, err := os.Open(container.Labels["com.docker.compose.project.config_files"])
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "${"+pair.Name+"}") {
			// get value from container shell

			//fmt.Println("Yes > ", pair.Name, scanner.Text())
			//
			//cmd := "docker exec -it " + container.ID + " echo ${" + pair.Name + "}"
			//fmt.Println(cmd)
			//echoCmd := exec.Command("/bin/sh", "-c", cmd)
			//out, err := echoCmd.Output()
			//fmt.Println(">>", string(out))
			//if err != nil {
			//	fmt.Println("err >>", err)
			//
			//	log.Fatal(err)
			//} else {
			//	fmt.Println(">>", string(out))
			//}

		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return LevelsList{
		Level{
			LevelType:   Level1,
			LevelString: Level1.String(),
			Descriptor:  "Set using docker compose run -e in the CLI",
			IsSet:       level1Value != nilValue,
			Value:       level1Value,
		},
		Level{
			LevelType:   Level2,
			LevelString: Level2.String(),
			Descriptor:  "Substituted from your shell",
			IsSet:       level2Value != nilValue,
			Value:       level2Value,
		},
		Level{
			LevelType:   Level3,
			LevelString: Level3.String(),
			Descriptor:  "Set using the environment attribute in the Compose file",
			IsSet:       level3Value != nilValue,
			Value:       level3Value,
		},
		Level{
			LevelType:   Level4,
			LevelString: Level4.String(),
			Descriptor:  "Use of the --env-file argument in the CLI",
			IsSet:       level4Value != nilValue,
			Value:       level4Value,
		},
		Level{
			LevelType:   Level5,
			LevelString: Level5.String(),
			Descriptor:  "Use of the env_file attribute in the Compose file",
			IsSet:       level5Value != nilValue,
			Value:       level5Value,
		},
		Level{
			LevelType:   Level6,
			LevelString: Level6.String(),
			Descriptor:  "Set using an .env file placed at base of your project directory",
			IsSet:       level6Value != nilValue,
			Value:       level6Value,
		},
		Level{
			LevelType:   Level7,
			LevelString: Level7.String(),
			Descriptor:  "Set in a container image in the ENV directive. Having any ARG or ENV setting in a Dockerfile evaluates only if there is no Docker Compose entry for environment, env_file or run --env.",
			IsSet:       level7Value != nilValue,
			Value:       level7Value,
		},
	}
}
