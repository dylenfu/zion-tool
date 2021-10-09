package flag

import (
	"github.com/dylenfu/zion-tool/pkg/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"math/big"
	"strings"
	"time"
)

//GetFlagName deal with short flag, and return the flag name whether flag name have short name
func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}

func Flag2string(ctx *cli.Context, f cli.Flag) string {
	fn := GetFlagName(f)
	data := ctx.String(fn)
	return data
}

func Flag2address(ctx *cli.Context, f cli.Flag) common.Address {
	data := Flag2string(ctx, f)
	return common.HexToAddress(data)
}

func Flag2big(ctx *cli.Context, f cli.Flag) *big.Int {
	fn := GetFlagName(f)
	data := ctx.String(fn)
	return math.String2BigInt(data)
}

func Flag2Uint64(ctx *cli.Context, f cli.Flag) uint64 {
	fn := GetFlagName(f)
	data := ctx.Uint64(fn)
	return data
}

func Flag2Duration(ctx *cli.Context, f cli.Flag) (time.Duration, error) {
	fn := GetFlagName(f)
	data := ctx.String(fn)

	tmp, err := time.ParseDuration(data)
	if err != nil {
		return 0, err
	}
	return time.Duration(tmp), nil
}
