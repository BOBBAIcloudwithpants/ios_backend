# ios_backend

## 启动方法

### 方法1：直接运行
- 开启 `GOMODULE`
    - Linux 下方法为 `export GO111MODULE=on`

- 在根目录下执行:
```
go run main.go
```

### 方法2：通过docker
前提是您已经配置好了docker环境
- build 后端镜像
```
 docker build -t ios_backend .
```

- 运行镜像
这里提供了以前台模式运行的范例，您可以根据需求改成后台模式运行
```
 docker run -it -p 5000:5000 ios_backend
```

