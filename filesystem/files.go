package filesystem

import (
	"bufio"
	"os"

	"github.com/khan-lau/kutils/container/klists"
)

func ReadFile(filePath string) (*klists.KList[string], error) {
	return ReadLinesWithBufferSize(filePath, 64*1024)
}

func ReadLinesDefaultBufferSize(filePath string) (*klists.KList[string], error) {
	// 按行读文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineList := klists.New[string]()
	scanner := bufio.NewScanner(file) // 默认缓冲区大小为64KB
	for scanner.Scan() {
		lineList.PushBack(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return lineList, nil
}

func ReadLines(filePath string) (*klists.KList[string], error) {
	// 按行读文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineList := klists.New[string]()

	reader := bufio.NewReaderSize(file, 10*1024*1024) //缓冲区10M
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineList.PushBack(string(line))
	}

	return lineList, nil
}

func ReadLinesWithBufferSize(filePath string, bufferSize int) (*klists.KList[string], error) {
	// 按行读文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineList := klists.New[string]()
	if bufferSize < 1 {
		bufferSize = 4 * 1024
	}

	reader := bufio.NewReaderSize(file, bufferSize) //缓冲区4k
	// for {
	// 	line, _, err := reader.ReadLine()
	// 	if err != nil {
	// 		break
	// 	}
	// 	lineList.PushBack(string(line))
	// }
	for {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) < 1 {
			break
		}
		// line = kstrings.TrimSpace(line)
		lineList.PushBack(line)
	}

	return lineList, nil
}
