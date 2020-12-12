## shuax-chrome-upgrade

> 一个可以升级 [shuax chrome 便携版](https://shuax.com/project/chrome/)  的工具

> [github](https://github.com/hezhizheng/shuax-chrome-upgrade)

### 功能
- 简单交互式操作
- 自动检测最新的 `shuax chrome` 版本
- 用户决定是否进行升级操作(自动下载、解压、重命名文件等)
- windows下 配合.bat文件 实现开机自动检测更新功能

### 流程

### 使用
自定义config.json配置文件(shuax chrome 的安装目录)

例：假如我的shuax chrome 安装解压目录为

![free-pic](https://vkceyugu.cdn.bspapp.com/VKCEYUGU-imgbed/b13f6f47-c970-4814-884a-8b30342a5808.png
)

那么 local_chrome_path 就定义为 `E:\\chrome`。如下：
```
# 参数说明
{
  "comments": {
    "comment1": "// 这些是注释不用理会，local_chrome_path：本地chrome安装路径"
  },
  "app": {
    "local_chrome_path": "E:\\chrome"
  }
}
```

编译 (windows提供编译好的文件 shuax-chrome-upgrade.7z
下载 [releases](https://github.com/hezhizheng/shuax-chrome-upgrade/releases) )

![free-pic](https://vkceyugu.cdn.bspapp.com/VKCEYUGU-imgbed/58d3ddf7-5060-457d-a996-f9fa5a4cefd5.png
)

自动编译
```
go build
```

运行
- 请不要随意更改`shuax chrome`原本的目录结构
- 保证编译的文件与 config.json、7z.dll、7z.exe 文件 在同级目录
- 执行 ./shuax-chrome-upgrade.exe 或者双击启动，根据提示输入指令完成升级

升级

![free-pic](https://vkceyugu.cdn.bspapp.com/VKCEYUGU-imgbed/b771db6b-022f-48cd-9aec-f93b159fc2a7.png)

无需升级

![free-pic](https://vkceyugu.cdn.bspapp.com/VKCEYUGU-imgbed/783179c0-d2ff-4012-b2bc-d1808109434c.png)


windows 开机自动检测(创建.bat文件)

./shuax-chrome-upgrade.bat

创建快捷方式，设定开机自启即可