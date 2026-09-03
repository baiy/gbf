package times

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtc(t *testing.T) {
	time.Local = time.UTC

	tests := []struct {
		layout   string
		str      string
		expected string
	}{
		{time.DateTime, "2023-04-01 12:00:00", "2023-04-01 12:00:00"},
		{"2006-01-02 15:04:05.999999999", "2023-04-01 12:00:00.123456", "2023-04-01 12:00:00.123456"},
		{time.RFC3339, "2023-04-01T01:00:00Z", "2023-04-01T01:00:00Z"},
		{time.RFC3339, "2023-04-01T01:00:00+01:00", "2023-04-01T00:00:00Z"},
		{time.RFC3339, "2023-04-01T01:00:00-01:00", "2023-04-01T02:00:00Z"},
		{time.RFC3339Nano, "2023-04-01T00:00:00.123Z", "2023-04-01T00:00:00.123Z"},
		{time.RFC3339Nano, "2023-04-01T01:00:00.123+01:00", "2023-04-01T00:00:00.123Z"},
		{time.RFC3339Nano, "2023-04-01T01:00:00.123-01:00", "2023-04-01T02:00:00.123Z"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			pt, err := Parse(tt.str)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, tt.expected, pt.Format(tt.layout))
				assert.Equal(t, time.Local.String(), pt.Location().String())
			}
		})
	}
}

func TestChina(t *testing.T) {
	tz, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = tz

	tests := []struct {
		layout   string
		str      string
		expected string
	}{
		{time.DateTime, "2023-04-01 12:00:00", "2023-04-01 12:00:00"},
		{"2006-01-02 15:04:05.999999999", "2023-04-01 12:00:00.123456", "2023-04-01 12:00:00.123456"},
		{time.RFC3339, "2023-04-01T01:00:00Z", "2023-04-01T09:00:00+08:00"},
		{time.RFC3339, "2023-04-01T01:00:00+01:00", "2023-04-01T08:00:00+08:00"},
		{time.RFC3339, "2023-04-01T01:00:00-01:00", "2023-04-01T10:00:00+08:00"},
		{time.RFC3339Nano, "2023-04-01T00:00:00.123Z", "2023-04-01T08:00:00.123+08:00"},
		{time.RFC3339Nano, "2023-04-01T01:00:00.123+01:00", "2023-04-01T08:00:00.123+08:00"},
		{time.RFC3339Nano, "2023-04-01T01:00:00.123-01:00", "2023-04-01T10:00:00.123+08:00"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			pt, err := Parse(tt.str)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, tt.expected, pt.Format(tt.layout))
				assert.Equal(t, time.Local.String(), pt.Location().String())
			}
		})
	}
}

func TestUnix(t *testing.T) {
	tz, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = tz

	tests := []struct {
		v        int64
		expected string
	}{
		{1787555662, "2026-08-24 15:14:22.000000"},
		{-1262304000, "1930-01-01 08:00:00.000000"},

		{1787555662000, "2026-08-24 15:14:22.000000"},
		{-1262304000000, "1930-01-01 08:00:00.000000"},

		{1787555662000000, "2026-08-24 15:14:22.000000"},
		{-1262304000000000, "1930-01-01 08:00:00.000000"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.v), func(t *testing.T) {
			pt, err := Parse(tt.v)
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, tt.expected, pt.Format(time.DateTime+".000000"))
				assert.Equal(t, time.Local.String(), pt.Location().String())
			}
		})
	}
}
