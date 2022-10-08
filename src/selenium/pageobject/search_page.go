package pageobject

import (
	"github.com/tebeka/selenium"
	"gostudy/src/selenium/basepage"
)

const (
	searchEleXpath = "//*[@id=\"root\"]/div/div[1]/div[1]/label/input"
)

const (
	expResult = "//*[@id=\"root\"]/div/div[2]/div[1]/div[1]/div/div[2]/div[1]/div[2]/div/p"
)

type SearchPage struct {
	Base basepage.Base
}

func (search *SearchPage) ClickSearch() {
	err := search.Base.FindElementByXpath(searchEleXpath).Click()
	if err != nil {
		return
	}
}

func (search *SearchPage) InputContent(content string) {
	err := search.Base.FindElementByXpath(searchEleXpath).SendKeys(content)
	if err != nil {
		return
	}
}

func (search *SearchPage) FindExpResult() selenium.WebElement {
	ele := search.Base.FindElementByXpath(expResult)
	return ele

}
