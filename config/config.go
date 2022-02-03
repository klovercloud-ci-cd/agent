package config

import (
	"github.com/joho/godotenv"
	"github.com/klovercloud-ci-cd/agent/enums"
	"log"
	"os"
	"strings"
)

// RunMode refers to run mode.
var RunMode string

// ServerPort refers to server port.
var ServerPort string

// EventStoreUrl refers to event-bank url.
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

// Token refers to oauth token for service to service communication.
var Token string


// EnableOpenTracing set true if opentracing is needed.
var EnableOpenTracing bool

// InitEnvironmentVariables initializes environment variables
func InitEnvironmentVariables() {
	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		RunMode = string(enums.DEVELOP)
	}

	if RunMode != string(enums.PRODUCTION) {
		//Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("ERROR:", err.Error())
			return
		}
	}
	log.Println("RUN MODE:", RunMode)
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


	if os.Getenv("ENABLE_OPENTRACING")==""{
		EnableOpenTracing=false
	}else{
		if strings.ToLower(os.Getenv("ENABLE_OPENTRACING")) == "true" {
			EnableOpenTracing = true
		} else {
			EnableOpenTracing = false
		}
	}
}
