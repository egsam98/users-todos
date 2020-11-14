package time

import (
	"time"
)

// Конвертация UNIX-time в time.Time без учета временной зоны
func UtcFromUnix(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}
