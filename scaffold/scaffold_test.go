package scaffold

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScaffold(t *testing.T) {

	tempDir, err := ioutil.TempDir(filepath.Join(Gopath, "src"), "test")

	if !filepath.IsAbs(tempDir) {
		tempDir, err = filepath.Abs(tempDir)
		assert.NoError(t, err)
	}

	fmt.Printf("tempDir:%s\n", tempDir)
	assert.NoError(t, New(true).Generate(tempDir))

	defer os.RemoveAll(tempDir) // clean up
}
