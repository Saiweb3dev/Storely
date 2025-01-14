package utils

// MBToBytes converts megabytes to bytes.
func MBToBytes(mb float64) int64 {
	return int64(mb * 1024 * 1024)
}

// BytesToMB converts bytes to megabytes.
func BytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
