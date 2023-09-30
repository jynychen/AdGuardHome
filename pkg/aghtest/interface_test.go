package aghtest_test

import (
	"github.com/jynychen/AdGuardHome/pkg/aghtest"
	"github.com/jynychen/AdGuardHome/pkg/filtering"
)

// Put interface checks that cause import cycles here.

// type check
var _ filtering.Resolver = (*aghtest.Resolver)(nil)
