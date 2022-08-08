package xerr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zhaolion/gostack/util/log"
	"github.com/zhaolion/gostack/util/stringutil"
)

// Custom biz custom error
func Custom(msg string) error {
	return wrapStack(&CustomError{msg: msg}, 1)
}

// Customf biz custom error, formats according to a format specifier
func Customf(format string, a ...interface{}) error {
	return wrapStack(&CustomError{msg: fmt.Sprintf(format, a...)}, 1)
}

// CustomError 业务错误代码，不应该返回 500 错误
// 同时，这个错误不会上传到 NewRelic
type CustomError struct {
	msg string
}

func (e *CustomError) Error() string {
	return e.msg
}

func (e *CustomError) CustomError() bool {
	return true
}

// IsCustomError 判断是否为业务错误
func IsCustomError(err error) (error, bool) {
	type custom interface {
		CustomError() bool
	}

	_, ok := Cause(err).(custom)
	return err, ok
}

// IgnoreDuplicateEntryError ignore duplicate entry error,
// return nil if error message contain Duplicate entry,
// return err in other errors.
func IgnoreDuplicateEntryError(err error, msg ...interface{}) error {
	if IsUniqueError(err) {
		if len(msg) != 0 {
			bs, _ := json.Marshal(msg)
			log.Debugf("ignore duplicate record %s", stringutil.BytesToString(bs))
		}

		return nil
	}

	return wrapStack(err, 1)
}

// IsUniqueError 判断是否是 sql unique 错误
func IsUniqueError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
