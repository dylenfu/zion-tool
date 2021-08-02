package journal

import (
	"fmt"

	"github.com/urfave/cli"
)

// HandleNative try native governance call `epoch`
func HandleNative(ctx *cli.Context) error {
	fmt.Println("start to call epoch", "start", true)

	// load config instance
	c, err := getConfig(ctx)
	if err != nil {
		return err
	}

	// generate master account
	fmt.Println("try to generate master account...")
	master, err := generateMasterAccount(c)
	if err != nil {
		return err
	}

	epoch, err := master.Epoch()
	if err != nil {
		return err
	}
	fmt.Printf("epoch is %d\r\n", epoch)
	return nil
}
