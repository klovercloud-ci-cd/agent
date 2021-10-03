package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var ServerPort string
var EventStoreUrl string
var EventStoreToken string
var PullSize string
var IsK8 string
var AgentName string
func InitEnvironmentVariables(){
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	EventStoreUrl=os.Getenv("EVENT_STORE_URL")
	if strings.HasPrefix(EventStoreUrl, "/") {
		EventStoreUrl = strings.TrimPrefix(EventStoreUrl, "/")
	}
	ServerPort = os.Getenv("SERVER_PORT")
	AgentName=os.Getenv("AGENT_NAME")
	EventStoreToken=os.Getenv("EVENT_STORE_URL_TOKEN")
	PullSize=os.Getenv("PULL_SIZE")
}