// prefix defines all Redis key prefixes
package prefix

const (
	SESSION           = "session:"
	PROBLEM_STATEMENT = "problem:stmt:"

	CORE_SUBMISSION_META = "core:submission:meta:"
	CORE_QUEUE = "core:queue:"
	CORE_TIMER_SUBMISSION = "core:timer:submission:"
	CORE_TIMER_SESSION = "core:timer:session:"
	CORE_JUDGER_SESSION = "core:judger:session:"
)
