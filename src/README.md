测试报告：
安装：go install github.com/vakenbolt/go-test-report/
执行：go test -json | go-test-report -o ../test_report.html
-o : 指定在某个位置生成指定名称的报告

测试报告：
安装：go install github.com/Thatooine/go-test-html-report@latest
使用：go test -v -cover -json  ./... | go-test-html-report

打包为linux下可执行文件（mac）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main-go-linux apitest.go
打包为linux下可执行文件（windows）
设置：
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build

进入测试case所在目录，执行如下：
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -c
运行指定测试函数：
./cases.test -test.run ^Test_Two$ -test.v

解析proto文件
protoc --go_out=./ ./new_api.proto


gomod爆红：GOPROXY=https://goproxy.cn,direct