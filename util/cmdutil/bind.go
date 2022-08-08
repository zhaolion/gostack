package cmdutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhaolion/gostack/util/log"
	"github.com/zhaolion/gostack/util/svcutil"
)

type CommandHandle func(ctx context.Context, args []string) error

type CommandOption func(*cobra.Command)

func Addr() string {
	if v := os.Getenv("CMD_HTTP_ADDR"); v != "" {
		return v
	}

	return ":8000"
}

// Add command to main command
func Add(mainCmd *cobra.Command, use string, handle func(cmd *cobra.Command, args []string), short string) {
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Run:   handle,
	}
	mainCmd.AddCommand(command)
}

// WaitFor binding sub command for running tasks until function executed or receives stop signals
// 绑定子命令，并且提供持续运行能力(内部函数结束就会终止，一般用于不断更新 异步任务系统 - K8S Deployment)
func WaitFor(mainCmd *cobra.Command, use string, f func(stop <-chan struct{}) error, short string, options ...CommandOption) {
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("[%s] task start", use)
			defer func() {
				log.Infof("[%s] task finished", use)
			}()

			// 持续运行，直到内部逻辑退出
			svcutil.WaitFor(Addr(), f)
		},
	}

	for _, opt := range options {
		opt(command)
	}

	mainCmd.AddCommand(command)
}

// Bind binding sub command for running tasks forever
// 绑定子命令，并且提供持续运行能力(结束了仍然持续运行，一般用于不断更新 K8S Job)
func Bind(mainCmd *cobra.Command, use string, handle CommandHandle, short string, options ...CommandOption) {
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			log.Infof("[%s] task start", use)
			defer func() {
				log.Infof("[%s] task finished", use)
			}()

			svcutil.NeverStop(Addr(), func() {
				startTime := time.Now()
				var notification string
				if err := handle(ctx, args); err != nil {
					notification = fmt.Sprintf(cmdErrTemplate, use, args, time.Now().Sub(startTime), err)
				} else {
					notification = fmt.Sprintf(cmdFinishedTemplate, use, args, time.Now().Sub(startTime))
				}

				log.Info(notification)
			})
		},
	}

	for _, opt := range options {
		opt(command)
	}

	mainCmd.AddCommand(command)
}

// Instant binding sub command for running tasks once
// # 绑定子命令，只能运行一次就停止
func Instant(mainCmd *cobra.Command, use string, handle CommandHandle, short string, options ...CommandOption) {
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			startTime := time.Now()
			var notification string
			err := handle(ctx, args)
			if err != nil {
				notification = fmt.Sprintf(cmdErrTemplate, use, args, time.Now().Sub(startTime), err)
			} else {
				notification = fmt.Sprintf(cmdFinishedTemplate, use, args, time.Now().Sub(startTime))
			}
			// 反馈报错信号
			if err != nil {
				log.WithError(err).Fatal(notification)
			} else {
				log.Info(notification)
			}
		},
	}

	for _, opt := range options {
		opt(command)
	}

	mainCmd.AddCommand(command)
}

var cmdFinishedTemplate = `
## cmd successfully finished
%s done
args %v
执行时间: %s`

var cmdErrTemplate = `
## cmd failed with error
[task] %s done
args %v
执行时间: %s
错误: %s`
