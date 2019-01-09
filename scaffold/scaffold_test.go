package scaffold

import (
	"io/ioutil"
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

	assert.NoError(t, New(true).Generate(tempDir))

	//defer os.RemoveAll(tempDir) // clean up
}
