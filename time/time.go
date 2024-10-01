package time

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	t0 "time"
)

var (
	Year   = "2006"
	Month  = "01"
	Day    = "02"
	Hour   = "15"
	Minute = "04"
	Second = "05"

	FmtYMdHmsSuS = "2006-01-02 15:04:05.000000"
	FmtYMdHmsS   = "2006-01-02 15:04:05.000"
	FmtYMdHms    = "2006-01-02 15:04:05"
	FmtYMdHm     = "2006-01-02 15:04"
	FmtYMdH      = "2006-01-02 15"
	FmtYMd       = "2006-01-02"
	FmtYM        = "2006-01"
	FmtY         = "2006"
	FmtYYYYMMdd  = "20060102"

	FmtCnYMdHmsSuS = "2006年01月02日 15时04分05秒000000微秒"
	FmtCnYMdHmsS   = "2006年01月02日 15时04分05秒000毫秒"
	FmtCnYMdHms    = "2006年01月02日 15时04分05秒"
	FmtCnYMdHm     = "2006年01月02日 15时04分"
	FmtCnYMdH      = "2006年01月02日 15时"
	FmtCnYMd       = "2006年01月02日"
	FmtCnYM        = "2006年01月"
	FmtCnY         = "2006年"

	FmtHmsSMore = "15:04:05.000000000"
	FmtHmsS     = "15:04:05.000"
	FmtHms      = "15:04:05"
	FmtHm       = "15:04"
	FmtH        = "15"

	EmptyTime = t0.Time{}

	yRegex         = regexp.MustCompile("^(\\d){4}$")
	yyyyMmDdRegex  = regexp.MustCompile("^(\\d){4}(\\d){2}(\\d){2}$")
	ymRegex        = regexp.MustCompile("^(\\d){4}-(\\d){2}$")
	ymdRegex       = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2}$")
	ymdHRegex      = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}$")
	ymdHmRegex     = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}$")
	ymdHmsRegex    = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}:(\\d){2}$")
	ymdHmsSRegex   = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}:(\\d){2}.(\\d){3}$")
	ymdHmsSuSRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}:(\\d){2}.(\\d){6}$")
)

type FPCDateTime float64

var DefaultLocation *t0.Location

const (
	DateDelta   = 693594 // Days between 1/1/0001 and 12/31/1899
	HoursPerDay = 24
	MinsPerHour = 60
	SecsPerMin  = 60
	MSecsPerSec = 1000

	MinsPerDay  = HoursPerDay * MinsPerHour
	SecsPerHour = SecsPerMin * MinsPerHour
	SecsPerDay  = MinsPerDay * SecsPerMin
	MSecsPerDay = SecsPerDay * MSecsPerSec

	OneMillisecond  = FPCDateTime(1) / MSecsPerDay
	HalfMilliSecond = OneMillisecond / 2

	JulianEpoch = FPCDateTime(-2415018.5)
	UnixEpoch   = JulianEpoch + FPCDateTime(2440587.5)

	ApproxDaysPerMonth = 30.4375
	ApproxDaysPerYear  = 365.25
)

func init() {
	_DefaultLocation, _ := t0.LoadLocation("Asia/Shanghai")
	DefaultLocation = _DefaultLocation
}

func TimeToStringYmd(t t0.Time) string {
	return t.Format(FmtYMd)
}

func TimeToStringYmdHms(t t0.Time) string {
	return t.Format(FmtYMdHms)
}

func TimeToStringYmdHmsS(t t0.Time) string {
	return t.Format(FmtYMdHmsS)
}

func TimeToStringFormat(t t0.Time, format string) string {
	return t.Format(format)
}

func ParseTimeYmd(timeStr string) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMd, timeStr, t0.Local)
}

func ParseTimeYmdHms(timeStr string) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMdHms, timeStr, t0.Local)
}

func ParseTimeYmdHmsS(timeStr string) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMdHmsS, timeStr, t0.Local)
}

func ParseTimeYmdHmsSuS(timeStr string) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMdHmsSuS, timeStr, t0.Local)
}

func ParseTimeYmdHmsLoc(timeStr string, loc *t0.Location) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMdHms, timeStr, loc)
}

func ParseTimeYmdHmsSLoc(timeStr string, loc *t0.Location) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMdHmsS, timeStr, loc)
}

func ParseTimeYmdHmsSusLoc(timeStr string, loc *t0.Location) (t0.Time, error) {
	return t0.ParseInLocation(FmtYMdHmsSuS, timeStr, loc)
}

// TimeInMillis 13位java时间戳
func TimeInMillis() int64 {
	return t0.Now().UnixMilli()
}

// TimeInSeconds 10位Unix时间戳
func TimeInSeconds() int64 {
	return t0.Now().Unix()
}

// TimeInMicro 16位时间戳
func TimeInMicro() int64 {
	return t0.Now().UnixMicro()
}

// TimeInNano 19位时间戳
func TimeInNano() int64 {
	return t0.Now().UnixNano()
}

//====================================================================================
// 时间日期函数
//====================================================================================

func CurrentMinuteOfDay() int {
	t := t0.Now()
	return t.Hour()*60 + t.Minute()
}

func CurrentSecondOfDay() int {
	t := t0.Now()
	return t.Hour()*3600 + t.Minute()*60 + t.Second()
}

func MinuteOfDay(t t0.Time) int {
	return t.Hour()*60 + t.Minute()
}

func SecondOfDay(t t0.Time) int {
	return t.Hour()*3600 + t.Minute()*60 + t.Second()
}

func IsLeapYear(year int) bool {
	return (year%4 == 0) && ((year%100 != 0) || (year%400 == 0))
}

func YearsBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours() / 24 / ApproxDaysPerYear))
}

func MonthsBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours() / 24 / ApproxDaysPerMonth))
}

func DaysBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours() / 24))
}

func HoursBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours()))
}

func MinutesBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Minutes()))
}

func SecondsBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Seconds()))
}

func MilliSecondsBetween(now, then t0.Time) int64 {
	return now.Sub(then).Milliseconds()
}

func WithInPastYears(now, then t0.Time, years int) bool {
	return YearsBetween(now, then) <= years
}

func WithInPastMonths(now, then t0.Time, months int) bool {
	return MonthsBetween(now, then) <= months
}

func WithInPastDays(now, then t0.Time, days int) bool {
	return DaysBetween(now, then) <= days
}

func WithInPastHours(now, then t0.Time, hours int) bool {
	return HoursBetween(now, then) <= hours
}

func WithInPastMinutes(now, then t0.Time, minutes int) bool {
	return MinutesBetween(now, then) <= minutes
}

func WithInPastSeconds(now, then t0.Time, seconds int) bool {
	return SecondsBetween(now, then) <= seconds
}

func WithInPastMilliSeconds(now, then t0.Time, milliSeconds int64) bool {
	return MilliSecondsBetween(now, then) <= milliSeconds
}

func YearSpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours() / 24 / ApproxDaysPerYear
}

func MonthSpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours() / 24 / ApproxDaysPerMonth
}

func DaySpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours() / 24
}

func HourSpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours()
}

func MinuteSpan(now, then t0.Time) float64 {
	return now.Sub(then).Minutes()
}

func SecondSpan(now, then t0.Time) float64 {
	return now.Sub(then).Seconds()
}

func MilliSecondSpan(now, then t0.Time) int64 {
	return now.Sub(then).Milliseconds()
}

func Now() t0.Time {
	return t0.Now()
}

func AddHour(times t0.Time, plusOrMinus string, hours string) t0.Time {
	h, _ := t0.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, hours))
	return times.Add(h)
}

func AddMinutes(times t0.Time, plusOrMinus string, minutes string) t0.Time {
	h, _ := t0.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, minutes))
	return times.Add(h)
}

func AddSeconds(times t0.Time, plusOrMinus string, seconds string) t0.Time {
	h, _ := t0.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, seconds))
	return times.Add(h)
}

func AddTime(times t0.Time, plusOrMinus string, timeStr string) t0.Time {
	h, _ := t0.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, timeStr))
	return times.Add(h)
}

func AddDays(times t0.Time, days int) t0.Time {
	return times.AddDate(0, 0, days)
}

func AddMonths(times t0.Time, month int) t0.Time {
	return times.AddDate(0, month, 0)
}

func AddYears(times t0.Time, year int) t0.Time {
	return times.AddDate(year, 0, 0)
}

func ParseTime(timeStr string) (t0.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	timeStr = strings.TrimSpace(strings.ReplaceAll(timeStr, "\\'", " "))

	if timeStr == "" {
		return EmptyTime, errors.New("时间字段为空")
	}
	if yRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtY, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if yyyyMmDdRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYYYYMMdd, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYM, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymdRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYMd, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymdHRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYMdH, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymdHmRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYMdHm, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymdHmsRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYMdHms, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymdHmsSRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYMdHmsS, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else if ymdHmsSuSRegex.MatchString(timeStr) {
		if times, err := t0.ParseInLocation(FmtYMdHmsSuS, timeStr, DefaultLocation); err == nil {
			return times, nil
		} else {
			log.Printf("解析时间错误, err: %v", err)
			return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, err: %v", err))
		}
	} else {
		log.Printf("解析时间错误, time: %v", timeStr)
		return EmptyTime, errors.New(fmt.Sprintf("解析时间错误, time: %v", timeStr))
	}
	return EmptyTime, errors.New("内部解析器格式暂时不匹配")
}

func IsTimeEmpty(time t0.Time) bool {
	return time == EmptyTime
}

func NumToTimeDuration(num int, duration t0.Duration) t0.Duration {
	int64Num, _ := strconv.ParseInt(fmt.Sprintf("%v", num), 10, 64)
	return t0.Duration(int64Num * duration.Nanoseconds())
}

func MillisecondToTime(milliseconds int64) t0.Time {
	seconds := milliseconds / 1000
	nanoseconds := (milliseconds % 1000) * 1e6
	return t0.Unix(seconds, nanoseconds)
}

// GetMiddleTime 获取两个时间的中间时间
func GetMiddleTime(t1, t2 t0.Time) t0.Time {
	return t1.Add(t2.Sub(t1) / 2)
}

// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
const baseOriginSecond = 1136185445

// 2006-01-02 距离 1900-01-01的天数
const baseDiffDay = 38719

// ExcelDateToTime 将 Excel 日期序列号（示例：45342）转换为 Go 的 time.Time 对象
func ExcelDateToTime(excelDate int) string {
	return t0.Unix(int64(baseOriginSecond+(excelDate-baseDiffDay)*24*3600), 0).Format(FmtCnYMd)
}
