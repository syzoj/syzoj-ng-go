package api

const sessionQuery = `
query sessionQuery($token: string) {
	session(func: eq(session.token, $token)) {
		uid
		session.auth_user {
			uid
			user.username
		}
	}
}
`

const sessionByUidQuery = `
query sessionByUidQuery($sessUid: string) {
	session(func: uid($sessUid)) {
		uid
		session.auth_user {
			uid
			user.username
		}
	}
}
`

const CheckUserNameQuery = `
query CheckUserName($userName: string) {
	user(func: eq(user.username, $userName)) {
		uid
	}
}
`

const LoginQuery = `
query Login($userName: string, $password: string) {
	user(func: eq(user.username, $userName)) {
		uid
		check: checkpwd(user.password, $password)
	}
}
`

const MyProblemQuery = `
query MyProblem($userId: string) {
	problems(func: uid($userId)) @normalize {
		~problem.owner {
			id: problem.id
			title: problem.title
			create_time: problem.create_time
		}
	}
}
`

const ViewProblemDbQuery = `
query ViewProblemDbQuery($problemId: string, $userId: string) {
	problem(func: eq(problem.id, $problemId)) @normalize {
		uid: uid
		title: problem.title@.
		statement: problem.statement@.
		token: problem.token
		owner: problem.owner {
			owner_uid: uid
		}
	}
	problemset(func: eq(problemset.name, "public")) @filter(uid_in(problemset.owner, $userId)) {
		uid: uid
	}
}
`

const CheckProblemPublicizeQuery = `
query CheckProblemPublicizeQuery($problemId: string, $userId: string, $name: string) {
	problem(func: eq(problem.id, $problemId)) @normalize {
		uid
		owner: problem.owner {
			owner_uid: uid
		}
	}
	problemset(func: eq(problemset.name, "public")) @filter(uid_in(problemset.owner, $userId)) {
		uid: uid
	}
	name(func: eq(problemsetentry.short_name, $name)) {
		uid: uid
	}
}
`

const CheckProblemCanSubmitQuery = `
query CheckProblemCanSubmitQuery($problemId: string) {
	problem(func: eq(problem.id, $problemId)) {
		uid
	}
}
`

const ProblemsQuery = `
{
	problems(func: has(problemsetentry.short_name)) @normalize {
		name: problemsetentry.short_name
		problemsetentry.problem {
			title: problem.title@.
			id: problem.id
			tags: problem.tags
			submit_num: count(~submission.problemsetentry)
		}
	}
}
`
const SubmissionViewQuery = `
query SubmissionViewQuery($submissionId: string) {
	submission(func: eq(submission.id, $submissionId)) @normalize {
		status: submission.status
		message: submission.message
		score: submission.score
		language: submission.language
		code: submission.code
		submit_time: submission.submit_time
		submission.owner {
			submitter_name: user.username
		}
		submission.problem {
			problem_id: problem.id
			problem_title: problem.title
		}
	}
}
`

const MySubmissionQuery = `
query MySubmissionQuery($userId: string) {
	submissions(func: uid($userId)) @normalize {
		~submission.owner {
			submission_id: submission.id
			submission_status: submission.status
			submit_time: submission.submit_time
			submission.problem {
				problem_id: problem.id
				problem_title: problem.title
			}
		}
	}
}
`
