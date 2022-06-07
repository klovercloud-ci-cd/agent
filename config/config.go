package config

import (
	"github.com/joho/godotenv"
	"github.com/klovercloud-ci-cd/agent/enums"
	"log"
	"os"
	"strconv"
	"strings"
)

// RunMode refers to run mode.
var RunMode string

// ServerPort refers to server port.
var ServerPort string

// ApiServiceUrl refers to api service url.
var ApiServiceUrl string

// PullSize refers to number of job to be pulled each time.
var PullSize int64

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

// KafkaPublisherEnabled set true if kafka publisher is enabled
var KafkaPublisherEnabled bool

// LighthouseEnabled set true if lighthouse is enabled
var LighthouseEnabled bool

// CurrentConcurrentJobs running jobs count.
var CurrentConcurrentJobs int64

// TerminalBaseUrl base url of terminal.
var TerminalBaseUrl string

// TerminalApiVersion api version of terminal.
var TerminalApiVersion string

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
	ApiServiceUrl = os.Getenv("API_SERVICE_URL")
	if strings.HasPrefix(ApiServiceUrl, "/") {
		ApiServiceUrl = strings.TrimPrefix(ApiServiceUrl, "/")
	}
	ServerPort = os.Getenv("SERVER_PORT")
	AgentName = os.Getenv("AGENT_NAME")
	TerminalBaseUrl = os.Getenv("TERMINAL_BASE_URL")
	TerminalApiVersion=os.Getenv("TERMINAL_API_VERSION")
	err := error(nil)
	PullSize, err = strconv.ParseInt(os.Getenv("PULL_SIZE"), 10, 64)
	if err != nil {
		PullSize = 4
	}
	log.Println(PullSize)
	CurrentConcurrentJobs=0
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

	if os.Getenv("ENABLE_OPENTRACING") == "" {
		EnableOpenTracing = false
	} else {
		if strings.ToLower(os.Getenv("ENABLE_OPENTRACING")) == "true" {
			EnableOpenTracing = true
		} else {
			EnableOpenTracing = false
		}
	}

	if os.Getenv("KAFKA_PUBLISHER_ENABLED") == "true" {
		KafkaPublisherEnabled = true
	} else {
		KafkaPublisherEnabled = false
	}

	if os.Getenv("LIGHTHOUSE_ENABLED") == "true" {
		LighthouseEnabled = true
	} else {
		LighthouseEnabled = false
	}
}
