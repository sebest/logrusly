# Loggly Hooks for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Usage

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sebest/logrusly"
)

var logglyToken string = "YOUR_LOGGLY_TOKEN"

func main() {
	log := logrus.New()
	hook := logrusly.NewLogglyHook(logglyToken, "www.hostname.com", logrus.WarnLevel, "tag1", "tag2")
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"name": "joe",
		"age":  42,
	}).Error("Hello world!")

	// Flush is automatic for panic/fatal
	// Just make sure to Flush() before exiting or you may loose up to 5 seconds
	// worth of messages.
	hook.Flush()
}
```
