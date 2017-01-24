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
	Host         string            // Changelog server host
	Port         string            // Changelog server port
	Endpoint     string            // Endpoint defaults to /api/events
	Category     string            // Category defaults to misc
	Severity     string            // Severity defaults to "INFO"
	AuthUser     string            // BasicAuth user
	AuthPassword string            // BasicAuth password
	ExtraHeaders map[string]string // Extra headers for auth for example
	ExtraFields  map[string]string // Extra fields
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
		Host:         host,
		Port:         port,
		Endpoint:     endpoint,
		Category:     category,
		Severity:     severity,
		ExtraHeaders: map[string]string{},
		ExtraFields:  map[string]string{},
	}
}

func (c *Changelog) buildUrl() string {
	var url string

	if c.Port != "" {
		url = fmt.Sprintf("%s:%s%s", c.Host, c.Port, c.Endpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.Host, c.Endpoint)
	}

	return url
}

func (c *Changelog) AddExtraHeaders(extra_headers map[string]string) {
	for k, v := range extra_headers {
		c.ExtraHeaders[k] = v
	}
}

func (c *Changelog) AddExtraFields(extra_fields map[string]string) {
	for k, v := range extra_fields {
		c.ExtraFields[k] = v
	}
}

func (c *Changelog) buildMessage(message string) (fields map[string]string) {
	now := time.Now()
	secs := fmt.Sprintf("%d", now.Unix())

	fields = map[string]string{"criticality": fmt.Sprintf("%d", deferSeverity(c.Severity)), "unix_timestamp": secs, "category": c.Category, "description": message}

	for k, v := range c.ExtraFields {
		fields[k] = v
	}

	return fields

}

func (c *Changelog) UseBasicAuth(username string, password string) {
	c.AuthUser = username
	c.AuthPassword = password
}

func (c *Changelog) Send(message string) (response string, err error) {

	fields := c.buildMessage(message)
	jsonMessage := new(bytes.Buffer)
	err = json.NewEncoder(jsonMessage).Encode(fields)
	if err != nil {
		return "NOK", err
	}

	req, err := http.NewRequest("POST", c.buildUrl(), jsonMessage)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-changelog/client")

	// basic auth if username and password is set
	if c.AuthUser != "" && c.AuthPassword != "" {
		req.SetBasicAuth(c.AuthUser, c.AuthPassword)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), err

}
