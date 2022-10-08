package cases

import (
	"fmt"
	"github.com/tebeka/selenium"
	"gostudy/src/selenium/basepage"
	"gostudy/src/selenium/pageobject"
	"gostudy/src/selenium/utils"
	"log"
	"testing"
	"time"
)

const (
	chromePath  = "/root/test/selenium_test/chromedriver"
	pathWindows = "/Users/eden/go/src/gostudy/src/selenium/chromedriver"
	port        = 1111
	reqUrl      = "https://bf5376facdb172cb79.html"
)

var (
	Driver   selenium.WebDriver
	BasePage basepage.Base
	Handles  []string
)

func TestMain(m *testing.M) {
	var err error
	opts := []selenium.ServiceOption{
		//selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	//selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(pathWindows, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	//defer service.Stop()
	// set browser as chrome
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	//chromeCaps := chrome.Capabilities{
	//	Args: []string{
	//		"--headless",
	//		"--start-maximized",
	//		"--window-size=1200x600",
	//		"--no-sandbox",
	//		//"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	//		//"--disable-gpu",
	//		//"--disable-impl-side-painting",
	//		//"--disable-gpu-sandbox",
	//		//"--disable-accelerated-2d-canvas",
	//		//"--disable-accelerated-jpeg-decoding",
	//		//"--test-type=ui",
	//
	//	},
	//}
	//caps.AddChrome(chromeCaps)

	// remote to selenium server
	if Driver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port)); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	BasePage = basepage.Base{Driver: Driver}
	err = Driver.Get(reqUrl)

	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}
	err1 := Driver.SetImplicitWaitTimeout(10 * time.Second)
	if err1 != nil {
		return
	}
	err2 := Driver.MaximizeWindow("")
	if err2 != nil {
		return
	}
	log.Println("开始登陆...")
	login := pageobject.LoginPage{Base: BasePage}
	login.ClearContext()
	login.InputUserName("eden888")
	login.ClickLogin()
	time.Sleep(5 * time.Second)
	Handles, _ = Driver.WindowHandles()
	m.Run()
	_ = service.Stop()
	_ = Driver.Quit()
	fmt.Printf("\033[33m====== 总成功数: %v  总失败数：%v ======\033[0m\n\n", utils.OkNum, utils.FailNum)
	if len(utils.FailCaseList) != 0 {
		fmt.Println("\033[31m失败的用例:\033[0m")
		for _, caseName := range utils.FailCaseList {
			fmt.Printf("\033[31m...funcName: %v...\033[0m\n", caseName)
		}
	}

}
