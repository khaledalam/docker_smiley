package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-contrib/cors"
	"github.com/gorilla/websocket"
	"log"
	"main/entities"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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

		var containers entities.ContainersList

		for _, container_ := range containers_ {

			reader, err := cli.ContainerInspect(context.Background(), container_.ID)
			if err != nil {
				return
			}

			container := entities.Container{
				Container:     container_,
				Envs:          entities.GetEnvsOfContainer(cli, container_),
				LabelApp:      container_.Labels["com.docker.compose.project"],
				ContainerJson: reader,
				Processes:     []os.Process{},
				Logs:          entities.GetContainerLogs(cli, container_),
			}

			containers = append(containers, container)
		}

		jsonRes, err := json.Marshal(containers.GroupByApp())

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		fmt.Fprint(c.Writer, string(jsonRes))
	})

	r.Run()

	// lsof -n -i :8010

}
