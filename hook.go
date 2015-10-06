package logrusly

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/segmentio/go-loggly"
)

// LogglyHook to send logs to the Loggly service.
type LogglyHook struct {
	client *loggly.Client
	host   string
	levels []logrus.Level
}

// NewLogglyHook creates a Loogly hook to be added to an instance of logger.
func NewLogglyHook(token string, host string, level logrus.Level, tags ...string) *LogglyHook {
	client := loggly.New(token, tags...)
	client.Defaults = loggly.Message{}

	// sigc := make(chan os.Signal, 1)
	// signal.Notify(sigc,
	// 	syscall.SIGHUP,
	// 	syscall.SIGINT,
	// 	syscall.SIGTERM,
	// 	syscall.SIGQUIT)
	// go func() {
	// 	s := <-sigc
	// 	if s != nil {
	// 		client.Flush()
	// 	}
	// }()

	levels := []logrus.Level{}
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	return &LogglyHook{
		client: client,
		host:   host,
		levels: levels,
	}
}

// Tag exposes the go-loggly .Tag() functionality
func (hook *LogglyHook) Tag(tags string) {
	hook.client.Tag(tags)
}

// Fire sends the event to Loggly
func (hook *LogglyHook) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	level := entry.Level.String()

	logglyMessage := loggly.Message{
		"timestamp": entry.Time.UTC().Format(time.RFC3339Nano),
		"level":     strings.ToUpper(level),
		"message":   entry.Message,
		"host":      hook.host,
		"data":      data,
	}

	err := hook.client.Send(logglyMessage)
	if err != nil {
		log := logrus.New()
		log.WithFields(logrus.Fields{
			"source": "loggly",
			"error":  err.Error(),
		}).Warn("Failed to send error to Loggly")
		return err
	}

	if level == "fatal" || level == "panic" {
		hook.Flush()
	}

	return nil
}

// Flush sends buffered events to Loggly.
func (hook *LogglyHook) Flush() {
	hook.client.Flush()
}

// Levels returns the list of logging levels that we want to send to Loggly.
func (hook *LogglyHook) Levels() []logrus.Level {
	return hook.levels
}
