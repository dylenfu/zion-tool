/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"time"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
)

// propose new epoch and vote until success
func Epoch() bool {
	var param struct {
		ProposerNodeIndex        int
		OldParticipantsIndexList []int
		VoterIndexList           []int
		NewParticipantsIndexList []int
		NextEpochStartHeight     uint64
	}

	log.Split("start to change epoch")

	if err := config.LoadParams("test_epoch.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	proposer, err := generateAccount(param.ProposerNodeIndex)
	if err != nil {
		log.Errorf("failed to generate proposer, err: %v", err)
		return false
	}

	log.Info("start to propose...")
	newParticipantPeers, err := getPeers(param.NewParticipantsIndexList)
	if err != nil {
		log.Errorf("failed to generate proposal peers, err: %v", err)
		return false
	}
	tx, err := proposer.Propose(param.NextEpochStartHeight, newParticipantPeers)
	if err != nil {
		log.Errorf("failed to propose, err: %v", err)
		return false
	}
	log.Split("validator %s propose, hash %s", proposer.Address().Hex(), tx.Hex())

	epoch, err := getProposalReceipt(proposer, tx)
	if err != nil {
		log.Errorf("failed to get proposal receipt, err: %v", err)
		return false
	}
	log.Split("epoch as follow:\r\n%s", epoch.String())

	log.Info("start to vote...")
	voters, err := generateAccounts(param.VoterIndexList)
	if err != nil {
		log.Errorf("failed to generate voter account, err: %v", err)
		return false
	}

	epochHash := epoch.Hash()
	curEpoch, err := voters[0].Epoch()
	if err != nil {
		log.Errorf("failed to get current epoch, err: %v", err)
		return false
	}

	epochID := curEpoch.ID + 1
	for _, voter := range voters {
		if tx, err := voter.Vote(epochID, epochHash); err != nil {
			log.Errorf("voter %s vote failed", voter.Address().Hex())
		} else {
			log.Infof("voter %s voted, hash %s", voter.Address().Hex(), tx.Hex())
		}
		time.Sleep(5 * time.Second)
	}

	return true
}
