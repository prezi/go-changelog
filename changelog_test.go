package changelog

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

// TestSeverity verifies that the returned severity is match matching
func TestSeverity(t *testing.T) {
	assertEqual(t, 1, deferSeverity("INFO"), fmt.Sprintf("INFO should be %d", INFO))
	assertEqual(t, 2, deferSeverity("NOTIFICATION"), fmt.Sprintf("NOTIFICATION should be %d", NOTIFICATION))
	assertEqual(t, 3, deferSeverity("WARNING"), fmt.Sprintf("WARNING should be %d", WARNING))
	assertEqual(t, 4, deferSeverity("ERROR"), fmt.Sprintf("ERROR should be %d", ERROR))
	assertEqual(t, 5, deferSeverity("CRITICAL"), fmt.Sprintf("CRITICAL should be %d", CRITICAL))
}

// TestNewWithDefaultValues verifies that the default variables are properly initialized
func TestNewWithDefaultValues(t *testing.T) {
	c := New("", "", "", "", "")
	assertEqual(t, "http://localhost", c.host, "Default host should be http://localhost")
	assertEqual(t, "", c.port, "Default port should be empty")
	assertEqual(t, "/api/events", c.endpoint, "Default endpoint should be /api/events")
	assertEqual(t, "misc", c.category, "Default category should be misc")
	assertEqual(t, "INFO", c.severity, "Default severity should be INFO")
}

// TestNewWithCustomValues verifies that the custom variables are properly initialized
func TestNewWithCustomValues(t *testing.T) {
	c := New("https://serverurl", "8080", "/customurl/api/events", "production", "WARNING")
	assertEqual(t, "https://serverurl", c.host, "Host should be https://serverurl")
	assertEqual(t, "8080", c.port, "Port should be 8080")
	assertEqual(t, "/customurl/api/events", c.endpoint, "Endpoint should be /customurl/api/events")
	assertEqual(t, "production", c.category, "Category should be production")
	assertEqual(t, "WARNING", c.severity, "Severity should be WARNING")
}

// TestBuildUrl will verify that the buildUrl generates the proper URL
func TestBuildUrl(t *testing.T) {
	c := New("https://server.tld", "9000", "", "", "")
	assertEqual(t, c.buildUrl(), "https://server.tld:9000/api/events", "The URL should be https://server.tld:9000/api/events")
}

// TestAddExtraHeaders will verify thet the function can be safely called
func TestAddExtraHeaders(t *testing.T) {
	c := New("", "", "", "", "")
	extraHeaders := map[string]string{"username": "foo", "password": "bar"}
	c.AddExtraHeaders(extraHeaders)
	assertEqual(t, "foo", c.extraHeaders["username"], "Username should be foo")
	assertEqual(t, "bar", c.extraHeaders["password"], "Password should be bar")
}

// TestAddExtraHeaders will verify thet the function can be safely called
func TestAddExtraFields(t *testing.T) {
	c := New("", "", "", "", "")
	extraFields := map[string]string{"environment": "production"}
	c.AddExtraFields(extraFields)
	assertEqual(t, "production", c.extraFields["environment"], "Environment should be production")
}

// TestSend will verify thata Send can be safely called
func TestSend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost/api/events",
		httpmock.NewStringResponder(200, `OK`))
	c := New("", "", "", "misc", "")
	response, err := c.Send("Message")
	assertEqual(t, "OK", response, "Response should be OK")
	assertEqual(t, nil, err, "err should be nil")
}

// TestBuildMessageWithNoCustomFields will verify that everything works as expected
func TestBuildMessageWithNoCustomFields(t *testing.T) {
	c := New("", "", "", "", "")
	fields := c.buildMessage("Hello")
	val, ok := fields["criticality"]
	assertEqual(t, true, ok, "criticality should be present")
	assertEqual(t, val, fmt.Sprintf("%d", deferSeverity("INFO")), "val should be 1")
	val, ok = fields["category"]
	assertEqual(t, true, ok, "category should be present")
	assertEqual(t, val, "misc", "val should be misc")
	_, ok = fields["unix_timestamp"]
	assertEqual(t, true, ok, "unix_timestamp should be present")
	val, ok = fields["description"]
	assertEqual(t, true, ok, "description should be present")
	assertEqual(t, val, "Hello", "val should be Hello")
	_, ok = fields["environment"]
	assertEqual(t, false, ok, "custom field environment shouldn't be present")
}

// TestBuildMessageWithNoCustomFields will verify that everything works as expected
func TestBuildMessageWithCustomFields(t *testing.T) {
	c := New("", "", "", "", "")
	extraFields := map[string]string{"environment": "production"}
	c.AddExtraFields(extraFields)
	fields := c.buildMessage("Hello")
	val, ok := fields["criticality"]
	assertEqual(t, true, ok, "criticality should be present")
	assertEqual(t, val, fmt.Sprintf("%d", deferSeverity("INFO")), "val should be 1")
	val, ok = fields["category"]
	assertEqual(t, true, ok, "category should be present")
	assertEqual(t, val, "misc", "val should be misc")
	_, ok = fields["unix_timestamp"]
	assertEqual(t, true, ok, "unix_timestamp should be present")
	val, ok = fields["description"]
	assertEqual(t, true, ok, "ok should be true")
	assertEqual(t, val, "Hello", "description should be present")
	val, ok = fields["environment"]
	assertEqual(t, true, ok, "custom field environment should be present")
	assertEqual(t, val, "production", "val should be production")
}
