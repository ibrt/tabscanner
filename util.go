package tabscanner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	moneyRegexp = regexp.MustCompile(`^([0-9]+)([\.,:][0-9]+)?$`)
)

// ParseMoney parses a money amount assuming the given number of decimals. It returns an integer representing a fixed
// point money value (e.g. ("10.010", 2) -> 1001). Excess decimals are truncated.
func ParseMoney(v string, decimals int) (int64, error) {
	matches := moneyRegexp.FindAllStringSubmatch(v, -1)
	if len(matches) != 1 {
		return -1, fmt.Errorf("invalid value '%v'", v)
	}

	wholePart := matches[0][1]
	decPart := matches[0][2]

	// no decimal part
	if decPart == "" {
		whole, err := strconv.ParseInt(wholePart, 10, 64)
		if err != nil {
			return -1, err
		}

		for i := 0; i < decimals; i++ {
			whole *= 10
		}

		return whole, nil
	}

	decPart = decPart[1:] // remove the comma

	if len(decPart) > decimals {
		decPart = decPart[:decimals]
	}

	if len(decPart) < decimals {
		decPart = decPart + strings.Repeat("0", decimals-len(decPart))
	}

	res, err := strconv.ParseInt(wholePart+decPart, 10, 64)
	if err != nil {
		return -1, err
	}

	return res, nil
}
