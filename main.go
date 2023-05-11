/**
 * @author  Khaled Alam
 */

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-contrib/cors"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"io"
	"log"
	"main/entity"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

func getEnvStringParts(env string) (string, string) {
	parts := strings.Split(env, "=")
	key := parts[0]
	val := strings.Join(parts[1:], "=")
	return key, val
}

func getEnvLevels(container types.Container, pair types.PluginEnv) entity.LevelsList {

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
		return strings.Replace(name, "/", "", 1)
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

	return entity.LevelsList{
		entity.Level{
			LevelType:   entity.Level1,
			LevelString: entity.Level1.String(),
			Descriptor:  "Set using docker compose run -e in the CLI",
			IsSet:       level1Value != nilValue,
			Value:       level1Value,
		},
		entity.Level{
			LevelType:   entity.Level2,
			LevelString: entity.Level2.String(),
			Descriptor:  "Substituted from your shell",
			IsSet:       level2Value != nilValue,
			Value:       level2Value,
		},
		entity.Level{
			LevelType:   entity.Level3,
			LevelString: entity.Level3.String(),
			Descriptor:  "Set using the environment attribute in the Compose file",
			IsSet:       level3Value != nilValue,
			Value:       level3Value,
		},
		entity.Level{
			LevelType:   entity.Level4,
			LevelString: entity.Level4.String(),
			Descriptor:  "Use of the --env-file argument in the CLI",
			IsSet:       level4Value != nilValue,
			Value:       level4Value,
		},
		entity.Level{
			LevelType:   entity.Level5,
			LevelString: entity.Level5.String(),
			Descriptor:  "Use of the env_file attribute in the Compose file",
			IsSet:       level5Value != nilValue,
			Value:       level5Value,
		},
		entity.Level{
			LevelType:   entity.Level6,
			LevelString: entity.Level6.String(),
			Descriptor:  "Set using an .env file placed at base of your project directory",
			IsSet:       level6Value != nilValue,
			Value:       level6Value,
		},
		entity.Level{
			LevelType:   entity.Level7,
			LevelString: entity.Level7.String(),
			Descriptor:  "Set in a container image in the ENV directive. Having any ARG or ENV setting in a Dockerfile evaluates only if there is no Docker Compose entry for environment, env_file or run --env.",
			IsSet:       level7Value != nilValue,
			Value:       level7Value,
		},
	}
}

func getEnvsOfContainer(cli *client.Client, container types.Container) map[string]entity.Env {
	reader, err := cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		return nil
	}

	//fmt.Println(container.Names[0], reader.Args)

	envs := map[string]entity.Env{}
	for _, env := range reader.Config.Env {

		var key, val = getEnvStringParts(env)

		pair := types.PluginEnv{
			Name:  key,
			Value: &val,
		}

		envs[key] = entity.Env{
			PluginEnv: pair,
			Levels:    getEnvLevels(container, pair),
		}
	}
	return envs
}

func getContainerLogs(cli *client.Client, container types.Container) []entity.Log {
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

	var logs []entity.Log

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		logs = append(logs, entity.Log{Line: scanner.Text()})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return logs
}

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:"message"`
}

func main() {

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ws", func(c *gin.Context) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if !errors.Is(err, nil) {
			log.Println(err)
		}
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				log.Printf("[ws.Close()] error occurred: %v", err)
			}
		}(ws)

		log.Println("WS Connected!")

		for {
			var message Message
			err := ws.ReadJSON(&message)
			if !errors.Is(err, nil) {
				log.Printf("[ws.ReadJSON(&message)] error occurred: %v", err)
				break
			}
			log.Println(message)

			// send message from server
			if err := ws.WriteJSON(message); !errors.Is(err, nil) {
				log.Printf("[ws.WriteJSON(message)] error occurred: %v", err)
			}
		}
	})

	r.GET("/env/list", func(c *gin.Context) {

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		containers_, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		var containers entity.DockerContainersList

		for _, container_ := range containers_ {

			reader, err := cli.ContainerInspect(context.Background(), container_.ID)
			if err != nil {
				return
			}

			container := entity.DockerContainer{
				Container:     container_,
				Envs:          getEnvsOfContainer(cli, container_),
				LabelApp:      container_.Labels["com.docker.compose.project"],
				ContainerJson: reader,
				Processes:     []os.Process{},
				Logs:          getContainerLogs(cli, container_),
			}

			containers = append(containers, container)
		}

		jsonRes, err := json.Marshal(containers.GroupByApp())

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		fmt.Fprint(c.Writer, string(jsonRes))
	})

	r.Run()

}
