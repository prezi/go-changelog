package changelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	INFO         = 1
	NOTIFICATION = 2
	WARNING      = 3
	ERROR        = 4
	CRITICAL     = 5
)

// Changelog struct to hold the provided options
type Changelog struct {
	host         string            // changelog server host
	port         string            // changelog server port
	endpoint     string            // endpoint defaults to /api/events
	category     string            // category defaults to misc
	severity     string            // severity defaults to "INFO"
	authUser     string            // basicAuth user
	authPassword string            // basicAuth password
	extraHeaders map[string]string // extra headers for auth for example
	extraFields  map[string]string // extra fields
}

var severityLookup = map[string]int{
	"INFO":         INFO,
	"NOTIFICATION": NOTIFICATION,
	"WARNING":      WARNING,
	"ERROR":        ERROR,
	"CRITICAL":     CRITICAL,
}

func deferSeverity(severity string) int {
	return severityLookup[severity]
}

func New(host string, port string, endpoint string, category string, severity string) *Changelog {
	if host == "" {
		host = "http://localhost"
	}

	if endpoint == "" {
		endpoint = "/api/events"
	}

	if category == "" {
		category = "misc"
	}

	if severity == "" {
		severity = "INFO"
	}

	return &Changelog{
		host:         host,
		port:         port,
		endpoint:     endpoint,
		category:     category,
		severity:     severity,
		extraHeaders: map[string]string{},
		extraFields:  map[string]string{},
	}
}

func (c *Changelog) buildUrl() string {
	var url string

	if c.port != "" {
		url = fmt.Sprintf("%s:%s%s", c.host, c.port, c.endpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.host, c.endpoint)
	}

	return url
}

func (c *Changelog) AddExtraHeaders(extra_headers map[string]string) {
	for k, v := range extra_headers {
		c.extraHeaders[k] = v
	}
}

func (c *Changelog) AddExtraFields(extra_fields map[string]string) {
	for k, v := range extra_fields {
		c.extraFields[k] = v
	}
}

func (c *Changelog) buildMessage(message string) (fields map[string]string) {
	now := time.Now()
	secs := fmt.Sprintf("%d", now.Unix())

	fields = map[string]string{"criticality": fmt.Sprintf("%d", deferSeverity(c.severity)), "unix_timestamp": secs, "category": c.category, "description": message}

	for k, v := range c.extraFields {
		fields[k] = v
	}

	return fields

}

func (c *Changelog) UseBasicAuth(username string, password string) {
	c.authUser = username
	c.authPassword = password
}

func (c *Changelog) Send(message string) (response string, err error) {

	fields := c.buildMessage(message)
	jsonMessage := new(bytes.Buffer)
	err = json.NewEncoder(jsonMessage).Encode(fields)
	if err != nil {
		return "NOK", err
	}

	req, err := http.NewRequest("POST", c.buildUrl(), jsonMessage)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-changelog/client")

	// basic auth if username and password is set
	if c.authUser != "" && c.authPassword != "" {
		req.SetBasicAuth(c.authUser, c.authPassword)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), err

}
