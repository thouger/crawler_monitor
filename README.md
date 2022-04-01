# crawler_monitor

### 爬虫监控系统

#### 使用方法
1. python3 -m venv venv
    1. 创建venv这个虚拟环境文件夹,这一步操作会影响后面的进入虚拟环境
2. 填写两个配置文件:utils/config/spider_tasks.toml和utils/config/db.toml
    1. spider_tasks.toml下的爬虫任务文件,按照例子的格式来
    2. db.toml需要插入爬虫状态的数据库账号密码等等,暂时只支持mysql
3. 爬虫输出日志要有"实际上插入的数据 %s 条"关键字,用来统计最后爬虫爬取到的数据量(可选)
    1. 这是我的写法:"共导出 %s 条数据 到 %s,实际上插入的数据 %s 条,重复 %s 条" % (1,2,3,4)
4. 创建爬虫监控数据库crawler_done,里面要有:字段(待添加)
5. 运行下面的命令

##### 运行命令
1. go run daemon/daemon.go monitor 会启动监控程序,读取utils/config/spider_tasks.toml下的爬虫任务进行长期守护式的监控
2. go run daemon/daemon.go run name "test" 运行spider_tasks.toml下name为test的爬虫
3. go run daemon/daemon.go run id 1 运行spider_tasks.toml下id为1的爬虫
4. go run daemon/daemon.go run all 运行当天所有爬虫程序

#### 文件树介绍
1. daemon/daemon.go是运行文件
2. daemon/monitor/monitor.go是监控文件
3. daemon/child/process.go是子进程启动文件
4. daemon/cli/cli.go是存放命令行相关指令的文件

### 主要功能
1. 每个一个小时自动执行git pull命令
2. 启动长期运行的守护进程通过执行shell命令执行爬虫任务
3. [x] 收集日志统计保存在log文件夹下(已取消,用过爬虫框架代替了)
4. 判断爬虫运行途中是否出错还是成功执行完毕:
    1. 爬虫出错:将发送信息发送到企业微信(目前只支持企业微信)
    2. 无论爬虫功能是否,都会在爬虫启动时,以及结束时(出错或成功)往数据库插入消息表示爬虫启动以及完成或出错.如果是完成,会记录爬虫程序爬取过多少的数据量
5. 数据库连接是用连接池功能,最大连接数100,闲置时是20

### todo
1. 列出现在正在运行的爬虫任务有哪一些
2. 除了启动子进程以外的操作应该分离出去或者加个大异常跳过,又或者专门启动进程去操作这些步骤.
3. 爬虫出错发送通知人要手动去源码那里改
4. 有人用就去改
5. 待添加
