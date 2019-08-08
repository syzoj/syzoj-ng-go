package main

func (app *App) tryLock(key string, val string) (string, bool) {
	app.locksMu.Lock()
	defer app.locksMu.Unlock()
	oval, ok := app.locks[key]
	if ok {
		return oval, false
	}
	app.locks[key] = val
	return "", true
}

func (app *App) unlock(key string) {
	app.locksMu.Lock()
	defer app.locksMu.Unlock()
	delete(app.locks, key)
}
