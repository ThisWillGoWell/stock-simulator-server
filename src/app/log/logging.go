package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

var Alerts = logrus.New()
var Log = logrus.New()

type Fields logrus.Fields

func Init(alertOut io.Writer) {
	Alerts.SetOutput(alertOut)
	Alerts.SetFormatter(&logrus.TextFormatter{})
}
