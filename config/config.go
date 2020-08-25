package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configurations
type Config struct {
	Database  Database
	General   General
	Messaging Messaging
}

// Database - holds credentials for Database
type Database struct {
	dbUser     string
	dbPswd     string
	dbHost     string
	dbPort     string
	dbName     string
	testDBHost string
	testDBName string
}

// General - holds general data for confoguring
type General struct {
	HTTPRequestTimeout int
	Timezone           string
	TimeFullFormat     string
}

// Messaging -  holds credentials for Slack
type Messaging struct {
	SlackBaseURL   string
	SlackAuthToken string
	SlackChannelID string
	SlackLinkNames bool
}

// Get configurations
func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.Database.dbUser, "dbuser", os.Getenv("DB_USER"), "DB user name")
	flag.StringVar(&conf.Database.dbPswd, "dbpswd", os.Getenv("DB_PASSWORD"), "DB pass")
	flag.StringVar(&conf.Database.dbPort, "dbport", os.Getenv("DB_PORT"), "DB port")
	flag.StringVar(&conf.Database.dbHost, "dbhost", os.Getenv("DB_HOST"), "DB host")
	flag.StringVar(&conf.Database.dbName, "dbname", os.Getenv("DB_NAME"), "DB name")
	flag.StringVar(&conf.Database.testDBHost, "testdbhost", os.Getenv("TEST_DB_HOST"), "test database host")
	flag.StringVar(&conf.Database.testDBName, "testdbname", os.Getenv("TEST_DB_NAME"), "test database name")

	rqTimeout, _ := strconv.Atoi(os.Getenv("HTTTP_REQUEST_TIMEOUT"))

	flag.IntVar(&conf.General.HTTPRequestTimeout, "httprequesttimeour", rqTimeout, "http request timeout")
	flag.StringVar(&conf.General.Timezone, "timezone", os.Getenv("TIMEZONE"), "Application Timezone")
	flag.StringVar(&conf.General.TimeFullFormat, "timefullformat", os.Getenv("TIMEFULLFORMAT"), "Format Time")

	slackLinkNames, _ := strconv.ParseBool(os.Getenv("SLACK_LINK_NAMES"))

	flag.StringVar(&conf.Messaging.SlackBaseURL, "slackbaseurl", os.Getenv("SLACK_BASE_URL"), "slack base url")
	flag.StringVar(&conf.Messaging.SlackAuthToken, "slackauthtoken", os.Getenv("SLACK_AUTH_TOKEN"), "slack auth token")
	flag.StringVar(&conf.Messaging.SlackChannelID, "slackchannelid", os.Getenv("SLACK_CHANNEL_ID"), "slack channel id")
	flag.BoolVar(&conf.Messaging.SlackLinkNames, "slacklinknames", slackLinkNames, "slack link names")

	flag.Parse()

	return conf
}

// GetDBConnStr handler
func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.Database.dbHost, c.Database.dbName)
}

// GetTestDBConnStr handler
func (c *Config) GetTestDBConnStr() string {
	return c.getDBConnStr(c.Database.testDBHost, c.Database.testDBName)
}

// getDBConnStr
func (c *Config) getDBConnStr(dbhost, dbname string) (dbURI string) {

	// If the optional DB_TCP_HOST environment variable is set, it contains
	// the IP address and port number of a TCP connection pool to be created,
	// such as "127.0.0.1:3306". If DB_TCP_HOST is not set, a Unix socket
	// connection pool will be created instead.
	if dbhost == "127.0.0.1" {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", c.Database.dbUser, c.Database.dbPswd, c.Database.dbHost, c.Database.dbName)
	} else {
		dbURI = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true", c.Database.dbUser, c.Database.dbPswd, c.Database.dbHost, c.Database.dbName)
	}

	return
}
