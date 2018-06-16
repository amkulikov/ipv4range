package ipv4range

import (
	"net"
	"encoding/binary"
)

type IPv4 uint32

func netIPtoIPv4(ip net.IP) IPv4 {
	return IPv4(binary.BigEndian.Uint32(ip.To4()[:4]))
}

func (i IPv4) ToIP() (ip net.IP) {
	ip = make(net.IP, 4)
	binary.BigEndian.PutUint32(ip[:], uint32(i))
	return
}

type IPRange struct {
	Left, Right IPv4
}

func NewIPRange(left, right IPv4) *IPRange {
	return &IPRange{Left: left, Right: right}
}

func NewIPRangeByIPs(left, right net.IP) *IPRange {
	return NewIPRange(
		netIPtoIPv4(left),
		netIPtoIPv4(right),
	)
}

func (r *IPRange) Contains(ip IPv4) bool {
	return ip >= r.Left && ip <= r.Right
}

func (r *IPRange) ContainsIP(ip net.IP) bool {
	t := netIPtoIPv4(ip)
	return t >= r.Left && t <= r.Right
}

func (r *IPRange) Subnets() (nets []*net.IPNet) {
	const (
		maxMask = 32
		allOnes = ^IPv4(0)
	)

	left := r.Left
	right := r.Right

	for left <= right {
		maskSize := maxMask
		mask := allOnes
		for mask > 0 {
			temp := mask << 1
			if (left&temp) != left || (left | ^temp) > right {
				break
			}
			mask = temp
			maskSize--
		}
		nets = append(nets, &net.IPNet{IP: left.ToIP(), Mask: net.CIDRMask(maskSize, maxMask)})
		left |= ^mask
		if left+1 < left {
			break
		}
		left++
	}
	return
}
