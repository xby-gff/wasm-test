# wasm-test
go version 1.22.1
1、设置变量:
set GOOS=js 
set GOARCH=wasm
2、构建wasm :
go build -o .\static\main.wasm .\assets\main.go
3、启动服务：
go run .\server\server.go
4、访问页面：
http://localhost:8080
