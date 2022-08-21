package netutils

import (
	"errors"
	"strconv"
	"strings"
)

const (
	IPv4      = "IPv4"
	IPv6      = "IPv6"
	IPUnknown = "Neither"
)

func IpTo32(ipStr string) (uint32, error) {
	split := strings.Split(ipStr, ".")
	if len(split) != 4 {
		return 0, errors.New("invalid ip address")
	}
	var ip uint32 = 0
	for i := 0; i < 4; i++ {
		str := split[3-i]
		if len(str) == 0 || len(str) > 3 || len(str) > 1 && str[0] == '0' {
			return 0, errors.New("invalid ip address")
		}
		n, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
		if n>>8 != 0 {
			return 0, errors.New("invalid ip address")
		}
		ip += uint32(n) << (i * 8)
	}
	return ip, nil
}

func ValidIPAddress(queryIP string) string {
	i := strings.IndexByte(queryIP, '.')
	if i != -1 {
		_, err := IpTo32(queryIP)
		if err != nil {
			return IPUnknown
		}
		return IPv4
	} else {
		splits := strings.Split(queryIP, ":")
		empty := false
		for _, split := range splits {
			if split == "" {
				if !empty {
					empty = true
				} else {
					return IPUnknown
				}
			}
			if len(split) > 4 {
				return IPUnknown
			}
			for _, s := range split {
				if s >= 'a' && s <= 'f' || s >= 'A' && s <= 'F' || s >= '0' && s <= '9' {
					continue
				}
				return IPUnknown
			}
			parseInt, err := strconv.ParseUint(split, 16, 64)
			if err != nil || parseInt < 0 || parseInt > 0xffff {
				return IPUnknown
			}
		}
		if empty && len(splits) < 8 || len(splits) == 8 {
			return IPv6
		}
		return IPUnknown
	}
}
