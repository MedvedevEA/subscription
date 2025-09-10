package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

const (
	CustomDateFormat = "01-2006"
	DbFormat         = "2006-01-02"
)

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(CustomDateFormat, s)
	if err != nil {
		return err
	}
	cd.Time = t

	return nil
}
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte("null"), nil
	}
	return fmt.Appendf(nil, `"%s"`, cd.Time.Format(CustomDateFormat)), nil

}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		cd.Time = v
	case string:
		t, err := time.Parse(CustomDateFormat, v)
		if err != nil {
			return err
		}
		cd.Time = t
	case []byte:
		t, err := time.Parse(CustomDateFormat, string(v))
		if err != nil {
			return err
		}
		cd.Time = t
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	return nil
}

func (cd CustomDate) Value() (driver.Value, error) {
	if cd.IsZero() {
		return nil, nil
	}
	return cd.Format(DbFormat), nil
}
