# bbs_backend

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
 docker build -t bbs_backend .
```

- 运行镜像
这里提供了以前台模式运行的范例，您可以根据需求改成后台模式运行
```
 docker run -it -p 5000:5000 bbs_backend
```

## 功能
- 用户能够自由浏览公开的`forum`中的内容
- 用户能够创建自己的公开/私有forum
    - 对于private forum, 创建者能够邀请其他用户加入，或者踢出现有用户，不可被关注
    - 对于public forum, 所有用户都可见，可被关注
- `forum` 包含两个板块，
    - `post` 板块：包含若干帖子，每个帖子都有`title`, `content` 以及一些关联的文件；每个帖子都可以被评论
    - `hole` 板块：类似于树洞，hole中的帖子都是匿名的，都只有文字信息，不可被评论


## 数据库设计
目前设计了8张表：
- user：用户表，存放用户名，密码，邮箱，头像等等，`is_admin` 字段表明是否是系统管理员
- forum: 论坛表，存放论坛的相关信息
- post：forum 下的两个板块之一，存放了若干帖子的title，content，like(点赞数)
- hole：forum 下的另一个板块，存放了匿名帖子的信息
- comment：评论表，存放了对于post的评论
- file：文件表，存放了文件id，网络地址和其所属的post（比如post中包含的图片会被放在file中）
- star：内容收藏表，是user与post的中间表（收藏是一个多对多关系，所以需要用一张中间表来表示）。
- forum_user：论坛成员表，是forum和user的中间表。如果forum为public，则表项的含义是用户关注了这个forum；如果forum为private，则表项的含义是用户是该forum的成员，role字段表示其在forum中的角色(user/admin/owner)
