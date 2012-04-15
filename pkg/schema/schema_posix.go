//+build linux darwin netbsd freebsd openbsd
//+build !appengine

package schema

import (
	"os"
	"syscall"
	"time"
)

var _ = syscall.Close

func init() {
	populateSchemaStat = append(populateSchemaStat, populateSchemaUnix)
}

func populateSchemaUnix(m map[string]interface{}, fi os.FileInfo) {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return
	}
	m["unixOwnerId"] = st.Uid
	if user := getUserFromUid(int(st.Uid)); user != "" {
		m["unixOwner"] = user
	}
	m["unixGroupId"] = st.Gid
	if group := getGroupFromGid(int(st.Gid)); group != "" {
		m["unixGroup"] = group
	}

	// Include the ctime too, if it differs.
	sec, nsec := st.Ctim.Unix()
	ctime := time.Unix(sec, nsec)
	if sec != 0 && !ctime.Equal(fi.ModTime()) {
		m["unixCtime"] = RFC3339FromTime(ctime)
	}
}