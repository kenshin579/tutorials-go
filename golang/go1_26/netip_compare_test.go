package go1_26_test

import (
	"fmt"
	"net/netip"
	"slices"
	"testing"
)

func TestNetipPrefixCompare(t *testing.T) {
	// CIDR 표기법 서브넷 목록
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("192.168.1.0/24"),
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("172.16.0.0/12"),
		netip.MustParsePrefix("10.0.1.0/24"),
		netip.MustParsePrefix("192.168.0.0/16"),
	}

	fmt.Println("--- 정렬 전 ---")
	for _, p := range prefixes {
		fmt.Println(" ", p)
	}

	// Go 1.26: netip.Prefix.Compare로 정렬
	slices.SortFunc(prefixes, netip.Prefix.Compare)

	fmt.Println("--- 정렬 후 ---")
	for _, p := range prefixes {
		fmt.Println(" ", p)
	}

	// 첫 번째는 10.0.0.0/8이어야 함
	if prefixes[0].String() != "10.0.0.0/8" {
		t.Errorf("expected first to be 10.0.0.0/8, got %s", prefixes[0])
	}
}

func TestNetipPrefixCompareIPv6(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("fe80::/10"),
		netip.MustParsePrefix("::1/128"),
		netip.MustParsePrefix("2001:db8::/32"),
		netip.MustParsePrefix("::/0"),
	}

	slices.SortFunc(prefixes, netip.Prefix.Compare)

	fmt.Println("--- IPv6 정렬 결과 ---")
	for _, p := range prefixes {
		fmt.Println(" ", p)
	}

	// ::/0 이 첫 번째
	if prefixes[0].String() != "::/0" {
		t.Errorf("expected first to be ::/0, got %s", prefixes[0])
	}
}

func TestNetipPrefixCompareMixed(t *testing.T) {
	// IPv4, IPv6 혼합 정렬
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("192.168.0.0/16"),
		netip.MustParsePrefix("::1/128"),
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("fe80::/10"),
	}

	slices.SortFunc(prefixes, netip.Prefix.Compare)

	fmt.Println("--- IPv4+IPv6 혼합 정렬 ---")
	for _, p := range prefixes {
		fmt.Println(" ", p)
	}

	// IPv4가 IPv6보다 먼저 정렬됨
	if !prefixes[0].Addr().Is4() {
		t.Error("expected IPv4 addresses to sort before IPv6")
	}
}
