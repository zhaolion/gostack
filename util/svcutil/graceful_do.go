package svcutil

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zhaolion/gostack/util/healthcheck"
	"github.com/zhaolion/gostack/util/log"
	"github.com/zhaolion/gostack/util/waitutil"
)

// StandBy graceful doUtilStop func and with HTTP health check at addr
// it will block and stop when function is finished
func StandBy(addr string, f func()) {
	stop := WaitSignals()

	grace := &GracefulDo{}
	done := grace.Do(addr, f)

	for {
		select {
		case <-done:
			return
		case <-stop:
			return
		}
	}
}

// NeverStop graceful doUtilStop func and with HTTP health check at addr
// it will block and never stop
func NeverStop(addr string, f func()) {
	// wait for signal
	stop := WaitSignals()
	grace := &GracefulDo{}
	grace.DoUtilStop(addr, stop, f)
}

type GracefulDo struct {
	once sync.Once
}

func (g *GracefulDo) Do(addr string, f func()) <-chan struct{} {
	stop := make(chan struct{})

	// start health at once
	go g.withHealthCheck(addr, stop)

	go func() {
		defer close(stop)

		func() {
			defer waitutil.HandleCrash()
			f()
		}()
	}()

	return stop
}

func (g *GracefulDo) DoUtilStop(addr string, stop <-chan struct{}, f func()) {
	// start health at once
	go g.withHealthCheck(addr, stop)

	select {
	case <-stop:
		return
	default:
	}

	func() {
		defer waitutil.HandleCrash()
		f()
	}()

	// NOTE: b/c there is no priority selection in golang
	// it is possible for this to race, meaning we could
	// trigger t.C and stopCh, and t.C select falls through.
	// In order to mitigate we re-check stopCh at the beginning
	// of every loop to prevent extra executions of f().
	<-stop
}

func (g *GracefulDo) withHealthCheck(addr string, stop <-chan struct{}) {
	g.once.Do(func() {
		HTTPHealthCheck(addr, stop)
	})
}

// WaitFor 简单版本(并不优雅) 处理退出, 需要相关处理函数 f 能够阻塞执行
func WaitFor(addr string, f func(stop <-chan struct{}) error) {
	stop := WaitSignals()
	quit := make(chan struct{})
	go func() {
		HTTPHealthCheck(addr, quit)
	}()

	if err := f(stop); err != nil {
		log.Errorf("%+v", err)
	}
	quit <- struct{}{}
}

// HTTPHealthCheck HTTP 模式健康检查，会阻塞执行
func HTTPHealthCheck(addr string, stop <-chan struct{}) {
	server := &http.Server{Addr: addr, Handler: healthcheck.NewHandler()}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("[BgHealthCheck] health server close with err: %+v", err)
		}
	}()

	<-stop
	server.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("[BgHealthCheck] stop server graceful stop with err: %+v", err)
	}
}

// WaitSignals 监听退出信号
func WaitSignals() chan struct{} {
	stop := make(chan struct{})

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		<-quit
		close(stop)
	}()

	return stop
}
