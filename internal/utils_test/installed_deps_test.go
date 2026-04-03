package utils_test

import (
	"reflect"
	"testing"

	"github.com/t-shah02/pacview/internal/utils"
)

func TestParseQiStdout_sample(t *testing.T) {
	const fixture = `Name            : foo
Version         : 1-1
Description     : one
Install Date    : Mon 01 Jan 2024
Depends On      : bar  baz
Required By     : None

Name            : bar
Version         : 2-2
Description     : two
Install Date    : Tue 02 Jan 2024
Depends On      : None
Required By     : foo
`

	got := utils.ParseQiStdout(fixture)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}

	want0 := utils.PacmanPackage{
		Name:        "foo",
		Version:     "1-1",
		Description: "one",
		InstalledAt: "Mon 01 Jan 2024",
		DependsOn:   []string{"bar", "baz"},
		RequiredBy:  nil,
	}
	if !reflect.DeepEqual(got[0], want0) {
		t.Errorf("pkg[0] = %#v\nwant %#v", got[0], want0)
	}

	want1 := utils.PacmanPackage{
		Name:        "bar",
		Version:     "2-2",
		Description: "two",
		InstalledAt: "Tue 02 Jan 2024",
		DependsOn:   nil,
		RequiredBy:  []string{"foo"},
	}
	if !reflect.DeepEqual(got[1], want1) {
		t.Errorf("pkg[1] = %#v\nwant %#v", got[1], want1)
	}
}
