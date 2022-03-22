package utils

func Swap[S ~[]E, E any](arr S, i, j int) {
	t := arr[i]
	arr[i] = arr[j]
	arr[j] = t
}
