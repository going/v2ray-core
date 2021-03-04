package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	BYTE = 1.0 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

func GetDetectedSize(bytes int64) string {
	unit := ""
	value := float32(bytes)

	switch {

	case bytes >= TERABYTE:
		unit = "TB"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "GB"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KB"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
		value = value / BYTE
	case bytes == 0:
		return "0"

	}

	result := fmt.Sprintf("%.2f", value)
	result = strings.TrimSuffix(result, ".00")
	return fmt.Sprintf("%s%s", result, unit)
}

func MD5(text string) []byte {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return ctx.Sum(nil)
}

func GetRandomString(len1 int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len1; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
