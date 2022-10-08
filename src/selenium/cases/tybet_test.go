package cases

import (
	"fmt"
	"gostudy/src/selenium/pageobject"
	"gostudy/src/selenium/utils"
	"log"
	"testing"
	"time"
)

//进行早盘投注
func TestZaoPanSimpleBet(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("早盘投注异常：", err)
		}
	}()
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	log.Println("早盘单关投注...")
	tyBet := pageobject.TYBetPage{Base: BasePage}
	tyBet.ClickTYBet()
	tyBet.ClickZaoPan()
	time.Sleep(1000 * time.Millisecond)
	tyBet.ClickZaoPanBet()
	time.Sleep(1000 * time.Millisecond)
	tyBet.ClearInputBox()
	tyBet.InputZaoPanBetAmount("6")
	tyBet.SubmitZaoPanBet()
	//var c selenium.Condition = func(wd selenium.WebDriver) (bool, error) {
	//	element, _ := wd.FindElement(selenium.ByXPATH, pageobject.ExpTYBetPageEle)
	//	return element.IsDisplayed()
	//}
	//_ = Driver.WaitWithTimeoutAndInterval(c, 5000, 1000)
	time.Sleep(5 * time.Second)
	ele := tyBet.GetZaoPanExpRes()
	expText, _ := ele.Text()
	log.Println("expText:", expText)
	utils.AssertEqual(expText, "投注成功", utils.RunFuncName())
	//assert.Equal(t, expText, "投注成功")
	//关闭下注弹窗
	tyBet.ClickConfirm()

}

//早盘串关投注
func TestZaoPanMoreBet(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("早盘串关异常：", err)
		}
	}()
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	log.Println("早盘串关投注...")
	tyBet := pageobject.TYBetPage{Base: BasePage}
	tyBet.ClickTYBet()
	tyBet.ClickZaoPan()
	time.Sleep(1000 * time.Millisecond)
	tyBet.ClickZaoPanMoreBet()
	time.Sleep(1000 * time.Millisecond)
	tyBet.ClickTouZhuDan()
	tyBet.ClickChuanGuan()
	tyBet.InputMoreBetAmount("7")
	tyBet.ClickSubmit()
	//var c selenium.Condition = func(wd selenium.WebDriver) (bool, error) {
	//	element, _ := wd.FindElement(selenium.ByXPATH, pageobject.ExpTYBetPageEle)
	//	return element.IsDisplayed()
	//}
	//_ = Driver.WaitWithTimeoutAndInterval(c, 5000, 1000)
	time.Sleep(5 * time.Second)
	ele := tyBet.GetZaoPanExpRes()
	expText, _ := ele.Text()
	utils.AssertEqual(expText, "投注成功", utils.RunFuncName())
	//关闭下注弹窗
	tyBet.ClickConfirm()

}

//滚球投注
func TestGunQiuSimpleBet(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("滚球投注异常：", err)
		}
	}()
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	log.Println("滚球单关投注...")
	tyBet := pageobject.TYBetPage{Base: BasePage}
	tyBet.ClickGunQiu()
	tyBet.ClickGunQiuBet()
	time.Sleep(1000 * time.Millisecond)
	tyBet.ClearInputBox()
	tyBet.InputGunQiuBetAmount("5")
	tyBet.SubmitGunQiuBet()
	time.Sleep(5 * time.Second)
	ele := tyBet.GetGunQiuExpRes()
	expText, _ := ele.Text()
	utils.AssertEqual(expText, "投注成功", utils.RunFuncName())
	//assert.Equal(t, expText, "投注成功")
	//关闭下注弹窗
	tyBet.ClickConfirm()

}

//今日投注
func TestJinRiSimpleBet(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("今日投注异常：", err)
		}
	}()
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	log.Println("今日单关投注...")
	tyBet := pageobject.TYBetPage{Base: BasePage}
	tyBet.ClickJinRi()
	tyBet.ClickJinRiBet()
	time.Sleep(1000 * time.Millisecond)
	tyBet.ClearInputBox()
	tyBet.InputJinRiBetXpath("5")
	tyBet.SubmitJinRiBet()
	time.Sleep(5 * time.Second)
	ele := tyBet.GetJinRiExpRes()
	expText, _ := ele.Text()
	utils.AssertEqual(expText, "投注成功", utils.RunFuncName())
	//assert.Equal(t, expText, "投注成功")
	//关闭下注弹窗
	tyBet.ClickConfirm()
}

//点击投注记录
func TestBetRecord(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("投注记录异常：", err)
		}
	}()
	log.Println("点击投注记录...")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	betRecord := pageobject.BetRecordPage{
		Base: BasePage,
	}
	betRecord.ClickBetRecord()
	if err != nil {
		return
	}
	expEle := betRecord.FindExpBetRecordPageEle()
	actualText, _ := expEle.Text()
	//assert.Contains(t, actualText, "单注")
	utils.AssertContains("单注", actualText, utils.RunFuncName())
	//assert.Contains(t, actualText, "单注")
}
