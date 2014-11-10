package logrusly

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	// Use this fork until the fixes are accepted upstream https://github.com/segmentio/go-loggly/pull/6
	"github.com/sebest/go-loggly"
)

// LogglyHook to send log messages to the Loggly API. You must set:
type LogglyHook struct {
	client *loggly.Client
	host   string
}

func NewLogglyHook(token string, host string, tags ...string) *LogglyHook {
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

	return &LogglyHook{
		client: client,
		host:   host,
	}
}

func (hook *LogglyHook) Fire(entry *logrus.Entry) error {
	level := entry.Level.String()
	logglyMessage := loggly.Message{
		"timestamp": entry.Time.UTC().Format(time.RFC3339Nano),
		"level":     strings.ToUpper(level),
		"message":   entry.Message,
		"host":      hook.host,
		"data":      entry.Data,
	}

	err := hook.client.Send(logglyMessage)
	if err != nil {
		log := logrus.New()
		log.WithFields(logrus.Fields{
			"source": "loggly",
			"error":  err,
		}).Warn("Failed to send error to Loggly")
		return err
	}

	if level == "fatal" || level == "panic" {
		hook.Flush()
	}

	return nil
}

func (hook *LogglyHook) Flush() {
	hook.client.Flush()
}

func (hook *LogglyHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
