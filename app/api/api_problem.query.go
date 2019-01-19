package api

const ViewProblemQuery = `
query ViewProblemQuery($problemName: string) {
	problem(func: eq(problemsetentry.short_name, $problemName)) @normalize {
		name: problemsetentry.short_name
		problemsetentry.problem {
			id: problem.id
			title: problem.title
			content: problem.content
		}
	}
}
`
