package xerr

import (
	"context"
	"fmt"
	"strings"

	"github.com/zhaolion/gostack/util/log"
	"github.com/zhaolion/gostack/util/runtimeutil"
)

// 自定义处理过的函数放置在这里

// NoticePanic records panic
func NoticePanic(ctx context.Context, names ...string) {
	var name string
	if len(names) == 0 {
		name = runtimeutil.CallerFuncPos(1)
	} else {
		name = names[0]
	}
	if r := recover(); r != nil {
		var e error
		if err, ok := r.(error); ok {
			e = fmt.Errorf("[%s] %+v", name, err)
		} else if s, ok := r.(string); ok {
			e = New(s)
		} else {
			e = fmt.Errorf("%v", r)
		}

		log.WithError(e).Error("panic happened")
	}
}

// ReportError 顶层函数需要主动调用这个进行错误日志上报和报警
func ReportError(ctx context.Context, err error, messages ...string) error {
	if err == nil {
		return nil
	}

	originErr := Cause(err)
	if originErr == nil {
		return nil
	}

	// 忽略业务自定义错误
	_, ok := originErr.(*CustomError)
	if ok {
		return nil
	}

	e := Wrap(err, messages...)

	// 打个错误堆栈日志
	log.Ctx(ctx).Errorf("%+v", e)

	return e
}

func getOperation(skip int, operations ...string) string {
	if len(operations) == 0 {
		return runtimeutil.CallerFuncName(skip + 1)
	}
	return operations[0]
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
// 如果不传 message，默认会将调用函数写入 message
func WithMessage(err error, messages ...string) error {
	if err == nil {
		return nil
	}

	return &withMessage{
		cause: err,
		msg:   buildMessage(messages...),
	}
}

// WrapWithLog returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
// NOTE: 这个默认打错误日志哦
func WrapWithLog(err error, messages ...string) error {
	if err == nil {
		return nil
	}

	err = &withStack{
		&withMessage{
			cause: err,
			msg:   buildMessage(messages...),
		},
		callersWithErr(err),
	}

	log.Errorf("%+v", err)

	return err
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, messages ...string) error {
	if err == nil {
		return nil
	}

	return &withStack{
		&withMessage{
			cause: err,
			msg:   buildMessage(messages...),
		},
		callersWithErr(err),
	}
}

func buildMessage(messages ...string) string {
	if len(messages) == 0 {
		return runtimeutil.CallerFuncPos(2)
	}
	return strings.Join(messages, " ")
}

// StackWithLog annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
// NOTE: 这个默认打错误日志哦
func StackWithLog(err error) error {
	err = wrapStack(err, 1)
	if err != nil {
		log.Errorf("%+s", err)
	}

	return err
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
// NOTE: 这个没有打日志哦
func Stack(err error) error {
	return wrapStack(err, 1)
}

// ErrorEqual 判断 err 是否为 errors 某一个错误
func ErrorEqual(err error, errors ...error) bool {
	root := Cause(err)
	if root == nil {
		return false
	}

	for _, e := range errors {
		if root.Error() == e.Error() {
			return true
		}
	}

	return false
}
