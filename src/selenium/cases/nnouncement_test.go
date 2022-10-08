package cases

import (
	"gostudy/src/selenium/pageobject"
	"log"
	"testing"
)

func TestClickAnnouButton(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("异常：", err)
		}
	}()
	log.Println("点击公告按钮")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	annou := pageobject.AnnouPage{
		Base: BasePage,
	}
	annou.ClickAnnouButton()
}
