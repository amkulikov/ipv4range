package ipv4range

import (
	"testing"
	"net"
)

func TestIPv4_ToIP(t *testing.T) {
	testIPs := []net.IP{
		net.ParseIP("0.0.0.0"),
		net.ParseIP("192.168.0.2"),
		net.ParseIP("255.255.255.255"),
	}
	for _, ip := range testIPs {
		temp := netIPtoIPv4(ip).ToIP()
		if !ip.Equal(temp) {
			t.Errorf("got %s, %s expected", temp, ip)
		}
	}
}

func mustCIDR(cidr string) *net.IPNet {
	_, n, _ := net.ParseCIDR(cidr)
	return n
}

func TestIPRange_Subnets(t *testing.T) {
	ipRanges := []*IPRange{
		NewIPRangeByIPs(net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")),
		NewIPRangeByIPs(net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.5")),
		NewIPRangeByIPs(net.ParseIP("0.0.0.0"), net.ParseIP("255.255.255.255")),
		NewIPRangeByIPs(net.ParseIP("0.0.0.0"), net.ParseIP("4.0.0.0")),
	}
	expectedNets := [][]*net.IPNet{
		{mustCIDR("192.168.0.0/16")},
		{mustCIDR("192.168.0.1/32"), mustCIDR("192.168.0.2/31"), mustCIDR("192.168.0.4/31")},
		{mustCIDR("0.0.0.0/0")},
		{mustCIDR("0.0.0.0/6"), mustCIDR("4.0.0.0/32")},
	}

	for i, r := range ipRanges {
		subnets := r.Subnets()
		expected := expectedNets[i]
		for j, s := range subnets {
			if j >= len(expected) {
				t.Fatal("Too much subnets")
			}
			if s.String() != expected[j].String() {
				t.Errorf("got %s, %s expected", s, expected[j])
			}
		}
	}
}
