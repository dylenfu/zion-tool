module github.com/dylenfu/zion-tool

go 1.15

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/btcsuite/goleveldb v1.0.0
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/ethereum/go-ethereum v1.10.14
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.4
)

replace github.com/ethereum/go-ethereum v1.10.14 => ../Zion
