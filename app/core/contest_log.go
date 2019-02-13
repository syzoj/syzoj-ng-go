package core

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type contestHook struct {
	*Contest
}

func (c contestHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel}
}

func (c contestHook) Fire(e *logrus.Entry) error {
	log.WithFields(e.Data).WithField("contestId", c.id).WithTime(e.Time).Log(e.Level, e.Message)
	return nil
}

func (c *Contest) initLogger() {
	c.log = logrus.New()
	c.log.Out = ioutil.Discard
	c.log.SetLevel(log.Level)
	c.log.AddHook(contestHook{c})
}
