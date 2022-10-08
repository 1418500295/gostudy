package pageobject

import (
	"gostudy/src/selenium/basepage"
)

type BetRecordPage struct {
	Base basepage.Base
}

const (
	betRecordByXpath = "//*[@id=\"root\"]/div/nav/div[1]/button[2]/div"
)

func (bet *BetRecordPage) ClickBetRecord() {
	err := bet.Base.FindElementByXpath(betRecordByXpath).Click()
	if err != nil {
		return
	}
}
