package cases

import (
	"gostudy/src/selenium/pageobject"
	"gostudy/src/selenium/utils"
	"log"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("搜索异常：", err)
		}
	}()
	log.Println("搜索联赛/球队")
	err := Driver.SwitchWindow(Handles[1])
	if err != nil {
		return
	}
	search := pageobject.SearchPage{
		Base: BasePage,
	}
	search.ClickSearch()
	search.InputContent("足球")
	freezeWindow := "setTimeout(function(){debugger}, 3000)"
	time.Sleep(1 * time.Second)
	_, _ = Driver.ExecuteScript(freezeWindow, nil)
	resul := search.FindExpResult()
	content, _ := resul.Text()
	//assert.Contains(t, content, "足球")
	utils.AssertContains("足球", content, utils.RunFuncName())
	Driver.Refresh()

}
