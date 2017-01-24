# go-changelog

[![Build Status](https://travis-ci.org/woohgit/go-changelog.svg?branch=master)](https://travis-ci.org/woohgit/go-changelog)

go-changelog is a simple package to send events to a [changelog](https://github.com/prezi/changelog) server. Examples can be found below.

For more detail about the library and its features, reference your local godoc once installed.

Contributions welcome!

## Installation

```bash
go get github.com/woohgit/go-changelog
```

## Available Severities

1. INFO
2. NOTIFICATION
3. WARNING
4. ERROR
5. CRITICAL

## Features

* Custom fields in the message for custom changelog server
* Custom headers
* Support for BasicAuth

## Examples

```Go
package main

import (
	"github.com/woohgit/go-changelog"
)

func main() {
	// if the port is standard 443/80 just pass an empty string ""
	c := changelog.New("https://changelog.yourdomain.tld", "8081", "/api/events", "misc", "INFO")  // Build our new changelog client
    response, err := c.Send("Our server is behind https and listens on the port 8081!")
    if err != nil {
    	// recover
    } else {
    	// everything worked fine
    }

}
```

## Use BasicAuth

```Go
c.UseBasicAuth("username", "password")
response, err := c.Send("Server is behind BasicAuth!")
	if err != nil {
    	// recover
    } else {
    	// everything worked fine
    }
```

## Use Custom Fields

```Go
extraFields := map[string]string{"environment": "production"}
c.AddExtraFields(extraFields)
response, err := c.Send("Our server supports the environment field too!")
	if err != nil {
    	// recover
    } else {
    	// everything worked fine
    }
```

