package bitcoinbalance

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"testing"
	"time"
)

func TestBtcBalanceModule(t *testing.T) {
	InitBitcoinBalance("https://sequencer-be.regtest.trustless.computer")
	module := GetBitcoinBalanceModule()

	addr, _ := common.NewMixedcaseAddressFromString("0xB8Bc42098CC278e03302DbcF9138Cf04b22e6c04")
	fmt.Println(module.GetBalance([]common.Address{addr.Address()}))
	err := module.DeductBalance([]common.Address{addr.Address()}, []string{"100"}, []string{time.Now().String()}, "batch1111")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 1000)
	fmt.Println(module.GetBalance([]common.Address{addr.Address()}))

}
