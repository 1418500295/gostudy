package pageobject

import (
	"github.com/tebeka/selenium"
	"gostudy/src/selenium/basepage"
)

type BetRecordPage struct {
	Base basepage.Base
}

const (
	betRecordByXpath = "//*[@id=\"root\"]/div/nav/div[1]/button[2]/div"
)

const (
	expBetRecordPageEle = "//*[@id=\"root\"]/div/div[2]/section/div[2]/div[2]/div[1]/div/p"
)

func (bet *BetRecordPage) ClickBetRecord() {
	err := bet.Base.FindElementByXpath(betRecordByXpath).Click()
	if err != nil {
		return
	}
}

func (bet *BetRecordPage) FindExpBetRecordPageEle() selenium.WebElement {
	return bet.Base.FindElementByXpath(expBetRecordPageEle)
}
