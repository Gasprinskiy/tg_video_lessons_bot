package slice

func Map[T, R any](slice []T, resultingFunc func(item T) R) []R {
	result := make([]R, 0, len(slice))

	for _, item := range slice {
		result = append(result, resultingFunc(item))
	}

	return result
}
