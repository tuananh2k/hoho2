package library

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// GetCurrentTimeStamp returns the current timestamp in the format "2006-01-02 15:04:05"
func GetCurrentTimeStamp() string {
	// Get the current time
	currentTime := time.Now()

	// Format the time using the standard time format
	return currentTime.Format("2006-01-02 15:04:05")
}

// GetCurrentDate returns the current timestamp in the format "2006-01-02"
func GetCurrentDate() string {
	// Get the current time
	currentTime := time.Now()

	// Format the time using the standard time format
	return currentTime.Format("2006-01-02")
}

func MakeNanoTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}
func FormatDatabaseDateTime(dateTime string) string {
	var dt string
	t, err := time.Parse("2006-01-02T15:04:05Z", dateTime)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05+07:00", dateTime)
		if err != nil {
			return dateTime
		}
		t = t.Add(time.Hour * time.Duration(7))
	}
	dt = t.Format("2006-01-02 15:04:05")
	return dt
}
func FormatDatabaseDate(date string) string {
	var dt string
	t, err := time.Parse("2006-01-02T15:04:05Z", date)

	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05+07:00", date)
		if err != nil {
			return date
		}
		t = t.Add(time.Hour * time.Duration(7))
	}
	dt = t.Format("2006-01-02")
	return dt
}
func FormatDateTimeByPattern(dateTime string, oldPattern string, pattern string) string {
	var dt string
	t, err := time.Parse(oldPattern, dateTime)
	if err != nil {
		return dateTime
	}
	dt = t.Format(pattern)
	return dt
}

// Hàm kiểm tra 1 phần tử có trong slice hay không
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetTimeMiliseconds(dateTime string) int64 {
	if dateTime == "" {
		return UnixMilli(time.Now())
	}
	dateTimeS := strings.Split(dateTime, " ")
	dates := dateTimeS[0]
	times := dateTimeS[1]
	dates1 := strings.Split(dates, "-")
	times1 := strings.Split(times, ":")
	y, _ := strconv.Atoi(dates1[0])
	m, _ := strconv.Atoi(dates1[1])
	d, _ := strconv.Atoi(dates1[2])
	hh, _ := strconv.Atoi(times1[0])
	mm, _ := strconv.Atoi(times1[1])
	ss, _ := strconv.Atoi(times1[2])
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	t := time.Date(y, time.Month(m), d, hh, mm, ss, 0, loc)
	mili := UnixMilli(t)
	return mili
}

func UnixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Difference returns the difference of two string slices
func Difference(a, b []string) []string {
	// Create a map with b's elements
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}

	// Initialize diff to an empty string slice
	var diff []string

	// Compare and filter elements from a that are not present in mb
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}

// IsUrl takes a string as input and checks if the given string is a valid URL
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Hàm lấy goroutine process id
func Goid() (int, error) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	return id, err
}

// HashString takes an input string "s" and returns a hash value of the input string
func HashString(s string) string {
	// calculate the MD5 sum of the input bytearray
	h := md5.Sum([]byte(s))

	// return the hash value in a hexadecimal format
	return fmt.Sprintf("%x", h)
}

func ToString(data interface{}) string {
	if data == nil {
		return ""
	}
	return fmt.Sprintf("%v", data)
}

func GetSubStringInSingleQuote(str string) []string {
	rsl := make([]string, 0)
	re := regexp.MustCompile(`'((?:''|[^'])*)'`) // Define the pattern to match anything inside single quotes

	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		rsl = append(rsl, fmt.Sprintf("'%v'", match[1]))
	}
	return rsl
}

func ExtractSubstring(s string, openStr rune, closeStr rune) []string {
	var results []string
	stack := make([]int, 0)

	for i, c := range s {
		if c == openStr {
			stack = append(stack, i)
		} else if c == closeStr {
			if len(stack) > 0 {
				start := stack[len(stack)-1] + 1
				if len(stack) == 1 {
					results = append(results, s[start:i])
				}
				stack = stack[:len(stack)-1]
			}
		}
	}
	return results
}

func FindAllIndices(s string, substr string) []int {
	var indices []int

	index := 0
	i := 0
	for {
		// get index by regex
		re := regexp.MustCompile("\\b" + substr + "\\b")
		indexs := re.FindStringIndex(s[i:])
		if indexs == nil {
			break
		}
		index = indexs[0]

		if index == -1 {
			break
		}

		indices = append(indices, index+i)
		i += index + len(substr)
	}
	return indices
}

func GetSubStringByFunction(str string, funcName string, openStr rune, closeStr rune, includeFuncName bool) []string {
	rsl := make([]string, 0)
	mapPlaceHolder := make(map[string]string)

	// Tạm comment đoạn bên dưới do đang bị conflict về cú pháp: với trường hợp 'ref(...)' thì ref(...) bị coi là chuỗi, thay vì function, nên bị đưa vào placholder thay vì đi tìm
	// thay thế các chuỗi trong dấu nháy đơn bằng placeholder
	// strInSingleQoute := GetSubStringInSingleQuote(str)
	// for idx, v := range strInSingleQoute {
	// 	key := fmt.Sprintf("_SYMPER_PLACE_HOLDER_SINGLE_QUOUTE_%d_", idx)
	// 	mapPlaceHolder[key] = v
	// 	str = strings.ReplaceAll(str, v, key)
	// }

	// find all index of funcName in str
	indices := FindAllIndices(str, funcName)
	if len(indices) == 0 {
		return rsl
	}

	// Lấy các chuỗi con nằm giữa openStr và closeStr
	for i := 0; i < len(indices); i++ {
		if i == len(indices)-1 {
			strToFind := str[indices[i]:]
			rsl = append(rsl, ExtractSubstring(strToFind, openStr, closeStr)...)
			break
		}
		strToFind := str[indices[i]:indices[i+1]]
		rsl = append(rsl, ExtractSubstring(strToFind, openStr, closeStr)...)
	}

	// trả lại placeholder bằng chuỗi ban đầu
	for idx := range rsl {
		for k, v := range mapPlaceHolder {
			rsl[idx] = strings.ReplaceAll(rsl[idx], k, v)
		}
		if includeFuncName {
			rsl[idx] = fmt.Sprintf("%s%s%s%s", funcName, string(openStr), rsl[idx], string(closeStr))
		}
	}

	return rsl
}
