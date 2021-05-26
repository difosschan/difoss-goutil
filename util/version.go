package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Version struct {
	Major    uint16 // 主版本号
	Minor    uint16 // 次版本号
	Patch    uint16 // 修订版本号
	Addition string // 附加信息
}

var (
	errVersionFormat = errors.New("version format error")
)

// @param {string} ver 字符串格式的版本号，兼容格式 [v]1.2.3{foo}
func NewVersion(ver string) *Version {
	v := Version{}
	if err := v.Set(ver); err != nil {
		return nil
	}
	return &v
}

func (v *Version) FromDB(data []byte) error {
	return v.Set(string(data))
}

func (v *Version) ToDB() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) Set(ver string) error {
	reg := regexp.MustCompile(`[v|V]?(\d+)\.(\d+)\.(\d+)(\S*)$`)
	if isMatch := reg.MatchString(ver); ! isMatch {
		return errVersionFormat
	}

	match := reg.FindStringSubmatch(ver)
	var tmp uint64
	tmp, _ = strconv.ParseUint(match[1], 10, 16)
	v.Major = uint16(tmp)
	tmp, _ = strconv.ParseUint(match[2], 10, 16)
	v.Minor = uint16(tmp)
	tmp, _ = strconv.ParseUint(match[3], 10, 16)
	v.Patch = uint16(tmp)
	v.Addition = match[4]
	return nil
}

func (v *Version) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("v%d.%d.%d%s",
		v.Major, v.Minor, v.Patch, v.Addition)
}

func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}
	return v.Set(data)
}

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

func (v *Version) UnmarshalJSON(data []byte) error {
	l := len(data)
	s := string(data)
	if l == 0 || s == `""` {
		return errVersionFormat
	}
	if l < 2 || data[0] != '"' || data[l-1] != '"' {
		return errVersionFormat
	}
	return v.Set(s)
}

func (v *Version) IsZero() bool {
	return v == nil || (v.Major == 0 && v.Minor == 0 && v.Patch == 0)
}

func (v *Version) GreaterThan(w *Version) bool {
	if v == nil {
		return false
	}
	if w.IsZero() {
		// any (non-zero) > zero
		return ! v.IsZero()
	}
	// when Major is greater, return true immediately, vice versa. But both equal..
	if v.Major > w.Major {
		return true
	} else if v.Major < w.Major {
		return false
	}

	// when Minor is greater, return true immediately, vice versa. But both equal..
	if v.Minor > w.Minor {
		return true
	} else if v.Minor < w.Minor {
		return false
	}

	// when Patch is greater, return true immediately, vice versa. But both equal..
	if v.Patch > w.Patch {
		return true
	} else if v.Patch < w.Patch {
		return false
	}

	return false
}

func (v *Version) Equal(w *Version) bool {
	if w.IsZero() {
		return v.IsZero()
	}
	return v.Major == w.Major && v.Minor == w.Minor && v.Patch == w.Patch
}
