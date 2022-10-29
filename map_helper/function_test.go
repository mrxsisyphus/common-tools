package map_helper

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
	"testing"
)

func TestFromPairFunc(t *testing.T) {
	temp := "a=1&b=2&c=3"
	res := FromPairFunc(strings.Split(temp, "&"), func(t string) (string, int) {
		ss := strings.Split(t, "=")
		return ss[0], cast.ToInt(ss[1])
	})
	fmt.Println(res)
}
