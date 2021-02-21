// build any argument from array for postgres
func buildAnyArgumentFromArray(args interface{}) interface{} {
	return pq.Array(args)
}
