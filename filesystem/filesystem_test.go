package filesystem

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test_FileMainName(t *testing.T) {
	filePath := "D:\\Project\\Golang\\XianYi\\RedisSubscribe\\go.sum"
	fileName := filepath.Base(filePath)
	extName := filepath.Ext(fileName)

	result := FileMainName(filePath)

	fmt.Printf("mainName:%s, ext:%s\n", fileName, extName)
	fmt.Printf("result:%s\n", result)
}

func Test_ReExtName(t *testing.T) {
	filePath := "D:\\Project\\Golang\\XianYi\\RedisSubscribe\\go.sum"
	fileName := filepath.Base(filePath)
	extName := filepath.Ext(fileName)

	result := ReExtName(filePath, ".txt")

	fmt.Printf("fileName:%s, ext:%s\n", fileName, extName)
	fmt.Printf("result:%s\n", result)
}

func Test_HasPrefixWithSuffixFilesFromDir(t *testing.T) {
	srcFiles, err := HasPrefixWithSuffixFilesFromDir("C:/Private/Test/send", "", ".txt")
	if nil != err {
		t.Errorf("san dir error: %s", err.Error())
	}

	for idx, file := range srcFiles {
		fmt.Printf("%2.2d %s\n", idx, file.Name())
	}
}
