package pageobject

import (
	"github.com/tebeka/selenium"
	"gostudy/src/selenium/basepage"
)

const (
	walletRecordByXpath = "//*[@id=\"root\"]/div/nav/div[1]/button[3]/div"
)

const (
	expWalletRecordPageEle = "//*[@id=\"root\"]/div/div[2]/section/div[2]/div[1]/div/div[1]/p[1]"
)

type WalletRecordPage struct {
	Base basepage.Base
}

func (walletRecord *WalletRecordPage) ClickWalletRecord() {
	err := walletRecord.Base.FindElementByXpath(walletRecordByXpath).Click()
	if err != nil {
		return
	}
}

func (walletRecord *WalletRecordPage) FindExpWalletRecordPageEle() selenium.WebElement {
	return walletRecord.Base.FindElementByXpath(expWalletRecordPageEle)
}
