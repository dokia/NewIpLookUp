/*
 * log日志输出接口
 */
package lookuptree

import (
	"log"
)

func Debugf(format string, v ...interface{}) {
	log.SetPrefix("[Debug]")
	log.Printf(format, v)
}

func Debugln(format string) {
	log.SetPrefix("[Debug]")
	log.Println(format)
}

func Infof(format string, v ...interface{}) {
	log.SetPrefix("[Info]")
	log.Printf(format, v)
}

func Infoln(format string) {
	log.SetPrefix("[Info]")
	log.Println(format)
}

func Warnf(format string, v ...interface{}) {
	log.SetPrefix("[Warn]")
	log.Printf(format, v)
}

func Warnln(format string) {
	log.SetPrefix("[Warn]")
	log.Println(format)
}

func Errorf(format string, v ...interface{}) {
	log.SetPrefix("[Error]")
	log.Printf(format, v)
}

func Errorln(format string) {
	log.SetPrefix("[Error]")
	log.Println(format)
}
