//go:build linux

package aghnet

import (
	"github.com/jynychen/AdGuardHome/pkg/aghos"
)

func defaultHostsPaths() (paths []string) {
	paths = []string{"etc/hosts"}

	if aghos.IsOpenWrt() {
		paths = append(paths, "tmp/hosts")
	}

	return paths
}
