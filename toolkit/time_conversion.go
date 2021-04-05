package toolkit

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
	//"strconv"
)

// input  :"2016-07-27 08:46:15"

const (
	TFORMAT_YMDHMS         = "2006-01-02 15:04:05"
	TFORMAT_HM             = "15:04"
	TFORMAT_YMDHMS_M_PLAIN = "20060102150405"
	TFORMAT_YMDHM          = "200601021504"
)

func Str2Time(str string) (time.Time, error) {

	var err error

	// ts
	tmp := strings.Split(str, ",")
	if len(tmp) != 2 {
		return time.Time{}, errors.New("comma not found")
	}

	t0, err := time.Parse(TFORMAT_YMDHMS, tmp[0])
	if err != nil {
		return time.Time{}, err
	}
	dur, err := time.ParseDuration(tmp[1] + "ms")
	if err != nil {
		return time.Time{}, err
	}
	//	fmt.Println(tmp[1]+ "ms", ", dur=", dur.Nanoseconds())

	t0 = t0.Add(dur)

	return t0, nil
}

func MsStr2Time(str string) (time.Time, error) {

	var err error

	// ts
	tmp := strings.Split(str, ".")
	if len(tmp) != 2 {
		return time.Time{}, errors.New("comma not found")
	}

	t0, err := time.Parse(TFORMAT_YMDHMS, tmp[0])
	if err != nil {
		return time.Time{}, err
	}
	dur, err := time.ParseDuration(tmp[1] + "ms")
	if err != nil {
		return time.Time{}, err
	}
	//	fmt.Println(tmp[1]+ "ms", ", dur=", dur.Nanoseconds())

	t0 = t0.Add(dur)

	return t0, nil
}
func Time2MStr(t time.Time) string {

	s := t.Format(TFORMAT_YMDHMS_M_PLAIN)
	mil := t.Nanosecond() / 1e6

	//	fmt.Println("mils= " , mil, ",", t.Nanosecond())

	return fmt.Sprintf("%s%03d", s, mil)
}

func Time2SecStr(t time.Time) string {

	return t.Format(TFORMAT_YMDHMS)
}

func DiffInSec(start string, end string) (int, error) {
	var t0 time.Time
	var t1 time.Time
	var err error

	t0, err = time.Parse(TFORMAT_YMDHMS, start)
	if err != nil {
		return 0, err
	}

	t1, err = time.Parse(TFORMAT_YMDHMS, end)
	if err != nil {
		return 0, err
	}

	diff := int((t1.Unix() - t0.Unix() + 59) / 60)

	if diff < 0 {
		return 0, errors.New("start is greater than end")
	}

	return diff, nil
}

func DiffInSeconds(start string, end string) (int, error) {
	var t0 time.Time
	var t1 time.Time
	var err error

	t0, err = time.Parse(TFORMAT_YMDHMS, start)
	if err != nil {
		return 0, err
	}

	t1, err = time.Parse(TFORMAT_YMDHMS, end)
	if err != nil {
		return 0, err
	}

	diff := t1.Unix() - t0.Unix()

	if diff < 0 {
		return 0, errors.New("start is greater than end")
	}

	return int(diff), nil
}

func GetTimeStamp(t0 time.Time) string {
	ts := t0.Format("2006-01-02 15:04:05")
	mil := t0.Nanosecond() / 1e6
	ts = ts + fmt.Sprintf(".%03d", mil)
	return ts
}

func MidnightTime(t time.Time) time.Time {

	mn := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	return mn
}
func SecondsFromMidnight(t time.Time) int64 {
	mn := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	return t.Unix() - mn.Unix()
}

func FormatFileName(s string, t time.Time) string {

	var buf bytes.Buffer

	for true {
		idx := strings.Index(s, "%")
		if idx > 0 {
			buf.WriteString(s[0:idx])
			s = s[idx:]
		} else if idx == 0 {
			if strings.Index(s, "%Y") == 0 {
				buf.WriteString(fmt.Sprintf("%d", t.Year()))
				s = s[2:]
			} else if strings.Index(s, "%M") == 0 {
				buf.WriteString(fmt.Sprintf("%02d", t.Month()))
				s = s[2:]
			} else if strings.Index(s, "%D") == 0 {
				buf.WriteString(fmt.Sprintf("%02d", t.Day()))
				s = s[2:]
			} else if strings.Index(s, "%h") == 0 {
				buf.WriteString(fmt.Sprintf("%02d", t.Hour()))
				s = s[2:]
			} else if strings.Index(s, "%m") == 0 {
				buf.WriteString(fmt.Sprintf("%02d", t.Minute()))
				s = s[2:]
			} else if strings.Index(s, "%s") == 0 {
				buf.WriteString(fmt.Sprintf("%02d", t.Second()))
				s = s[2:]
			} else {
				// skip the %
				s = s[1:]
			}
		} else {
			// idx < 0
			buf.WriteString(s)
			break
		}
	}

	return buf.String()
}

func GetDateRange(str string) (t0 time.Time, t1 time.Time, err error) {

	t := time.Now().AddDate(0, 0, -1)

	yesterday := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	s := strings.Split(str, ",")

	// get the yyyy, mm , dd of last day
	if len(s) >= 1 && s[0] != "" {
		if t0, err = time.ParseInLocation("2006-01-02", s[0], time.Local); err != nil {
			return
		}
	}

	if len(s) >= 2 && s[1] != "" {
		if t1, err = time.ParseInLocation("2006-01-02", s[1], time.Local); err != nil {
			return
		}
	}

	fmt.Println(yesterday, t0, t1)

	// if both of them are zero, ignore the other, use default
	if t0.IsZero() && t1.IsZero() {
		// use default -- yesterday
		t0 = yesterday
		t1 = t0
	} else if !t0.IsZero() && !t1.IsZero() {
		//  allowed
	} else {
		err = errors.New("Invalid startdate or enddate ")
		return
	}
	// now we have t0, t1

	if t1.Before(t0) {
		err = errors.New("end date is before start date")
		return
	}

	if yesterday.Before(t1) {
		err = errors.New("end date is after yesterday")
		return
	}

	if t0.Year() < 2018 {
		err = errors.New("start date is less than 2018")
		return

	}

	return
}

func NormalizeTs(tstr string, fmt string, interval int64) (string, error) {
	// seconds to start time of this stat slot

	tm, err := time.ParseInLocation(fmt, tstr, time.Local)
	if err != nil {
		return "", err
	}

	ts := tm.Unix()

	t := ts/3600*3600 + ((ts%3600)+interval-1)/interval*interval

	tf := time.Unix(t, 0)
	return tf.Format(fmt), nil
}
