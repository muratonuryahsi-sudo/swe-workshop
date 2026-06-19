//go:build integration

package integration

import "strconv"

func uintToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
