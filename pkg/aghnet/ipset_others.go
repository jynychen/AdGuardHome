//go:build !linux

package aghnet

import (
	"github.com/jynychen/AdGuardHome/pkg/aghos"
)

func newIpsetMgr(_ []string) (mgr IpsetManager, err error) {
	return nil, aghos.Unsupported("ipset")
}
