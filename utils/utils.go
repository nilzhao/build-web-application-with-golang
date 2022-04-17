package utils

func SliceContains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func SliceDiff[T comparable](slice1 []T, slice2 []T) (diffs []T) {
	for _, v := range slice1 {
		if !SliceContains(slice2, v) {
			diffs = append(diffs, v)
		}
	}
	return
}
