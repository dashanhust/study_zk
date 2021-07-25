package core

import (
	"fmt"
	"path"
	"strings"
)

// 功能：去除路径中最后面的符号 /
// 如果是该路径本身为根路径 / 那么就保留
func cleanPath(p string) string {
	if p == "/" {
		return p
	}
	return strings.TrimSuffix(p, "/")
}

// 功能：将 节点路径root 和 节点名child 分别提取出来
// 例如： /a/b/c/d 那么root为 /a/b/c child为 d
func splitPath(p string) (root, child string) {
	root, child = path.Split(p)
	root = strings.TrimSuffix(root, "/")
	if root == "" {
		root = "/"
	}
	return
}

// 功能：打印帮助信息
func printHelp() {
	fmt.Println(`get <path>
ls <path>
create <path> [<data>]
set <path> [<data>]
delete <path>
connect <host:port>
addauth <scheme> <auth>
close
exit`)
}

// 功能：打印错误信息
func printRunError(err error) {
	fmt.Println(err)
}
