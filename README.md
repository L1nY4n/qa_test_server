# qa_test_server

## 服务部分
```shell
# 启动当前文件夹下的服务(main.go为入口) 
go run .  # go mian.go

#构建
go build

```


## 前端
---
### 技术栈
- **[pnpm](https://pnpm.io/)** - Fast, disk space efficient package manager
- **[vite](https://vitejs.dev/)** - Next Generation Frontend Tooling


```shell
# 进入web文件夹
cd qa-test-web

# 安装依赖
pnpm i 

# 开发模式 
pnpm  dev

#  前端编译打包，编译完成后得到 `dist` 文件夹，复制到 templates目录下作为 gin 的模板
pnpm build

```


## git 相关操作

```
# 拉取远端的更新并且合并到本地

git fetch upstream

# 合并到本地分支，先提交本地的改动，再合并代码 
git merge upstream/main 

# 解决冲突后提交并推送
commit  push
```