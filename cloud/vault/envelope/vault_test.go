package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDEKData_Valid(t *testing.T) {
	current, wrapped, err := parseDEKData(map[string]interface{}{
		"current_version": "2",
		"wrapped_v1":      "vault:v1:a",
		"wrapped_v2":      "vault:v1:b",
	})
	require.NoError(t, err)
	assert.Equal(t, uint8(2), current)
	assert.Equal(t, map[uint8]string{1: "vault:v1:a", 2: "vault:v1:b"}, wrapped)
}

func TestParseDEKData_MissingCurrentField(t *testing.T) {
	_, _, err := parseDEKData(map[string]interface{}{
		"current_version": "2",
		"wrapped_v1":      "vault:v1:a",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wrapped_v2")
}

func TestParseDEKData_NonNumericVersion(t *testing.T) {
	_, _, err := parseDEKData(map[string]interface{}{
		"current_version": "abc",
		"wrapped_v1":      "vault:v1:a",
	})
	assert.Error(t, err)
}

func TestParseDEKData_VersionOutOfRange(t *testing.T) {
	for _, v := range []string{"0", "256"} {
		t.Run(v, func(t *testing.T) {
			_, _, err := parseDEKData(map[string]interface{}{
				"current_version": v,
				"wrapped_v1":      "vault:v1:a",
			})
			assert.Error(t, err)
		})
	}
}

func TestNextVersion_UsesHighestExistingVersion(t *testing.T) {
	next, err := nextVersion(map[uint8]string{})
	require.NoError(t, err)
	assert.Equal(t, uint8(1), next)

	// current_version 을 2 로 되돌린 뒤에도 v3 를 덮어쓰지 않도록 가장 큰 버전 다음을 쓴다
	next, err = nextVersion(map[uint8]string{1: "a", 2: "b", 3: "c"})
	require.NoError(t, err)
	assert.Equal(t, uint8(4), next)

	_, err = nextVersion(map[uint8]string{255: "z"})
	assert.Error(t, err)
}

func TestParseDEKData_LeadingZeroVersion(t *testing.T) {
	_, _, err := parseDEKData(map[string]interface{}{
		"current_version": "1",
		"wrapped_v1":      "vault:v1:a",
		"wrapped_v01":     "vault:v1:b",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wrapped_v01")
}
