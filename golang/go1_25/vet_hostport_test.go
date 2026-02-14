package go1_25

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HostPort_잘못된방식_Sprintf(t *testing.T) {
	// ❌ 잘못된 방식: IPv6 주소에서 문제 발생
	// Go 1.25 vet의 hostport 분석기가 이 패턴을 감지함
	ipv4 := "127.0.0.1"
	ipv6 := "::1"
	port := 8080

	// IPv4는 우연히 동작하지만...
	resultV4 := fmt.Sprintf("%s:%d", ipv4, port)
	assert.Equal(t, "127.0.0.1:8080", resultV4)

	// IPv6는 잘못된 결과! (대괄호가 없음)
	resultV6 := fmt.Sprintf("%s:%d", ipv6, port)
	assert.Equal(t, "::1:8080", resultV6, "IPv6 주소가 대괄호로 감싸지지 않음")
}

func Test_HostPort_올바른방식_JoinHostPort(t *testing.T) {
	// ✅ 올바른 방식: net.JoinHostPort 사용
	ipv4 := "127.0.0.1"
	ipv6 := "::1"
	port := 8080

	// IPv4: 정상 동작
	resultV4 := net.JoinHostPort(ipv4, strconv.Itoa(port))
	assert.Equal(t, "127.0.0.1:8080", resultV4)

	// IPv6: 자동으로 대괄호 추가
	resultV6 := net.JoinHostPort(ipv6, strconv.Itoa(port))
	assert.Equal(t, "[::1]:8080", resultV6, "IPv6 주소가 대괄호로 올바르게 감싸져야 한다")
}

func Test_HostPort_다양한_주소(t *testing.T) {
	tests := []struct {
		host     string
		port     int
		expected string
	}{
		{"localhost", 80, "localhost:80"},
		{"127.0.0.1", 443, "127.0.0.1:443"},
		{"::1", 8080, "[::1]:8080"},
		{"2001:db8::1", 3000, "[2001:db8::1]:3000"},
		{"example.com", 9090, "example.com:9090"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s:%d", tt.host, tt.port), func(t *testing.T) {
			result := net.JoinHostPort(tt.host, strconv.Itoa(tt.port))
			assert.Equal(t, tt.expected, result)
		})
	}
}
