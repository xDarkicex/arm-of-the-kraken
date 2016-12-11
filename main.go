package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	load()
}

// Might be an okay port..
var port = "4980"

func main() {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	for {
		var object reply

		defer listen.Close()
		connection, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(connection)
		for scanner.Scan() {
			args := strings.Split(strings.ToLower(scanner.Text()), " ")
			switch args[0] {
			case "status":
				fileInfo, err := os.Stat("log/" + data.Env + ".log")
				if err != nil {
					fmt.Printf("Error getting stat %s", err)
				}
				i64, err := strconv.ParseInt(args[1], 10, 32)
				if err != nil {
					msg := fmt.Sprintf("Error converting java last file time too golanf time err: %s", err)
					object = reply{
						Error: true,
						Response: response{
							Message: msg,
						},
					}
				} else {
					time := time.Unix(i64, 0)
					modTime := fileInfo.ModTime()
					modString := fmt.Sprintf("%d", modTime.Unix())
					if modTime.After(time) {
						object = reply{
							Error: false,
							Response: response{
								Message: "There's been a change!",
								Data: []string{
									modString,
								},
							},
						}
					} else {
						object = reply{
							Error: false,
							Response: response{
								Message: "Log is fine",
								Data: []string{
									modString,
								},
							},
						}
					}
				}
				js, err := json.Marshal(object)
				fmt.Fprintln(connection, string(js))
			}
		}
	}
}

////////////////////////////
// Config loader
///////////////////////////
// Data struct for config.Data
var data config

// config is for config.Data
type config struct {
	Version string
	Port    int
	Host    string
	Secret  string
	Verbose bool
	Env     string
	Email   string
	SMTP    smtp
}

// SMTP for smtp settings
type smtp struct {
	Host     string
	Port     int
	Password string
}

//Load loads config file
func load() {
	config, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println(err, config)
	}
	// fmt.Println(config)
	err = json.Unmarshal(config, &data)
	if err != nil {
		fmt.Println(err)
	}
}

///////////////////////////
// Structs for JSON REPLY
//////////////////////////

type reply struct {
	Error    bool
	Response response
}
type response struct {
	Data    []string
	Message string
}
