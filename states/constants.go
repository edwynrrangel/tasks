package states

const (
	sqlGetStates = `
	SELECT * FROM task.states WHERE 1 = 1 %s
	`

	sqlGetStateTransitions = `
	SELECT
		st.next_state_id AS id,
		s.name
	FROM task.state_transitions st
	INNER JOIN task.states s 
		ON s.id = st.next_state_id
	WHERE 1 = 1 %s
	`
)
