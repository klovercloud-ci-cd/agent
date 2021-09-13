package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var ServerPort string
var EventStoreUrl string
var IsK8 string
func InitEnvironmentVariables(){
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	EventStoreUrl=os.Getenv("EVENT_STORE_URL")
	ServerPort = os.Getenv("SERVER_PORT")

}