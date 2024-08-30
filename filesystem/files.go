package filesystem

import (
	"bufio"
	"os"

	"github.com/khan-lau/kutils/container/klists"
)

func ReadFile(filePath string) (*klists.KList[string], error) {
	// 按行读文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineList := klists.New[string]()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineList.PushBack(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return lineList, nil
}
