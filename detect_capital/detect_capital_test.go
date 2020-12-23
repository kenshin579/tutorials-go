package detect_capital

import (
	"fmt"
	"testing"
)

func Test_detectCapitalUse1(t *testing.T) {
	use := DetectCapitalUse("USA")
	fmt.Println("use", use)
}
