package rediskey

import (
	"time"
)

// A Redis key template consists of a prefix and a suffix.
type RedisKey [2]string

func (k RedisKey) Format(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' || s[i] == '{' || s[i] == '}' {
			panic("RedisKey: Format: key contains ':', '{', or '}'")
		}
	}
	return k[0] + s + k[1]
}

var (
	SESSION = RedisKey{"{session:", "}"}

	CORE_QUEUE               = RedisKey{"{core:queue:", "}"}
	CORE_SUBMISSION_PROGRESS = RedisKey{"{core:submission:", "}:progress"}
	CORE_SUBMISSION_DATA     = RedisKey{"{core:submission:", "}:data"}
	CORE_SUBMISSION_CALLBACK = RedisKey{"{core:submission:", "}:callback"}
	CORE_SUBMISSION_RESULT   = RedisKey{"{core:submission:", "}:result"}

	MAIN_PROBLEM_SUBMITS      = RedisKey{"{main:problem:", "}:submits"}
	MAIN_PROBLEM_ACCEPTS      = RedisKey{"{main:problem:", "}:accepts"}
	MAIN_USER_LAST_SUBMISSION = RedisKey{"{main:user:", "}:last_submit"}
	MAIN_USER_LAST_ACCEPT     = RedisKey{"{main:user:", "}:last_accept"}
	MAIN_JUDGE_DONE           = "{main:judge_done}"
)

// Default expiry time for keys that are no longer active.
const DEFAULT_EXPIRE = time.Second * 86400 * 7 // a week
