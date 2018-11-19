package sharder

import (
	"context"
	"strconv"

	"0chain.net/datastore"
	"0chain.net/ememorystore"
	"0chain.net/round"
)

/*StoreRound - persists given round to ememory(rocksdb)*/
func (sc *Chain) StoreRound(ctx context.Context, r *round.Round) error {
	roundEntityMetadata := r.GetEntityMetadata()
	rctx := ememorystore.WithEntityConnection(ctx, roundEntityMetadata)
	defer ememorystore.Close(rctx)
	err := r.Write(rctx)
	if err != nil {
		return err
	}
	con := ememorystore.GetEntityCon(rctx, roundEntityMetadata)
	err = con.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (sc *Chain) GetMostRecentRoundFromDB(ctx context.Context) (*round.Round, error) {
	remd := datastore.GetEntityMetadata("round")
	rctx := ememorystore.WithEntityConnection(ctx, remd)
	defer ememorystore.Close(rctx)
	c := ememorystore.GetEntityCon(rctx, remd)
	r := remd.Instance().(*round.Round)
	iterator := c.Conn.NewIterator(c.ReadOptions)
	defer iterator.Close()
	iterator.SeekToLast()
	if !iterator.Valid() {
		roundNum, err := strconv.ParseInt(string(iterator.Key().Data()), 10, 64)
		if err == nil {
			r.Number = roundNum
			return r, nil
		}
		return r, err
	}
	return r, iterator.Err()
}
