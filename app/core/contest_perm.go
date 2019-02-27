package core

func (c *Contest) CheckListProblems(player *ContestPlayer) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running && player != nil
}

func (c *Contest) CheckViewProblem(player *ContestPlayer, name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running && player != nil
}

func (c *Contest) CheckSubmitProblem(player *ContestPlayer, name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running && player != nil
}
