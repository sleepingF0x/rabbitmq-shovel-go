package version

import (
	"fmt"
)

const (
	Major = "0"
	Minor = "1"
	Build = "2"
)

func Full() string {
	return fmt.Sprintf("%s.%s:%s", Major, Minor, Build)
}
