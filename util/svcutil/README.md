# svcutil 服务基础工具

该目录主要是用来做基础 CMD 的组件来使用

## 优雅退出和默认支持健康检查

支持退出:

* ctrl+c 退出,输出
* kill pid 输出
* kill -USR1 pid 输出
* kill -USR2 pid 输出

### StandBy

StandBy 会快速异步启动 HTTP 健康检查，和同步执行任务，**任务结束就会退出**，或者收到退出信号

```
package main

import (
	"fmt"
	"time"

	"git.llsapp.com/zhenghe/pkg/util/svcutil"
)

func main() {
	svcutil.StandBy(":8080", func() {
		fmt.Println("doing ...")
		time.Sleep(5 * time.Second)
		fmt.Println("done ...")
	})
}
```


### NeverStop

NeverStop 会快速异步启动 HTTP 健康检查，和同步执行任务，**执行结束也不会退出** (除非收到退出信号)

* block run
* health check
* never stop
* support exit when received signal

```
// NeverStop graceful do func and with HTTP health check at addr
// it will block and never stop
func NeverStop(addr string, f func())
```

demo:

```
package main

import (
	"fmt"
	"time"

	"git.llsapp.com/zhenghe/pkg/util/svcutil"
)

func main() {
	tick := time.NewTicker(60* time.Second)

	svcutil.GracefulRun(":8080", func() {
		for {
			select {
			case <-tick.C:
				fmt.Printf("hello at %s\n", time.Now())
			default:
			}
		}
	})

	tick.Stop()
}
```

