package netutils

import (
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	Ping("192.168.10.22", 1 * time.Second)
}
