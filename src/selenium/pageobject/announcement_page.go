package pageobject

import (
	"github.com/tebeka/selenium"
	"gostudy/src/selenium/basepage"
)

type AnnouPage struct {
	Base basepage.Base
}

const (
	checkAnnouXpath = "//*[@id=\"root\"]/div/div[1]/div[3]/div/div/div/div[2]/div"
)

func (annou *AnnouPage) CheckAnnouIsExist() bool {
	ele := annou.Base.FindElementByXpath(checkAnnouXpath)
	childEles, _ := ele.FindElements(selenium.ByTagName, "div")
	if len(childEles) != 0 {
		return true
	} else {
		return false
	}
}

const (
	annouButton = "//*[@id=\"root\"]/div/div[1]/div[3]/div/div/div/div[1]"
)

func (annou *AnnouPage) ClickAnnouButton() {
	err := annou.Base.FindElementByXpath(annouButton).Click()
	if err != nil {
		return
	}
}

const (
	weihuAnnouXpath  = "//*[@id=\"root\"]/div/div[1]/div/ul/li[2]/span"
	pankouAnnouXpath = "//*[@id=\"root\"]/div/div[1]/div/ul/li[3]/span"
	annouDetailXpath = "//*[@id=\"root\"]/div/div[2]/div[2]/div"
)

const (
	expResXpath = "//*[@id=\"root\"]/div/div[2]/div[1]/div/div"
)

func (annou *AnnouPage) ClickWeiHuAnnou() {
	err := annou.Base.FindElementByXpath(pankouAnnouXpath).Click()
	if err != nil {
		return
	}
}

func (annou *AnnouPage) GetAnnouList() []selenium.WebElement {
	paraentEle := annou.Base.FindElementByXpath(annouDetailXpath)
	childList, _ := paraentEle.FindElements(selenium.ByTagName, "div")
	return childList
}

func (annou *AnnouPage) GetExpRes() selenium.WebElement {
	return annou.Base.FindElementByXpath(expResXpath)
}
