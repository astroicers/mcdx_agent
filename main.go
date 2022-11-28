package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Agent struct {
	Agent_id string
	Team_id  string
	Subnet   []string
}

type Health_check struct {
	Health_check string
}

func httpGet() {
	resp, err := http.Get("http://127.0.0.1:60000/agent_info/id")
	if err != nil {
		fmt.Printf("Connect Error: %s", err.Error())
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Data Error: %s", err.Error())
		}
		fmt.Printf("Data: %s\n", string(body))
	}
}

func online_check(ip string) {
	resp, err := http.Get("http://" + ip + ":60000/health_check")
	if err != nil {
		fmt.Printf("Connect Error: %s", err.Error())
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Data Error: %s", err.Error())
		}
		var hc Health_check
		err = json.Unmarshal(body, &hc)
		if err != nil {
			fmt.Printf("Health Check Error: %s", err)
		}
		fmt.Printf("Data: %s\n", hc.Health_check)

	}
}

func server() {
	server := gin.Default()

	net_debug := server.Group("net_debug")
	net_debug.GET(":agent_id", func(c *gin.Context) {
		name := c.Param("agent_id")
		c.JSON(http.StatusOK, gin.H{"agent id": name})
	})

	agent_id_generator := server.Group("agent_id_generator")
	agent_id_generator.GET("/", func(c *gin.Context) {
		time_now := time.Now().Unix()
		hash_id := sha1.New()
		hash_id.Write([]byte(strconv.FormatInt(time_now, 10)))
		c.JSON(http.StatusOK, gin.H{"id": hash_id.Sum(nil)})
	})

	health_check := server.Group("health_check")
	health_check.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health_check": "OK"})
	})
	server.Run(":60000")
}

func main() {
	go server()
	for true {
		online_check("127.0.0.1")
		time.Sleep(1 * time.Second)
	}
}
