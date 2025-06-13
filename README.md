# 火烧云消息推送服务

每天都去看自己的城市是否有火烧云

太麻烦了！

于是我写了个程序，定时去获取火烧云指数，配置简单

如果火烧云满足你的阈值，直接将消息推送至你的微信

![](https://image.fengfengzhidao.com/rj_10259001f2cfba6630b318de3b4d39ee065.png)


程序使用go语言编写，项目运行很方便，支持docker部署

## 如何运行此项目
如果你不会go语言

那么可以看看我的go零基础入门课程

https://www.bilibili.com/video/BV1B6MwzGEYc/

搭建go环境
https://www.fengfengzhidao.com/article/ZtYtBIsBg90FB71eC4QU

这些都做好之后，配置好settings.yaml

使用 `go run main.go` 即可运行项目

## 消息推送
使用server酱推送到微信

注册server酱账号： https://sct.ftqq.com/

注册成功之后获取 sendKey，填入配置文件

免费账号有五次推送机会，对我们来说足够了

## 配置文件
```yaml
hsy:
  checkAod: 0.5 # 校验指标
  city: 长沙
  wxDate: 0 0 16 * * * # 晚霞的通知时间
  zxDate: 0 0 20 * * * # 朝霞的通知时间
serverBot:
  enable: true
  sendKey: # server酱上面的sendKey
```
修改你所在的城市，修改晚霞和朝霞的推送时间


## 如何部署
如果你有自己的云服务器
那么你可以直接构建docker镜像，然后使用docker运行容器