package dbanon

import (
	"github.com/sirupsen/logrus/hooks/test"
	"strconv"
	"testing"
	"time"
)

func TestGetForEtcCases(t *testing.T) {
	testLogger, hook := test.NewNullLogger()
	SetLogger(testLogger)

	val := "test"

	provider := NewProvider()
	_ = provider.Get("faker.Whoops1", &val)
	if hook.LastEntry().Message != "faker.Whoops1 does not match any known type" {
		t.Errorf("Unsupported provider not handled correctly")
	}

	_ = provider.Get("faker.Number().Whoops2()", &val)
	if hook.LastEntry().Message != "faker.Number().Whoops2() does not match any known type" {
		t.Errorf("Unsupported method not handled correctly")
	}

	to := time.Now()
	from := to.AddDate(-40, 0, 0)
	r11a := provider.Get("datetime", &val)
	r11Time, _ := time.Parse("2006-01-02 15:04:05", r11a)
	if r11Time.Before(from) || r11Time.After(to) {
		t.Errorf("%v not in expected range [%v, %v]", r11Time, from, to)
	}
}

func TestGetForLengthBasedOptions(t *testing.T) {
	provider := NewProvider()
	tests := map[string]struct {
		input  string
		wantGt int
		wantLt int
		isStr  bool
	}{
		"md5":                            {input: "md5", wantGt: 31, wantLt: 33, isStr: true},
	}

	val := "test"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := provider.Get(tc.input, &val)
			var compare int
			if tc.isStr {
				compare = len(got)
			} else {
				compare, _ = strconv.Atoi(got)
			}

			if compare < tc.wantGt || compare > tc.wantLt {
				t.Errorf("Expected %v to be greater than %v and less than %v", compare, tc.wantGt, tc.wantLt)
			}
		})
	}
}

func TestDynamicEmailProviderIgnoresEmail(t *testing.T) {
	provider := NewProvider()

	current := "blam@test.com"
	value := provider.Get("dynamic.email(@test.com)", &current)
	if current != value {
		t.Errorf("new:%s does not equal original:%s when a fake was not expected", value, current)
	}
}

func TestDynamicEmailProviderFakesEmail(t *testing.T) {
	provider := NewProvider()
	
	current := "blam@blam.com"
	value := provider.Get("dynamic.email(@test.com)", &current)
	if current == value {
		t.Errorf("new:%s equals original:%s when a fake was expected", value, current)
	}
}

func TestDynamicEmailProviderHandlesNoArgs(t *testing.T) {
	provider := NewProvider()
	
	current := "blam@blam.com"
	faked1 := provider.Get("dynamic.email()", &current)
	faked2 := provider.Get("dynamic.email", &current)
	if faked1 == current || faked2 == current {
		t.Error("unhandled empty args")
	}
}
