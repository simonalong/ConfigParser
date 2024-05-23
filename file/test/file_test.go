package test

import (
	"github.com/magiconair/properties/assert"
	"github.com/simonalong/gole/isc"
	"os"
	"path/filepath"
	"testing"

	"github.com/simonalong/gole/file"
)

func TestFile(t *testing.T) {
	// file.WriteFile("./sample.txt", "aaa")
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "sample.txt")

	file.AppendFile(path, "ccc")
	file.DeleteFile("sample.txt")
}

func TestExtract(t *testing.T) {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "sample.txt")
	p0 := file.ExtractFilePath(path)
	t.Logf("p0: %s", p0)
	n0 := file.ExtractFileName(path)
	t.Logf("n0: %s", n0)
	e0 := file.ExtractFileExt(path)
	t.Logf("e0: %s", e0)
	c0 := file.ChangeFileExt(path, "xyz")
	t.Logf("c0: %s", c0)
}

func TestCreatFile(t *testing.T) {
	file.CreateFile("./file/test.txt")
	file.CreateFile("./test2.txt")

	file.DeleteFile("./test2.txt")
	file.DeleteDirs("./file/")
}

func TestChild(t *testing.T) {
	f, _ := file.Child("../")
	for i := range f {
		t.Logf("file_name: %s", f[i].Name())
	}
}

func TestFileSize(t *testing.T) {
	assert.Equal(t, isc.ToInt64(40), file.Size("./assert_file_size.txt"))
}

func TestFileFormatSize(t *testing.T) {
	assert.Equal(t, "40.00B", file.SizeFormat("./assert_file_size.txt"))
}
