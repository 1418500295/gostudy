# shellcheck disable=SC2164
cd ./cases/
go test -json | go-test-report -o ../test_report.html
