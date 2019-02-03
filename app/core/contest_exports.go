package core

func (c *Contest) Unlock() {
	c.lock.Unlock()
}

func (c *Contest) RUnlock() {
	c.lock.RUnlock()
}

func (c *Contest) Running() bool {
	return c.running
}
