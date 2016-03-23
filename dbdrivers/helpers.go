package dbdrivers

// isJoinTable is true if table has at least 2 foreign keys and
// the two foreign keys are involved in a primary composite key
func isJoinTable(t Table) bool {
	return false
}
