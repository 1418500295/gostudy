package basepage

import "github.com/tebeka/selenium"

type Base struct {
	Driver selenium.WebDriver
}

func (base *Base) driver() selenium.WebDriver {
	return base.Driver
}

func (base *Base) FindElementByXpath(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByXPATH, element)
	return ele
}
func (base *Base) FindElementByName(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByName, element)
	return ele
}
func (base *Base) FindElementByClassName(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByClassName, element)
	return ele
}
func (base *Base) FindElementByID(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByID, element)
	return ele
}
func (base *Base) FindElementByCssSelector(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByCSSSelector, element)
	return ele
}

func (base *Base) FindElementByLinkText(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByLinkText, element)
	return ele
}

func (base *Base) FindElementByTagName(element string) selenium.WebElement {
	ele, _ := base.Driver.FindElement(selenium.ByTagName, element)
	return ele
}

func (base *Base) FindElementsByName(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByName, element)
	return ele
}
func (base *Base) FindElementsByClassName(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByClassName, element)
	return ele
}
func (base *Base) FindElementsByID(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByID, element)
	return ele
}
func (base *Base) FindElementsByCssSelector(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByCSSSelector, element)
	return ele
}
func (base *Base) FindElementsByXpath(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByXPATH, element)
	return ele
}
func (base *Base) FindElementsByLinkText(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByLinkText, element)
	return ele
}

func (base *Base) FindElementsByTagName(element string) []selenium.WebElement {
	ele, _ := base.Driver.FindElements(selenium.ByTagName, element)
	return ele
}
