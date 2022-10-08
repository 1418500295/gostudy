package pageobject

import (
	"github.com/tebeka/selenium"
	"gostudy/src/selenium/basepage"
)

//早盘投注
const (
	tyBetXpath          = "//*[@id=\"root\"]/div/nav/div[1]/button[1]/div"
	zaoPanClickXpath    = "//*[@id=\"root\"]/div/div[2]/div[1]/section/div[2]/button[3]/div"
	matchInfoByClass    = "flex flex-col bg-framework100 shadow-sm mb-underline"
	zaoPanBetByXpath    = "//*[@id=\"root\"]/div/div[2]/section/div[2]/div[1]/div[2]/div/div[2]/div/div[1]/div[2]/div[1]/div/div/div[1]/div/div"
	inputBetAmountXpath = "//*[@id=\"betslip\"]/div[2]/div[1]/div/div[2]/label/input"
	subMitBetXpath      = "//*[@id=\"betslip\"]/div[2]/div[2]/div/div[2]/button"
)
const (
	ExpTYBetPageEle = "//*[@id=\"betslip\"]/div/div[1]/p"
	confirmXpath    = "//*[@id=\"betslip\"]/div/div[3]/div/button[1]"
)

type TYBetPage struct {
	Base basepage.Base
}

func (zaoPan *TYBetPage) ClickTYBet() {
	err := zaoPan.Base.FindElementByXpath(tyBetXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) ClickZaoPan() {
	err := zaoPan.Base.FindElementByXpath(zaoPanClickXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) ClickZaoPanBet() {
	err := zaoPan.Base.FindElementByXpath(zaoPanBetByXpath).Click()
	if err != nil {
		return
	}
}

const (
	paraentXpath    = "//*[@id=\"root\"]/div/div[2]/section/div[2]"
	secondClassName = "flex flex-col bg-framework100 shadow-sm mb-underline"
	thirdXpath      = "flex justify-around items-center p-sm w-full border-b min-h-[36px] bg-framework100 font-bold clickable border-b-0.5"

	touzhudanXpath       = "//*[@id=\"betslip\"]/div[1]/div[1]"
	chuanguanXpath       = "//*[@id=\"betslip\"]/div[2]/div[1]/button[2]"
	inputAmountClassName = "//*[@id=\"betslip\"]/div[2]/div[2]"

	submitMoreBetXpath = "//*[@id=\"betslip\"]/div[2]/div[3]/div/div[2]/button"
)

func (zaoPan *TYBetPage) ClickZaoPanMoreBet() {
	paraEle := zaoPan.Base.FindElementByXpath(paraentXpath)
	eles, _ := paraEle.FindElements(selenium.ByCSSSelector, "[class=\"flex flex-col bg-framework100 shadow-sm mb-underline\"]")
	//3串
	for _, one := range eles[:3] {
		childEles, _ := one.FindElements(selenium.ByCSSSelector, "[class=\"flex justify-center items-center relative text-subtitle\"]")
		for _, value := range childEles {
			valueText, _ := value.Text()
			if valueText != "" {
				err := value.Click()
				if err != nil {
					return
				}
				break
			}
		}

	}
}

func (zaoPan *TYBetPage) ClickTouZhuDan() {
	err := zaoPan.Base.FindElementByXpath(touzhudanXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) ClickChuanGuan() {
	err := zaoPan.Base.FindElementByXpath(chuanguanXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) InputMoreBetAmount(amount string) {
	ele := zaoPan.Base.FindElementByXpath(inputAmountClassName)
	eles, _ := ele.FindElements(selenium.ByCSSSelector, "[placeholder=\"输入金额\"]")
	for _, value := range eles {
		err := value.Clear()
		if err != nil {
			return
		}
		err1 := value.SendKeys(amount)
		if err1 != nil {
			return
		}
	}

}

func (zaoPan *TYBetPage) ClickSubmit() {
	err := zaoPan.Base.FindElementByXpath(submitMoreBetXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) InputZaoPanBetAmount(amount string) {
	err := zaoPan.Base.FindElementByXpath(inputBetAmountXpath).SendKeys(amount)
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) SubmitZaoPanBet() {
	err := zaoPan.Base.FindElementByXpath(subMitBetXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) GetZaoPanExpRes() selenium.WebElement {
	return zaoPan.Base.FindElementByXpath(ExpTYBetPageEle)
}

func (zaoPan *TYBetPage) ClickConfirm() {
	err := zaoPan.Base.FindElementByXpath(confirmXpath).Click()
	if err != nil {
		return
	}
}

func (zaoPan *TYBetPage) ClearInputBox() {
	err := zaoPan.Base.FindElementByXpath(inputBetAmountXpath).Clear()
	if err != nil {
		return
	}
}

//滚球投注
const (
	gunQiuClickXpath  = "//*[@id=\"root\"]/div/div[2]/div[1]/section/div[2]/button[1]/div"
	gunQiuBetXpath    = "//*[@id=\"root\"]/div/div[2]/section/div[2]/div[1]/div[2]/div/div[2]/div/div[1]/div[2]/div[1]/div/div/div[1]/div"
	gunQIuAmountXpath = "//*[@id=\"betslip\"]/div[2]/div[1]/div/div[2]/label/input"
	gunQiuSubmitXpath = "//*[@id=\"betslip\"]/div[2]/div[2]/div/div[2]/button"
)

const (
	gunQiuExpResultXpath = "//*[@id=\"betslip\"]/div/div[1]/p"
)

func (gunQiu *TYBetPage) ClickGunQiu() {
	err := gunQiu.Base.FindElementByXpath(gunQiuClickXpath).Click()
	if err != nil {
		return
	}
}
func (gunQiu *TYBetPage) ClickGunQiuBet() {
	err := gunQiu.Base.FindElementByXpath(gunQiuBetXpath).Click()
	if err != nil {
		return
	}
}
func (gunQiu *TYBetPage) InputGunQiuBetAmount(amount string) {
	err := gunQiu.Base.FindElementByXpath(gunQIuAmountXpath).SendKeys(amount)
	if err != nil {
		return
	}
}
func (gunQiu *TYBetPage) SubmitGunQiuBet() {
	err := gunQiu.Base.FindElementByXpath(gunQiuSubmitXpath).Click()
	if err != nil {
		return
	}
}

func (gunQiu *TYBetPage) GetGunQiuExpRes() selenium.WebElement {
	return gunQiu.Base.FindElementByXpath(gunQiuExpResultXpath)
}

//今日投注
const (
	jinRiClickXpath  = "//*[@id=\"root\"]/div/div[2]/div[1]/section/div[2]/button[2]/div"
	jinRiBetXpath    = "//*[@id=\"root\"]/div/div[2]/section/div[2]/div[2]/div[1]/div[2]/div/div[2]/div/div[1]/div[2]/div[1]/div/div/div[1]/div"
	jinRiAmountXpath = "//*[@id=\"betslip\"]/div[2]/div[1]/div/div[2]/label/input"
	jinRiSubmitXpath = "//*[@id=\"betslip\"]/div[2]/div[2]/div/div[2]/button"
)
const (
	jinRiExpResXpath = "//*[@id=\"betslip\"]/div/div[1]/p"
)

func (jinRi *TYBetPage) ClickJinRi() {
	err := jinRi.Base.FindElementByXpath(jinRiClickXpath).Click()
	if err != nil {
		return
	}
}

func (jinRi *TYBetPage) ClickJinRiBet() {
	err := jinRi.Base.FindElementByXpath(jinRiBetXpath).Click()
	if err != nil {
		return
	}
}
func (jinRi *TYBetPage) InputJinRiBetXpath(amount string) {
	err := jinRi.Base.FindElementByXpath(jinRiAmountXpath).SendKeys(amount)
	if err != nil {
		return
	}
}
func (jinRi *TYBetPage) SubmitJinRiBet() {
	err := jinRi.Base.FindElementByXpath(jinRiSubmitXpath).Click()
	if err != nil {
		return
	}
}

func (jinRi *TYBetPage) GetJinRiExpRes() selenium.WebElement {
	return jinRi.Base.FindElementByXpath(jinRiExpResXpath)
}
