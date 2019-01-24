package utils

import (
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/panzhenyu12/flower/config"
)

func WaitForExit() os.Signal {
	return WaitForSignal(syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
}

func WaitForSignal(sources ...os.Signal) os.Signal {
	var s = make(chan os.Signal, 1)
	defer signal.Stop(s) //the second Ctrl+C will force shutdown

	signal.Notify(s, sources...)
	return <-s //blocked
}

//UnixToStr 毫秒级别时间戳转时间
func UnixToStr(timestamp int64) string {
	timestr := ""
	if timestamp <= 0 {
		return timestr
	}
	//time.Now().UnixNano() / 1000000
	//tm:=time.Unix
	tm := time.Unix(0, timestamp*1e6)
	return tm.Format("2006-01-02 15:04:05")
}

//TimeToTimestamp 时间转毫秒级别时间戳
func TimeToTimestamp(datetime time.Time) int64 {
	if datetime.IsZero() {
		return 0
	}
	return datetime.UnixNano() / 1e6
}
func GetTimeLocation() *time.Location {
	return time.FixedZone("serverzone", config.GetConfig().TimeZone*60*60)
}

func FromStringToInt64(s string) int64 {
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}

func FromInt64ToString(num int64) string {
	return strconv.FormatInt(int64(num), 10)
}

func GetDayTime() time.Time {
	//starttshour := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), 0, 0, 0, utils.GetTimeLocation())
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, GetTimeLocation())
}

func GetHourTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), 0, 0, 0, GetTimeLocation())
}

func FindNum(str string) string {
	reg := regexp.MustCompile("^[0-9]\\d*|0$")
	strs := reg.FindAllString(str, -1)
	if len(strs) == 0 {
		return str
	}
	return strs[0]
}
