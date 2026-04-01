package server

import "github.com/stockyard-dev/stockyard-almanac/internal/license"

type Limits struct {
	MaxEntries    int  // 0 = unlimited
	EmailNotify   bool // notify subscribers on publish
	CustomDomain  bool
	RetentionDays int
}

var freeLimits = Limits{
	MaxEntries:   20,
	EmailNotify:  false,
	CustomDomain: false,
	RetentionDays: 365,
}

var proLimits = Limits{
	MaxEntries:   0,
	EmailNotify:  true,
	CustomDomain: true,
	RetentionDays: 365,
}

func LimitsFor(info *license.Info) Limits {
	if info != nil && info.IsPro() {
		return proLimits
	}
	return freeLimits
}

func LimitReached(limit, current int) bool {
	if limit == 0 {
		return false
	}
	return current >= limit
}
