// @Title tool.go
// @Description $END$
// @Author Zero - 2024/8/11 01:00:18

package tool

func Contains[T comparable] (slice []T, target T) bool {
	for _, item := range slice {
		if item == target{
			return true
		}
	}
	return false
}
