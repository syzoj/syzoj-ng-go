package rediskey

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
)
