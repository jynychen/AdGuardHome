//go:build darwin || freebsd || openbsd

package aghnet

import "github.com/jynychen/AdGuardHome/pkg/aghos"

func canBindPrivilegedPorts() (can bool, err error) {
	return aghos.HaveAdminRights()
}
