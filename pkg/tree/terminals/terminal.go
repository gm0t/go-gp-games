package terminals

func extractDepth(extra []int) int {
	if len(extra) > 0 {
		return extra[0]
	}

	return 0
}
