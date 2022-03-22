package netutils

import (
	"fmt"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	succ := Ping("192.123.12.3", 1*time.Second)
	fmt.Println(succ)
}
