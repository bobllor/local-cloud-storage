package diskwriter

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/gologger"
)

func TestUndersizedChunkWrite(t *testing.T) {
	dir := t.TempDir()
	size := 512

	file := filepath.Join(dir, "test1.txt")

	dw := newTestDiskWriter(size)
	b := getTestByte(size)

	info, err := dw.WriteToDisk(file, b)
	assert.Nil(t, err)

	assert.Equal(t, info.Size(), int64(size))
}

func TestOversizeChunkWrite(t *testing.T) {
	dir := t.TempDir()
	chunkSize := 512

	file := filepath.Join(dir, "test1.txt")

	byteSize := 10240
	dw := newTestDiskWriter(chunkSize)
	b := getTestByte(byteSize)

	info, err := dw.WriteToDisk(file, b)
	assert.Nil(t, err)

	assert.Equal(t, info.Size(), int64(byteSize))
}

func TestVeryUndersizedChunkWrite(t *testing.T) {
	dir := t.TempDir()
	chunkSize := 1048

	byteSize := 30

	dw := newTestDiskWriter(chunkSize)
	b := getTestByte(byteSize)

	info, err := dw.WriteToDisk(dir+"/"+"test.txt1", b)
	assert.Nil(t, err)

	assert.Equal(t, info.Size(), int64(byteSize))
}

func TestFailWriteDirectoryPath(t *testing.T) {
	dir := t.TempDir()
	chunkSize := 512

	dw := newTestDiskWriter(chunkSize)
	b := getTestByte(2048)

	_, err := dw.WriteToDisk(dir, b)
	assert.NotNil(t, err)
}

func TestFailWriteNoData(t *testing.T) {
	dir := t.TempDir()
	chunkSize := 512

	dw := newTestDiskWriter(chunkSize)

	_, err := dw.WriteToDisk(dir, []byte{})
	assert.NotNil(t, err)
}

// newTestDiskWriter creates a new DiskWriter with a testing setup.
func newTestDiskWriter(chunkSize int) *DiskWriter {
	logger := gologger.NewLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime), gologger.Lsilent)

	config := config.NewConfig(logger)

	dw := NewDiskWriter(chunkSize, config)

	return dw
}

// getTestByte creates a new byte slice of arbitrary values based on
// the given size.
func getTestByte(size int) []byte {
	b := make([]byte, 0, size)
	r := rand.New(rand.NewSource(time.Now().Unix()))

	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	chars += strings.ToLower(chars)

	for range size {
		ranCh := r.Intn(len(chars))
		b = append(b, chars[ranCh])
	}

	return b
}
