package logging

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	log := New()

	if log == nil {
		t.Fatal("Unexpected nil logging instance")
	}

	if _, ok := log.Formatter.(*logrus.JSONFormatter); !ok {
		t.Fatalf(
			"Unexpected formatter type. Got: %T, Expected: %T",
			log.Formatter,
			&logrus.JSONFormatter{},
		)
	}

	logFile, ok := log.Out.(*os.File)
	if !ok {
		t.Fatalf(
			"Unexpected output type. Got: %T, Expected: %T",
			log.Out,
			&os.File{},
		)
	}

	if logFile.Fd() != os.Stdout.Fd() {
		t.Fatal("Log output not set to stdout")
	}

	if log.Level != logrus.InfoLevel {
		t.Fatalf(
			"Unexpeced log level. Got: %s, Expected: %s",
			log.Level.String(),
			logrus.InfoLevel.String(),
		)
	}
}

func TestLogger_WrappedLogger(t *testing.T) {
	wrappedLogger := New().WrappedLogger()

	if wrappedLogger == nil {
		t.Fatal("Unexpected nil wrapped logging instance")
	}
}
