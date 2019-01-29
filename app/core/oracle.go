package core

type KeyUserName string

type OracleLock struct {
	locked int32
	core   *Core
	keys   []interface{}
}

func (c *Core) initOracle() {
	c.oracle = make(map[interface{}]struct{})
}

// Locks the specified keys at once.
// The keys are either all locked or none locked.
// Returns an *OracleLock if success, otherwise nil.
// It is recommended to defer a call to *OracleLock.Release() immediately after the call.
// Afer the call the slice is owned by the oracle.
func (c *Core) LockOracle(keys []interface{}) *OracleLock {
	c.oracleLock.Lock()
	defer c.oracleLock.Unlock()
	for k, v := range keys {
		_, found := c.oracle[v]
		if found {
			// rollback
			for k2, v2 := range keys {
				if k2 == k {
					break
				}
				delete(c.oracle, v2)
			}
			return nil
		}
		c.oracle[v] = struct{}{}
	}
	return &OracleLock{core: c, keys: keys}
}

func (o *OracleLock) Release() {
	if o.core != nil {
		o.core.oracleLock.Lock()
		defer o.core.oracleLock.Unlock()
		for _, v := range o.keys {
			delete(o.core.oracle, v)
		}
        o.core = nil
	}
}
