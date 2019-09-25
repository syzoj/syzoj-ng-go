package app

func (a *App) checkRedis(data interface{}, err error) {
	if err != nil {
		log.WithError(err).Warning("failure while executing redis command")
	}
}
