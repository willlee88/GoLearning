package parse

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseDXDY(s string) (dx, dy int, err error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("want dx,dy")
	}
	dx, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	dy, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	return dx, dy, err
}
