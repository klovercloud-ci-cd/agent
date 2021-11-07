package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var ServerPort string
var EventStoreUrl string
var PullSize string
var IsK8 string
var AgentName string
var Publickey string
var EnableAuthentication bool
var Token string
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
	PullSize=os.Getenv("PULL_SIZE")
	Publickey=os.Getenv("PUBLIC_KEY")

	if os.Getenv("ENABLE_AUTHENTICATION")==""{
		EnableAuthentication=false
	}else{
		if strings.ToLower(os.Getenv("ENABLE_AUTHENTICATION"))=="true"{
			EnableAuthentication=true
		}else{
			EnableAuthentication=false
		}
	}
	Token=os.Getenv("TOKEN")
}