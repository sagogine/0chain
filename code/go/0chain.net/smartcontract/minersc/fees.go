package minersc

import (
	"errors"
	"fmt"
	"sort"

	"0chain.net/chaincore/block"
	cstate "0chain.net/chaincore/chain/state"
	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/util"

	. "0chain.net/core/logging"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
)

var (
	ErrExecutionStatsNotFound = errors.New("SmartContractExecutionStats stat not found")
)

func (msc *MinerSmartContract) activatePending(mn *MinerNode) {
	for id, pool := range mn.Pending {
		pool.Status = ACTIVE
		mn.Active[id] = pool
		mn.TotalStaked += int64(pool.Balance)
		delete(mn.Pending, id)
	}
}

// pay interests for active pools
func (msc *MinerSmartContract) payInterests(mn *MinerNode, gn *globalNode,
	balances cstate.StateContextI) (err error) {

	if !gn.canMint() {
		return // no mints anymore
	}

	// all active
	for _, pool := range mn.Active {
		var amount = state.Balance(float64(pool.Balance) * gn.InterestRate)
		if amount == 0 {
			continue
		}
		var mint = state.NewMint(ADDRESS, pool.DelegateID, amount)
		if err = balances.AddMint(mint); err != nil {
			return common.NewErrorf("pay_fees/pay_interests",
				"error adding mint for stake %v-%v: %v", mn.ID, pool.ID, err)
		}
		msc.addMint(gn, mint.Amount)
	}

	return
}

func (msc *MinerSmartContract) emptyPool(mn *MinerNode,
	pool *sci.DelegatePool, round int64, balances cstate.StateContextI) (
	resp string, err error) {

	mn.TotalStaked -= int64(pool.Balance)

	// transfer, empty
	var transfer *state.Transfer
	transfer, resp, err = pool.EmptyPool(ADDRESS, pool.DelegateID, nil)
	if err != nil {
		return "", fmt.Errorf("error emptying delegate pool: %v", err)
	}
	if err = balances.AddTransfer(transfer); err != nil {
		return "", fmt.Errorf("adding transfer: %v", err)
	}

	// delete fro muser node
	var un *UserNode
	if un, err = msc.getUserNode(pool.DelegateID, balances); err != nil {
		return "", fmt.Errorf("getting user node: %v", err)
	}
	delete(un.Pools, pool.ID)
	if err = un.save(balances); err != nil {
		return "", fmt.Errorf("saving user node: %v", err)
	}
	return
}

// unlock deleted pools
func (msc *MinerSmartContract) unlockDeleted(mn *MinerNode, round int64,
	balances cstate.StateContextI) (err error) {

	for id := range mn.Deleting {
		var pool = mn.Active[id]
		if _, err = msc.emptyPool(mn, pool, round, balances); err != nil {
			return common.NewError("pay_fees/unlock_deleted", err.Error())
		}
		delete(mn.Active, id)
		delete(mn.Deleting, id)
	}

	return
}

// unlock all delegate pools of offline node
func (msc *MinerSmartContract) unlockOffline(mn *MinerNode,
	balances cstate.StateContextI) (err error) {

	mn.Deleting = make(map[string]*sci.DelegatePool) // reset

	// unlock all pending
	for id, pool := range mn.Pending {
		if _, err = msc.emptyPool(mn, pool, 0, balances); err != nil {
			return common.NewError("pay_fees/unlock_offline", err.Error())
		}
		delete(mn.Pending, id)
	}

	// unlock all active
	for id, pool := range mn.Active {
		if _, err = msc.emptyPool(mn, pool, 0, balances); err != nil {
			return common.NewError("pay_fees/unlock_offline", err.Error())
		}
		delete(mn.Active, id)
	}

	if err = mn.save(balances); err != nil {
		return
	}

	return
}

func (msc *MinerSmartContract) viewChangePoolsWork(gn *globalNode,
	mb *block.MagicBlock, round int64, balances cstate.StateContextI) (
	err error) {

	var miners, sharders *MinerNodes
	if miners, err = msc.getMinersList(balances); err != nil {
		return fmt.Errorf("getting all miners list: %v", err)
	}
	sharders, err = msc.getShardersList(balances, AllShardersKey)
	if err != nil {
		return fmt.Errorf("getting all sharders list: %v", err)
	}

	var (
		mbMiners   = make(map[string]struct{}, mb.Miners.Size())
		mbSharders = make(map[string]struct{}, mb.Miners.Size())

		minersOffline, shardersOffline []*MinerNode
	)

	for _, k := range mb.Miners.Keys() {
		mbMiners[k] = struct{}{}
	}

	for _, k := range mb.Sharders.Keys() {
		mbSharders[k] = struct{}{}
	}

	// miners
	for _, mn := range miners.Nodes {
		if mn, err = msc.getMinerNode(mn.ID, balances); err != nil {
			return fmt.Errorf("missing miner node: %v", err)
		}
		if err = msc.payInterests(mn, gn, balances); err != nil {
			return
		}
		if err = msc.unlockDeleted(mn, round, balances); err != nil {
			return
		}
		msc.activatePending(mn)
		if _, ok := mbMiners[mn.ID]; !ok {
			minersOffline = append(minersOffline, mn)
			continue
		}
		// save excluding offline nodes
		if err = mn.save(balances); err != nil {
			return
		}
	}

	// sharders
	for _, mn := range sharders.Nodes {
		if mn, err = msc.getSharderNode(mn.ID, balances); err != nil {
			return fmt.Errorf("missing sharder node: %v", err)
		}
		if err = msc.payInterests(mn, gn, balances); err != nil {
			return
		}
		if err = msc.unlockDeleted(mn, round, balances); err != nil {
			return
		}
		msc.activatePending(mn)
		if _, ok := mbSharders[mn.ID]; !ok {
			shardersOffline = append(shardersOffline, mn)
			continue
		}
		// save excluding offline nodes
		if err = mn.save(balances); err != nil {
			return
		}
	}

	// unlockOffline
	for _, mn := range minersOffline {
		if err = msc.unlockOffline(mn, balances); err != nil {
			return
		}
	}

	for _, mn := range shardersOffline {
		if err = msc.unlockOffline(mn, balances); err != nil {
			return
		}
	}

	return
}

func (msc *MinerSmartContract) payFees(t *transaction.Transaction,
	inputData []byte, gn *globalNode, balances cstate.StateContextI) (
	resp string, err error) {

	var pn *PhaseNode
	if pn, err = msc.getPhaseNode(balances); err != nil {
		return
	}

	if err = msc.setPhaseNode(balances, pn, gn); err != nil {
		return "", common.NewErrorf("pay_fees",
			"error inserting phase node: %v", err)
	}

	var block = balances.GetBlock()
	if block.Round == gn.ViewChange && !msc.SetMagicBlock(balances) {
		return "", common.NewErrorf("pay_fee",
			"can't set magic block round=%d viewChange=%d",
			block.Round, gn.ViewChange)
	}

	if t.ClientID != block.MinerID {
		return "", common.NewError("pay_fee", "not block generator")
	}

	if block.Round <= gn.LastRound {
		return "", common.NewError("pay_fee", "jumped back in time?")
	}

	// the block generator
	var mn *MinerNode
	if mn, err = msc.getMinerNode(block.MinerID, balances); err != nil {
		return "", common.NewErrorf("pay_fee", "can't get generator: %v", err)
	}

	var (
		// block reward -- mint for the block
		blockReward = state.Balance(
			float64(gn.BlockReward) * gn.RewardRate,
		)
		charger, restr   = mn.splitByServiceCharge(blockReward)
		minerr, sharderr = gn.splitByShareRatio(restr)
		// fees         -- total fees for the block
		fees             = msc.sumFee(block, true)
		chargef, restf   = mn.splitByServiceCharge(fees)
		minerf, sharderf = gn.splitByShareRatio(restf)
		// intermediate response
		iresp string
	)

	// pay for the generator (charge + share ratio)
	iresp, err = msc.mintStakeHolders(minerr+charger, mn, gn, false, balances)
	if err != nil {
		return "", err
	}
	resp += iresp
	// mint for the generator (charge + share ratio)
	iresp, err = msc.payStakeHolders(minerf+chargef, mn, false, balances)
	if err != nil {
		return "", err
	}
	resp += iresp
	// pay and mint rest for block sharders
	iresp, err = msc.paySharders(sharderf, sharderr, block, gn, balances)
	if err != nil {
		return "", err
	}
	resp += iresp

	// view change stuff
	if block.Round == gn.ViewChange {
		var mb = balances.GetBlock().MagicBlock
		err = msc.viewChangePoolsWork(gn, mb, block.Round, balances)
		if err != nil {
			return "", err
		}
	}

	if err = mn.save(balances); err != nil {
		return "", common.NewErrorf("pay_fees",
			"saving generator node: %v", err)
	}

	gn.setLastRound(block.Round)
	if err = gn.save(balances); err != nil {
		return "", common.NewErrorf("pay_fees",
			"saving global node: %v", err)
	}

	return resp, nil
}

func (msc *MinerSmartContract) sumFee(b *block.Block,
	updateStats bool) state.Balance {

	var totalMaxFee int64
	var feeStats metrics.Counter
	if stat := msc.SmartContractExecutionStats["feesPaid"]; stat != nil {
		feeStats = stat.(metrics.Counter)
	}
	for _, txn := range b.Txns {
		totalMaxFee += txn.Fee
	}

	if updateStats && feeStats != nil {
		feeStats.Inc(totalMaxFee)
	}
	return state.Balance(totalMaxFee)
}

func (msc *MinerSmartContract) mintStakeHolders(value state.Balance,
	node *MinerNode, gn *globalNode, isSharder bool,
	balances cstate.StateContextI) (resp string, err error) {

	if !gn.canMint() {
		return // can't mint anymore
	}

	if value == 0 {
		return // nothing to mint
	}

	if isSharder {
		node.Stat.SharderRewards += value
	} else {
		node.Stat.GeneratorRewards += value
	}

	var totalStaked = node.TotalStaked

	for _, pool := range node.orderedActivePools() {
		var (
			ratio    = float64(pool.Balance) / float64(totalStaked)
			userMint = state.Balance(float64(value) * ratio)
		)

		Logger.Info("mint delegate",
			zap.Any("pool", pool),
			zap.Any("mint", userMint))

		if userMint == 0 {
			continue // avoid insufficient minting
		}

		var mint = state.NewMint(ADDRESS, pool.DelegateID, userMint)
		if err = balances.AddMint(mint); err != nil {
			resp += fmt.Sprintf("pay_fee/minting - adding mint: %v", err)
			continue
		}
		msc.addMint(gn, mint.Amount)

		pool.TotalPaid += mint.Amount

		if pool.High < mint.Amount {
			pool.High = mint.Amount
		}

		if pool.Low == -1 || pool.Low > mint.Amount {
			pool.Low = mint.Amount
		}

		resp += string(mint.Encode())
	}

	return resp, nil
}

func (msc *MinerSmartContract) payStakeHolders(value state.Balance,
	node *MinerNode, isSharder bool,
	balances cstate.StateContextI) (resp string, err error) {

	if value == 0 {
		return // nothing to pay
	}

	if isSharder {
		node.Stat.SharderFees += value
	} else {
		node.Stat.GeneratorFees += value
	}

	var totalStaked = node.TotalStaked

	for _, pool := range node.orderedActivePools() {
		var (
			ratio   = float64(pool.Balance) / float64(totalStaked)
			userFee = state.Balance(float64(value) * ratio)
		)

		Logger.Info("pay delegate",
			zap.Any("pool", pool),
			zap.Any("fee", userFee))

		if userFee == 0 {
			continue // avoid insufficient transfer
		}

		var transfer = state.NewTransfer(ADDRESS, pool.DelegateID, userFee)
		if err = balances.AddTransfer(transfer); err != nil {
			return "", fmt.Errorf("adding transfer: %v", err)
		}

		pool.TotalPaid += transfer.Amount
		pool.NumRounds++

		if pool.High < transfer.Amount {
			pool.High = transfer.Amount
		}

		if pool.Low == -1 || pool.Low > transfer.Amount {
			pool.Low = transfer.Amount
		}

		resp += string(transfer.Encode())
	}

	return resp, nil
}

func (msc *MinerSmartContract) getBlockSharders(block *block.Block,
	balances cstate.StateContextI) (sharders []*MinerNode, err error) {

	var sids = balances.GetBlockSharders(block.PrevBlock)
	sort.Strings(sids)

	sharders = make([]*MinerNode, 0, len(sids))

	for _, sid := range sids {
		var sn *MinerNode
		sn, err = msc.getSharderNode(sid, balances)
		if err != nil && err != util.ErrValueNotPresent {
			return nil, fmt.Errorf("unexpected error: %v", err)
		}
		sharders, err = append(sharders, sn), nil // even if it's nil, reset err
	}

	return
}

// pay fees and mint sharders
func (msc *MinerSmartContract) paySharders(fee, mint state.Balance,
	block *block.Block, gn *globalNode, balances cstate.StateContextI) (
	resp string, err error) {

	var sharders []*MinerNode
	if sharders, err = msc.getBlockSharders(block, balances); err != nil {
		return // unexpected error
	}

	// fess and mint
	var (
		partf = state.Balance(float64(fee) / float64(len(sharders)))
		partm = state.Balance(float64(mint) / float64(len(sharders)))
	)

	// part for every sharder
	for _, sh := range sharders {

		var sresp string
		sresp, err = msc.payStakeHolders(partf, sh, true, balances)
		if err != nil {
			return "", common.NewErrorf("pay_fees/pay_sharders",
				"paying block sharder fees: %v", err)
		}

		resp += sresp

		sresp, err = msc.mintStakeHolders(partm, sh, gn, true, balances)
		if err != nil {
			return "", common.NewErrorf("pay_fees/mint_sharders",
				"minting block sharder reward: %v", err)
		}

		resp += sresp

		if err = sh.save(balances); err != nil {
			return "", common.NewErrorf("pay_fees/pay_sharders",
				"saving sharder node: %v", err)
		}
	}

	return
}
