package utl

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	l "github.com/stevenb256/log"
)

// AreStringSliceSame returns true if string slices are the same
func AreStringSliceSame(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, a1 := range s1 {
		if s2[i] != a1 {
			return false
		}
	}
	return true
}

// Itoa converts int to string
func Itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// Atoi - atoi that doesn't fail and instead returns 0
func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if l.Check(err) {
		return -1
	}
	return i
}

// Percent computes percent value
func Percent(v1, v2 int) int {
	if 0 == v2 {
		return 0
	}
	return int((float64(v1) / float64(v2)) * float64(100))
}

// HashBytes returns md5 hash of bytes
func HashBytes(buf []byte) string {
	h := md5.Sum(buf)
	return base64.StdEncoding.EncodeToString(h[:])
}

// HashString - gets a hash of a string
func HashString(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%X", h.Sum32())
}

// MinUint16 min of two uint16
func MinUint16(x, y uint16) uint16 {
	if x < y {
		return x
	}
	return y
}

// MaxUint16 Max of two uint32
func MaxUint16(x, y uint16) uint16 {
	if x > y {
		return x
	}
	return y
}

// SendError - send error over channel if not nil
func SendError(chError chan error, err error) error {
	if nil != chError {
		chError <- err
	}
	return err
}

// Execute - runs a command
func Execute(wait bool, dir, app string, args ...string) error {

	// change directory
	err := os.Chdir(dir)
	if l.Check(err) {
		return err
	}

	// start command
	command := exec.Command(app, args...)

	// give current in/out
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// wait
	if true == wait {
		err = command.Run()
		if l.Check(err) {
			return err
		}
	} else {
		err = command.Start()
		if l.Check(err) {
			return err
		}
	}

	// done
	return nil
}

// MoveFile - copies file and then deletes the source file
func MoveFile(srcPath, dstPath string) error {
	err := CopyFile(srcPath, dstPath)
	if l.Check(err) {
		return err
	}
	return os.Remove(srcPath)
}

// CopyFileWithJoin - same as copy file but joins src/dst paths with src/dst file names
func CopyFileWithJoin(srcPath, srcFile, dstPath, dstFile string) error {
	return CopyFile(Join(srcPath, srcFile), Join(dstPath, dstFile))
}

// CopyFile copys source file path to dst file path
func CopyFile(srcPath, dstPath string) error {

	// make path
	err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if l.Check(err) {
		return err
	}

	// open source
	srcFile, err := os.Open(srcPath)
	if l.Check(err) {
		return err
	}
	defer srcFile.Close()

	// kill destination
	os.Remove(dstPath)

	// get info of source
	info, err := srcFile.Stat()
	if l.Check(err) {
		return err
	}

	// write to dest
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR, info.Mode())
	if l.Check(err) {
		return err
	}
	defer dstFile.Close()

	// make sure mode really got set
	if runtime.GOOS != "windows" {
		err = dstFile.Chmod(info.Mode())
		if l.Check(err) {
			return err
		}
	}

	// copy it
	_, err = io.Copy(dstFile, srcFile)
	if l.Check(err) {
		os.Remove(dstPath)
		return err
	}

	// done
	return nil
}

// DoesFileExist - checks to see if file exists
func DoesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return true != os.IsNotExist(err)
}

// IsDirectory - checks to see if path is a director or a file
func IsDirectory(path string) bool {
	stat, err := os.Stat(path)
	return nil == err && true == stat.IsDir()
}

// Join - takes a set of strings and joins them into a path
func Join(a ...string) string {
	return filepath.Join(a...)
}

// Clean - calls go filepath clean method
func Clean(path string) string {
	return filepath.Clean(path)
}

// WriteFile writes buffer into path
func WriteFile(path string, buffer []byte) error {

	// make sure directory exists
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if l.Check(err) {
		return err
	}

	// create/lock the local file
	file, err := os.Create(path)
	if l.Check(err) {
		return err
	}
	defer file.Close()

	// write the contents
	_, err = file.Write(buffer)
	if l.Check(err) {
		return err
	}

	// truncate the file
	err = file.Truncate(int64(len(buffer)))
	if l.Check(err) {
		return err
	}

	// set the size
	return nil
}
