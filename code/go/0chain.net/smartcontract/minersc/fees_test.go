package minersc

import (
	"math/rand"
	"testing"

	"0chain.net/chaincore/block"
	"0chain.net/chaincore/node"
	"0chain.net/chaincore/state"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestClient struct {
	client   *Client
	delegate *Client
	stakers  []*Client
}

func createLFMB(miners []*TestClient, sharders []*TestClient) (
	b *block.Block) {

	b = new(block.Block)

	b.MagicBlock = block.NewMagicBlock()
	b.MagicBlock.Miners = node.NewPool(node.NodeTypeMiner)
	b.MagicBlock.Sharders = node.NewPool(node.NodeTypeSharder)

	for _, miner := range miners {
		b.MagicBlock.Miners.NodesMap[miner.client.id] = new(node.Node)
	}
	for _, sharder := range sharders {
		b.MagicBlock.Sharders.NodesMap[sharder.client.id] = new(node.Node)
	}
	return
}

func (msc *MinerSmartContract) setDKGMiners(t *testing.T,
	miners []*TestClient, balances *testBalances) {

	t.Helper()

	var global, err = msc.getGlobalNode(balances)
	require.NoError(t, err)

	var nodes *DKGMinerNodes
	nodes, err = msc.getMinersDKGList(balances)
	require.NoError(t, err)

	nodes.setConfigs(global)
	for _, miner := range miners {
		nodes.SimpleNodes[miner.client.id] = &SimpleNode{ID: miner.client.id}
		nodes.Waited[miner.client.id] = true
	}

	_, err = balances.InsertTrieNode(DKGMinersKey, nodes)
	require.NoError(t, err)
}

func Test_payFees(t *testing.T) {
	const sharderStakeValue, minerStakeValue, generatorStakeValue = 17, 19, 23
	const sharderStakersAmount, minerStakersAmount, generatorStakersAmount = 3, 5, 7
	const generatorIdx = 0

	var (
		balances = newTestBalances()
		msc      = newTestMinerSC()
		now      int64
		err      error

		miners   []*TestClient
		sharders []*TestClient
		generator  *TestClient
	)

	setConfig(t, balances)

	t.Run("add miners", func(t *testing.T) {
		generator = newClientWithStakers(true, t, msc, now,
			generatorStakersAmount, generatorStakeValue, balances)

		for idx := 0; idx < 10; idx++ {
			if idx == generatorIdx {
				miners = append(miners, generator)
			} else {
				miners = append(miners, newClientWithStakers(true, t, msc, now,
					minerStakersAmount, minerStakeValue, balances))
			}
			now += 10
		}
	})

	t.Run("add sharders", func(t *testing.T) {
		for idx := 0; idx < 10; idx++ {
			sharders = append(sharders, newClientWithStakers(false, t, msc, now,
				sharderStakersAmount, sharderStakeValue, balances))
			now += 10
		}
	})

	msc.setDKGMiners(t, miners, balances)
	balances.setLFMB(createLFMB(miners, sharders))

    for idx, miner := range miners {
        var stakeValue int64
        if idx == generatorIdx {
            stakeValue = generatorStakeValue
        } else {
            stakeValue = minerStakeValue
        }

        for _, staker := range miner.stakers {
            _, err = staker.callAddToDelegatePool(t, msc, now,
                stakeValue, miner.client.id, balances)

            require.NoError(t, err, "staking miner")
            assert.Zero(t, balances.balances[staker.id], "stakers' balances should be updated later")

            now += 10
        }

        assert.Zero(t, balances.balances[miner.client.id], "miner's balance shouldn't be changed yet")
        assert.Zero(t, balances.balances[miner.delegate.id], "miner's delegate balance shouldn't be changed yet")
    }

    for _, sharder := range sharders {
        for _, staker := range sharder.stakers {
            _, err = staker.callAddToDelegatePool(t, msc, now,
                sharderStakeValue, sharder.client.id, balances)

            require.NoError(t, err, "staking sharder")
            assert.Zero(t, balances.balances[staker.id], "stakers' balance should be updated later")

            now += 10
        }

        assert.Zero(t, balances.balances[sharder.client.id], "sharder's balance shouldn't be changed yet")
        assert.Zero(t, balances.balances[sharder.delegate.id], "sharder's balance shouldn't be changed yet")
    }

	msc.setDKGMiners(t, miners, balances)

	//t.Run("pay fees -> view change", func(t *testing.T) {
	//	for id, bal := range balances.balances {
	//		if id == ADDRESS {
	//			continue
	//		}
	//		assert.Zerof(t, bal, "unexpected balance: %s", id)
	//	}
	//
	//	setRounds(t, msc, 250, 251, balances)
	//	setMagicBlock(t, unwrapClients(miners), unwrapClients(sharders),
	//		balances)
	//
	//	var generator, blck = prepareGeneratorAndBlock(miners, 0, 251)
	//
	//	// payFees transaction
	//	now += 10
	//	var tx = newTransaction(generator.miner.id, ADDRESS, 0, now)
	//	balances.txn = tx
	//	balances.block = blck
	//	balances.blockSharders = selectRandom(sharders, 3)
	//
	//	var global, err = msc.getGlobalNode(balances)
	//	require.NoError(t, err, "getting global node")
	//
	//	_, err = msc.payFees(tx, nil, gn, balances)
	//	require.NoError(t, err, "pay_fees error")
	//
	//	// pools becomes active, nothing should be payed
	//
	//	for _, mn := range miners {
	//		assert.Zero(t, balances.balances[mn.miner.id],
	//			"miner balance")
	//		assert.Zero(t, balances.balances[mn.delegate.id],
	//			"miner delegate balance?")
	//		for _, st := range mn.stakers {
	//			assert.Zero(t, balances.balances[st.id], "stake balance?")
	//		}
	//	}
	//	for _, sh := range sharders {
	//		assert.Zero(t, balances.balances[sh.sharder.id],
	//			"sharder balance")
	//		assert.Zero(t, balances.balances[sh.delegate.id],
	//			"sharder delegate balance?")
	//		for _, st := range sh.stakers {
	//			assert.Zero(t, balances.balances[st.id], "stake balance?")
	//		}
	//	}
	//
	//	global, err = msc.getGlobalNode(balances)
	//	require.NoError(t, err, "can't get global node")
	//	assert.EqualValues(t, 251, global.LastRound)
	//	assert.EqualValues(t, 0, global.Minted)
	//})

	msc.setDKGMiners(t, miners, balances)

	t.Run("pay fees -> no fees", func(t *testing.T) {
		for id, bal := range balances.balances {
			if id == ADDRESS {
				continue
			}
			assert.Zerof(t, bal, "unexpected balance: %s", id)
		}

		setRounds(t, msc, 251, 501, balances)

		var generator, blck = prepareGeneratorAndBlock(miners, 0, 252)

		// payFees transaction
		now += 10
		var tx = newTransaction(generator.client.id, ADDRESS, 0, now)
		balances.txn = tx
		balances.block = blck
		balances.blockSharders = selectRandom(sharders, 3)

		var global, err = msc.getGlobalNode(balances)
		require.NoError(t, err, "getting global node")

		_, err = msc.payFees(tx, nil, global, balances)
		require.NoError(t, err, "pay_fees error")

		// pools active, no fees, rewards should be payed for
		// generator's and block sharders' stake holders

		var (
			expected = make(map[string]state.Balance)
			got      = make(map[string]state.Balance)
		)

		for _, miner := range miners {
			assert.Zero(t, balances.balances[miner.client.id])
			assert.Zero(t, balances.balances[miner.delegate.id])

			var val state.Balance = 0;
			if miner == generator {
				val = 77e7;
			}

			for _, staker := range miner.stakers {
				expected[staker.id] += val;
				got[staker.id] = balances.balances[staker.id]
			}
		}

		assert.Equal(t, len(expected), len(got), "sizes of balance maps")
		assert.Equal(t, expected, got, "balances")

		for _, sharder := range sharders {
			assert.Zero(t, balances.balances[sharder.client.id])
			assert.Zero(t, balances.balances[sharder.delegate.id])
			for _, st := range sharder.stakers {
				expected[st.id] += 0
				got[st.id] = balances.balances[st.id]
			}
		}

		for _, sharder := range filterClientsById(sharders, balances.blockSharders) {
			for _, staker := range sharder.stakers {
				expected[staker.id] += 21e7
			}
		}

		assert.Equal(t, len(expected), len(got), "sizes of balance maps")
		assert.Equal(t, expected, got, "balances")
	})

	// don't set DKG miners list, because no VC is expected

	// reset all balances
	balances.balances = make(map[string]state.Balance)

	//t.Run("pay fees -> with fees", func(t *testing.T) {
	//	setRounds(t, msc, 252, 501, balances)
	//
	//	var generator, blck = prepareGeneratorAndBlock(miners, 0, 253)
	//
	//	// payFees transaction
	//	now += 10
	//	var tx = newTransaction(generator.miner.id, ADDRESS, 0, now)
	//	balances.txn = tx
	//	balances.block = blck
	//	balances.blockSharders = selectRandom(sharders, 3)
	//
	//	// add fees
	//	tx.Fee = 100e10
	//	blck.Txns = append(blck.Txns, tx)
	//
	//	var global, err = msc.getGlobalNode(balances)
	//	require.NoError(t, err, "getting global node")
	//
	//	_, err = msc.payFees(tx, nil, global, balances)
	//	require.NoError(t, err, "pay_fees error")
	//
	//	// pools are active, rewards as above and +fees
	//
	//	var (
	//		expected = make(map[string]state.Balance)
	//		got      = make(map[string]state.Balance)
	//	)
	//
	//	for _, mn := range miners {
	//		assert.Zero(t, balances.balances[mn.miner.id])
	//		assert.Zero(t, balances.balances[mn.delegate.id])
	//		for _, st := range mn.stakers {
	//			if mn == generator {
	//				expected[st.id] += 77e7 + 11e10 // + generator fees
	//			} else {
	//				expected[st.id] += 0
	//			}
	//			got[st.id] = balances.balances[st.id]
	//		}
	//	}
	//
	//	for _, sh := range sharders {
	//		assert.Zero(t, balances.balances[sh.sharder.id])
	//		assert.Zero(t, balances.balances[sh.delegate.id])
	//		for _, st := range sh.stakers {
	//			expected[st.id] += 0
	//			got[st.id] = balances.balances[st.id]
	//		}
	//	}
	//
	//	for _, sh := range filterClientsById(sharders, balances.blockSharders) {
	//		for _, st := range sh.stakers {
	//			expected[st.id] += 21e7 + 3e10 // + block sharders fees
	//		}
	//	}
	//
	//	assert.Equal(t, len(expected), len(got), "sizes of balance maps")
	//	assert.Equal(t, expected, got, "balances")
	//})

	// don't set DKG miners list, because no VC is expected

	// reset all balances
	balances.balances = make(map[string]state.Balance)

	//t.Run("pay fees -> view change interests", func(t *testing.T) {
	//	setRounds(t, msc, 500, 501, balances)
	//
	//	var generator, blck = prepareGeneratorAndBlock(miners, 0, 501)
	//
	//	// payFees transaction
	//	now += 10
	//	var tx = newTransaction(generator.miner.id, ADDRESS, 0, now)
	//	balances.txn = tx
	//	balances.block = blck
	//	balances.blockSharders = selectRandom(sharders, 3)
	//
	//	// add fees
	//	var gn, err = msc.getGlobalNode(balances)
	//	require.NoError(t, err, "getting global node")
	//
	//	_, err = msc.payFees(tx, nil, gn, balances)
	//	require.NoError(t, err, "pay_fees error")
	//
	//	// pools are active, rewards as above and +fees
	//
	//	var (
	//		expected = make(map[string]state.Balance)
	//		got      = make(map[string]state.Balance)
	//	)
	//
	//	for _, mn := range miners {
	//		assert.Zero(t, balances.balances[mn.miner.id])
	//		assert.Zero(t, balances.balances[mn.delegate.id])
	//		for _, st := range mn.stakers {
	//			if mn == generator {
	//				expected[st.id] += 77e7 + 1e10
	//			} else {
	//				expected[st.id] += 1e10
	//			}
	//			got[st.id] = balances.balances[st.id]
	//		}
	//	}
	//
	//	for _, sh := range sharders {
	//		assert.Zero(t, balances.balances[sh.sharder.id])
	//		assert.Zero(t, balances.balances[sh.delegate.id])
	//		for _, st := range sh.stakers {
	//			expected[st.id] += 1e10
	//			got[st.id] = balances.balances[st.id]
	//		}
	//	}
	//
	//	for _, sh := range filterClientsById(sharders, balances.blockSharders) {
	//		for _, st := range sh.stakers {
	//			expected[st.id] += 21e7
	//		}
	//	}
	//
	//	assert.Equal(t, len(expected), len(got), "sizes of balance maps")
	//	assert.Equal(t, expected, got, "balances")
	//})

	t.Run("epoch", func(t *testing.T) {
		var global, err = msc.getGlobalNode(balances)
		require.NoError(t, err)

		var interestRate, rewardRate = global.InterestRate, global.RewardRate
		global.epochDecline()

		assert.True(t, global.InterestRate < interestRate)
		assert.True(t, global.RewardRate < rewardRate)
	})

}

func prepareGeneratorAndBlock(miners []*TestClient, idx int, round int64) (
	generator *TestClient, blck *block.Block) {

	generator = miners[idx]

	blck = block.Provider().(*block.Block)
	blck.Round = round                                // VC round
	blck.MinerID = generator.client.id                // block generator
	blck.PrevBlock = block.Provider().(*block.Block)  // stub

	return generator, blck
}

func unwrapClients(clients []*TestClient) (list []*Client) {
	list = make([]*Client, 0, len(clients))
	for _, mn := range clients {
		list = append(list, mn.client)
	}
	return
}

func selectRandom(clients []*TestClient, n int) (selection []string) {
	if n > len(clients) {
		panic("too many elements requested")
	}

	selection = make([]string, 0, n)

	var permutations = rand.Perm(len(clients))
	for i := 0; i < n; i++ {
		selection = append(selection, clients[permutations[i]].client.id)
	}
	return
}

func filterClientsById(clients []*TestClient, ids []string) (
	selection []*TestClient) {

	selection = make([]*TestClient, 0, len(ids))

	for _, client := range clients {
		for _, id := range ids {
			if client.client.id == id {
				selection = append(selection, client)
			}
		}
	}
	return
}
