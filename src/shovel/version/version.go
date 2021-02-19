package version

import (
	"fmt"
)

const (
	Major = "0"
	Minor = "1"
	Build = "3"
)

func Full() string {
	return fmt.Sprintf("%s.%s.%s", Major, Minor, Build)
}
