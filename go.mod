module github.com/dylenfu/zion-tool

go 1.15

require (
	github.com/btcsuite/goleveldb v1.0.0
	github.com/ethereum/go-ethereum v1.10.14
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.4
)

replace github.com/ethereum/go-ethereum v1.10.14 => ../Zion
