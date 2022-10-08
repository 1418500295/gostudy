package cases

import (
	"gostudy/src/selenium/pageobject"
	"log"
	"testing"
)

func TestBet(t *testing.T) {
	log.Println("点击投注记录")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	bet := pageobject.WalletRecordPage{
		Base: basePage,
	}
	bet.ClickWalletRecord()
}
