package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var ServerPort string
var IsK8 string
func InitEnvironmentVariables(){
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	ServerPort = os.Getenv("SERVER_PORT")

}