package bot

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
)

func hash(v interface{}) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprint(v)))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func formatBytes(bytes float64, decimals uint) string {
	if bytes == 0 {
		return "0 B"
	}

	k := 1000.0
	sizes := [9]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	i := math.Floor(math.Log(bytes) / math.Log(k))
	v := strconv.FormatFloat(bytes/math.Pow(k, i), 'f', 2, 64)

	return fmt.Sprintf("%v %v", v, sizes[int(i)])
}
