package pageobject

import "gostudy/src/selenium/basepage"

const (
	accountByXpath   string = "body > form:nth-child(4) > input:nth-child(1)"
	operatorById     string = "loginOperator"
	loginSiteLocById string = "loginSite"
	submitByXpath    string = "body > form:nth-child(4) > input[type=submit]:nth-child(5)"
)

type LoginPage struct {
	Base basepage.Base
}

func (login *LoginPage) ClearContext() {
	err := login.Base.FindElementByCssSelector(accountByXpath).Clear()
	if err != nil {
		return
	}
}
func (login *LoginPage) InputUserName(username string) {
	err := login.Base.FindElementByCssSelector(accountByXpath).SendKeys(username)
	if err != nil {
		return
	}
}

func (login *LoginPage) InputPassWord(pwd string) {
	err := login.Base.FindElementByID(operatorById).SendKeys(pwd)
	if err != nil {
		return
	}
}

func (login *LoginPage) InputSite(site string) {
	err := login.Base.FindElementByID(loginSiteLocById).SendKeys(site)
	if err != nil {
		return
	}
}

func (login *LoginPage) ClickLogin() {
	err := login.Base.FindElementByCssSelector(submitByXpath).Click()
	if err != nil {
		return
	}
}
