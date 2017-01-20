package settings

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"os"
	"time"
)

var environments = map[string]string{
	"production":    "settings/prod.json",
	"preproduction": "settings/pre.json",
	"tests":         "settings/tests.json",
}

type Settings struct {
	DebugMode          bool
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
	MongoAddrs         []string
	MongoTimeout       int
	MongoDatabase      string
	MongoUsername      string
	MongoPassword      string
	SmtpHost           string
	SmtpPort           int
	SmtpUser           string
	SmtpPassword       string
	DefaultLocale      string
	NoReplyFrom        string
	BaseFrontUrl       string
}

var settings Settings = Settings{}
var env = "preproduction"

func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

func GetEnvironment() string {
	return env
}

func Get() Settings {
	if &settings == nil {
		Init()
	}
	return settings
}

func IsTestEnvironment() bool {
	return env == "tests"
}

func GetMongoSession() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    Get().MongoAddrs,
		Timeout:  60 * time.Second,
		Database: Get().MongoDatabase,
		Username: Get().MongoUsername,
		Password: Get().MongoPassword,
	}
	// Connect to our local mongo
	mgoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	mgoSession.SetMode(mgo.Monotonic, true)

	return mgoSession
}

func GetSmtpDial() *gomail.Dialer {
	return gomail.NewDialer(Get().SmtpHost, Get().SmtpPort, Get().SmtpUser, Get().SmtpPassword)
}
