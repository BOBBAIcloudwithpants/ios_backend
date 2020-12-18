## 使用说明
可以参考：
https://razeencheng.com/post/go-swagger
https://github.com/swaggo/swag
如果懒得看就直接看controllers/user.go中的示例注解
swag工具安装：`go get -u github.com/swaggo/swag/cmd/swag`
（每次`run`之前运行`swag init`更新文档）
在`run`之后可以通过 http://localhost:5000/swagger/index.html 访问文档