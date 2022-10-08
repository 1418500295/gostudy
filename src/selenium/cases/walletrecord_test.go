package cases

import (
	"fmt"
	"gostudy/src/selenium/pageobject"
	"gostudy/src/selenium/utils"
	"log"
	"testing"
)

//钱包记录
func TestWalletRecord(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("钱包记录异常：", err)
		}
	}()
	log.Println("点击钱包记录...")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	bet := pageobject.WalletRecordPage{
		Base: BasePage,
	}
	bet.ClickWalletRecord()
	ele := bet.FindExpWalletRecordPageEle()
	expText, _ := ele.Text()
	utils.AssertEqual("投注总额:", expText, utils.RunFuncName())
	//assert.Equal(t, "投注总额:", expText)
}
