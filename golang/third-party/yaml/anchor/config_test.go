package anchor

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), ".")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	fmt.Println("dir", dir)
}

func Test_Config(t *testing.T) {
	cfg, err := New("config.yaml")
	assert.NoError(t, err)

	content, err := json.MarshalIndent(cfg, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(content))
}
