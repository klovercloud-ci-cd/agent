package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

// ServerPort refers to server port.
var ServerPort string

// EventStoreUrl refers to event-store url.
var EventStoreUrl string

// PullSize refers to number of job to be pulled each time.
var PullSize string

// IsK8 refers if application is running inside k8s.
var IsK8 string

// AgentName refers to name of the agent.
var AgentName string

// Publickey refers to publickey of EventStoreToken.
var Publickey string

// EnableAuthentication refers if service to service authentication is enabled.
var EnableAuthentication bool

// Token refers to jwt token for service to service communication.
var Token string

// InitEnvironmentVariables initializes environment variables
func InitEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	EventStoreUrl = os.Getenv("EVENT_STORE_URL")
	if strings.HasPrefix(EventStoreUrl, "/") {
		EventStoreUrl = strings.TrimPrefix(EventStoreUrl, "/")
	}
	ServerPort = os.Getenv("SERVER_PORT")
	AgentName = os.Getenv("AGENT_NAME")
	PullSize = os.Getenv("PULL_SIZE")
	Publickey = os.Getenv("PUBLIC_KEY")
	IsK8 = os.Getenv("IS_K8")
	if os.Getenv("ENABLE_AUTHENTICATION") == "" {
		EnableAuthentication = false
	} else {
		if strings.ToLower(os.Getenv("ENABLE_AUTHENTICATION")) == "true" {
			EnableAuthentication = true
		} else {
			EnableAuthentication = false
		}
	}
	Token = os.Getenv("TOKEN")
}
