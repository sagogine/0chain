package interestpoolsc

import (
	"context"
	"time"
	// "encoding/json"
	"net/url"

	"0chain.net/core/common"
)

func (ip *InterestPoolSmartContract) getPoolsStats(ctx context.Context, params url.Values) (interface{}, error) {
	un, err := ip.getUserNode(params.Get("client_id"))
	if err != nil {
		return nil, err
	}
	if len(un.Pools) == 0 {
		return common.NewError("failed to get stats", "no pools exist").Error(), nil
	}
	t := time.Now()
	stats := &poolStats{}
	for _, pool := range un.Pools {
		stat, err := ip.getPoolStats(pool, t)
		if err != nil {
			return "crap this shouldn't happen", nil
		}
		stats.addStat(stat)
	}
	return string(stats.encode()), nil
}

func (ip *InterestPoolSmartContract) getPoolStats(pool *interestPool, t time.Time) (*poolStat, error) {
	stat := &poolStat{}
	statBytes := pool.LockStats(t)
	err := stat.decode(statBytes)
	if err != nil {
		return nil, err
	}
	stat.ID = pool.ID
	stat.Locked = pool.IsLocked(t)
	stat.Balance = pool.Balance
	stat.InterestRate = pool.InterestRate
	stat.InterestEarned = pool.InterestEarned
	return stat, nil
}
