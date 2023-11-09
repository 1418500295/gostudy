打包exe带图标：
1.go get(install) github.com/akavel/rsrc
2.1、创建manifest文件
仍然以test.go举例，假设打包文件名test.exe，ico文件名test.ico
manifest文件名test.exe.manifest
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
    <assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="x86"
    name="controls"
    type="win32">
    </assemblyIdentity>

    <dependency>
        <dependentAssembly>
            <assemblyIdentity
            type="win32"
            name="Microsoft.Windows.Common-Controls"
            version="6.0.0.0"
            processorArchitecture="*"
            publicKeyToken="6595b64144ccf1df"
            language="*">
            </assemblyIdentity>
        </dependentAssembly>
    </dependency>
</assembly>

2.2 执行rsrc -manifest test.exe.manifest -ico test.ico -o test.exe.syso
3. go build -0 指定文件名如：a.exe


func main() {
	cy := cyclicbarrier.New(num)
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	wg.Add(num)
	wg1.Add(1)
	for i := 0; i < num; i++ {
		go func(i int) {
			wg1.Wait()
			fmt.Println("准备开始。。。")
			time.Sleep(3 * time.Second)
			//fmt.Println(time.Now().UnixMilli())
			fmt.Printf("第%v个协程启动了\n", i)
			err := cy.Await(context.Background())
			if err != nil {
				return
			}
			wg.Done()
			fmt.Println("我在吃饭")
		}(i)
	}
	fmt.Println("开始")
	wg1.Done()
	wg.Wait()
	fmt.Println("结束")
}
