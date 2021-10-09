package core

import (
	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/sdk"
)

func masterAccount() (*sdk.Account, error) {
	chainID := config.Conf.ChainID
	node := config.Conf.Nodes[0]
	return sdk.CustomNewAccount(chainID, node.Url, node.PrivateKey)
}
