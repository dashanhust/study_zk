package core

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

const timeFormat = "Mon Jan 02 15:04:05 GMT 2006"

func fmtTime(t int64) string {
	return time.Unix(t/1000, 0).UTC().Format(timeFormat)
}

func fmtStat(stat *zk.Stat) string {
	return fmt.Sprintf(`Czxid = 0x%x
Ctime = %s
Mzxid = 0x%x
Mtime = %s
Pzxid = 0x%x
Cversion = %d
Version = %d
Aversion = %d
EphemeralOwner = 0x%x
DataLength = %d
NumberChildren = %d`,
		stat.Czxid,
		fmtTime(stat.Ctime),
		stat.Mzxid,
		fmtTime(stat.Mtime),
		stat.Pzxid,
		stat.Cversion,
		stat.Version,
		stat.Aversion,
		stat.EphemeralOwner,
		stat.DataLength,
		stat.NumChildren)
}
