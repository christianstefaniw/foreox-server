package errors

import (
	"fmt"
	"log"
	"runtime/debug"
)

type Error struct {
	Inner               error
	Message, StackTrace string
	Misc                map[string]interface{}
}

type ErrorKey int64

var errorKey = ErrorKey(0)

func GetErrorKey() ErrorKey {
	return errorKey
}

func Wrap(err error, message string, msgArgs ...interface{}) Error {
	errorKey++
	return Error{
		Inner:      err,
		Message:    fmt.Sprintf(message, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

func PrintError(key ErrorKey, err Error) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, err.Message)
}

func (err Error) UnWrap() error {
	return err.Inner
}

func (err Error) Error() string {
	return err.Message
}
