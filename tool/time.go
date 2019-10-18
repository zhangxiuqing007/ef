package tool

import (
	"fmt"
	"time"
)

var weekdayStrs = [7]string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

//FormatTimeDetail 格式化时间显示,例：2019年8月10日 星期日 16:25:59
func FormatTimeDetail(t time.Time) string {
	return fmt.Sprintf("%d年%d月%d日 %s %d:%d:%d",
		t.Year(),
		t.Month(),
		t.Day(),
		weekdayStrs[t.Weekday()],
		t.Hour(),
		t.Minute(),
		t.Second())
}

//FormatTimeSimple 格式化时间显示,例：2019年8月10日
func FormatTimeSimple(t time.Time) string {
	return fmt.Sprintf("%d年%d月%d日", t.Year(), t.Month(), t.Day())
}
