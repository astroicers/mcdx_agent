package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/netip"
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

// func httpGet() {
// 	resp, err := http.Get("http://127.0.0.1:60000/agent_info/id")
// 	if err != nil {
// 		fmt.Printf("Connect Error: %s", err.Error())
// 	} else {
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Printf("Data Error: %s", err.Error())
// 		}
// 		fmt.Printf("Data: %s\n", string(body))
// 	}
// }

func Hosts(cidr string) ([]netip.Addr, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}

	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}

	if len(ips) < 2 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}

func get_api(ip string, dir string) {
	resp, err := http.Get("http://" + ip + ":60000/" + dir)
	if err != nil {
		// fmt.Printf("Connect Error: %s\n", err.Error())
		log.Fatalf("Connect Error: %v", err.Error())
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Data Error: %s\n", err.Error())
		}
		var hc Health_check
		err = json.Unmarshal(body, &hc)
		if err != nil {
			fmt.Printf("Health Check Error: %s\n", err)
		}
		fmt.Printf("Data: %s\n", hc.Health_check)

	}
}

func online_check(ip string) {

	ips, _ := Hosts("172.16.90.0/24")
	// fmt.Printf("%s", ips)
	for _, ip := range ips {
		// fmt.Printf("%s\n\n", ip)
		go get_api(ip.String(), "health_check")
	}

}

func server() {
	server := gin.Default()

	net_debug := server.Group("net_debug")
	net_debug.GET(":agent_id", func(c *gin.Context) {
		name := c.Param("agent_id")
		c.JSON(http.StatusOK, gin.H{"agent id": name})
	})

	id_generator := server.Group("id_generator")
	id_generator.GET("/", func(c *gin.Context) {
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
		time.Sleep(60 * time.Second)
	}
}
