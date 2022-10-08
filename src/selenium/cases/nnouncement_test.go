package cases

import (
	"gostudy/src/selenium/pageobject"
	"gostudy/src/selenium/utils"
	"log"
	"testing"
	"time"
)

func TestCheckAnnouIsExist(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("异常：", err)
		}
	}()
	log.Println("检查跑马灯公告是否有数据...")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	annou := pageobject.AnnouPage{
		Base: BasePage,
	}
	res := annou.CheckAnnouIsExist()
	if res == true {
		log.Println("跑马灯公告有数据")
	}
	utils.AssertEqual(true, res, utils.RunFuncName())

}
func TestClickAnnouButton(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("异常：", err)
		}
	}()
	log.Println("点击盘口公告...")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	annou := pageobject.AnnouPage{
		Base: BasePage,
	}
	annou.ClickAnnouButton()
	Handles, _ = Driver.WindowHandles()
	err1 := Driver.SwitchWindow(Handles[2])
	if err1 != nil {
		return
	}
	annou.ClickWeiHuAnnou()
	annouList := annou.GetAnnouList()
	log.Printf("存在%v条公告", len(annouList))
	expR, _ := annou.GetExpRes().Text()
	utils.AssertEqual(expR, "盘口公告", utils.RunFuncName())
	time.Sleep(5 * time.Second)
}
