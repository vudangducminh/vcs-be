package function

import (
	"errors"
)

func Divide(a int, b int) (float32, error) {
	if b == 0 {
		return 0, errors.New("Divide by 0")
	}
	var ans float32 = float32(a) / float32(b)
	return ans, nil
}
