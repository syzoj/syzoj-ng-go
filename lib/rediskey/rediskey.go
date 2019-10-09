package rediskey

import (
	"strings"
	"time"
)

type RedisKey []string

func (k RedisKey) Format(s ...string) string {
	if len(s) != len(k)-1 {
		panic("RedisKey: Format: wrong number of arguments")
	}
	b := &strings.Builder{}
	b.WriteString(k[0])
	for i, v := range s {
		for j := 0; j < len(v); j++ {
			if v[j] == ':' || v[j] == '{' || v[j] == '}' {
				panic("RedisKey: Format: key contains ':', '{', or '}'")
			}
		}
		b.WriteString(v)
		b.WriteString(k[i+1])
	}
	return b.String()
}

var (
	SESSION = RedisKey{"{session:", "}"}

	CORE_QUEUE               = RedisKey{"{core:queue:", "}"}
	CORE_SUBMISSION_PROGRESS = RedisKey{"{core:submission:", "}:progress"}
	CORE_SUBMISSION_DATA     = RedisKey{"{core:submission:", "}:data"}
	CORE_SUBMISSION_CALLBACK = RedisKey{"{core:submission:", "}:callback"}
	CORE_SUBMISSION_RESULT   = RedisKey{"{core:submission:", "}:result"}

	MAIN_PROBLEM_SUBMITS                 = RedisKey{"{main:problem:", "}:submits"}
	MAIN_PROBLEM_ACCEPTS                 = RedisKey{"{main:problem:", "}:accepts"}
	MAIN_USER_LAST_SUBMISSION            = RedisKey{"{main:user:", "}:last_submit"}
	MAIN_USER_LAST_ACCEPT                = RedisKey{"{main:user:", "}:last_accept"}
	MAIN_JUDGE_DONE                      = "{main:judge_done}"
	MAIN_EMAIL_PASSWORD_RECOVERY_RATELIM = RedisKey{"{main:email:", "}:password_recovery_ratelim"}
	MAIN_EMAIL_PASSWORD_RECOVERY_TOKEN   = RedisKey{"{main:email:", "}:password_recovery_token:", ""}
)

// Default expiry time for keys that are no longer active.
const DEFAULT_EXPIRE = time.Second * 86400 * 7 // a week
