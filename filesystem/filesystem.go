package filesystem

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// 获取程序执行路径
func GetExecPath() (string, error) {
	flag := false
	err := error(nil)
	file := ""
	if file, err = exec.LookPath(os.Args[0]); err == nil {
		flag = true
	} else {
		if _, err = os.Stat(os.Args[0]); err == nil {
			flag = true
		}
	}
	if flag {
		if file, err = filepath.Abs(file); err == nil {
			return strings.Replace(file, "\\", "/", -1), err
		}
	}
	return "", err
}

// 获取当前路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// 判断文件是否存在
func IsFileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func RemoveFile(path string) bool {
	return os.Remove(path) == nil
}

func RenameFile(path string, newPath string) bool {
	return os.Rename(path, newPath) == nil
}

// 修改文件扩展名
//
// @param filePath string 文件原路径
//
// @param newExtName string 新扩展名, 例如 ".exe", ".txt"
//
// @return bool
func FileReExtName(filePath string, newExtName string) bool {
	fileInfo, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	if fileInfo.IsDir() {
		return false
	}

	fileName := fileInfo.Name()
	// fileName := filePath
	extName := filepath.Ext(fileName)
	result := ""
	index := strings.LastIndex(filePath, extName)
	if index != -1 {
		result = filePath[:index] + newExtName
	} else {
		if strings.HasSuffix(filePath, ".") {
			result = filePath + strings.Replace(newExtName, ".", "", 1)
		} else {
			result = filePath + newExtName
		}
	}
	return os.Rename(filePath, result) == nil
}

// 修改文件扩展名, 只修改字符串, 不实际修改磁盘文件
//
// @param filePath string 文件原路径
//
// @param newExtName string 新扩展名, 例如 ".exe", ".txt"
//
// @return string
func ReExtName(filePath string, newExtName string) string {

	fileName := filepath.Base(filePath)
	extName := filepath.Ext(fileName)

	result := ""
	index := strings.LastIndex(filePath, extName)
	if index != -1 {
		result = filePath[:index] + newExtName
	} else {
		if strings.HasSuffix(filePath, ".") {
			result = filePath + strings.Replace(newExtName, ".", "", 1)
		} else {
			result = filePath + newExtName
		}
	}
	return result
}

// 获取文件名, 不包含文件扩展名
//
// @param filePath string 文件路径
//
// @return string
func FileMainName(filePath string) string {
	fileName := filepath.Base(filePath)
	extName := filepath.Ext(fileName)

	result := fileName
	if extName != "" {
		result = strings.Replace(fileName, extName, "", 1)
	} else {
		if strings.HasSuffix(filePath, ".") {
			result = strings.Replace(fileName, ".", "", 1)
		}
	}
	return result
}

// 从目录获取排序好的全部文件, 不包含目录, 不遍历子目录?
//
// @parameter path 目录路径
//
// @parameter prefix 过滤条件, 以`prefix`开始, 空字符串表示忽略该条件
//
// @parameter suffix 过滤条件, 以`suffix`结束, 空字符串表示忽略该条件
//
// @return []os.DirEntry 文件信息切片, 如失败则len()为0
//
// @return error 错误信息, 如正常则为nil
func HasPrefixWithSuffixFilesFromDir(path, prefix string, suffix string) ([]os.DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs, err := f.ReadDir(-1)
	if err == nil {
		// 排除目录
		index := int(0)
		files := make([]os.DirEntry, len(dirs))
		for _, item := range dirs {
			if item.IsDir() { // 跳过目录
				continue
			} else {
				// 与正则表达式是否匹配
				if strings.HasPrefix(item.Name(), prefix) && strings.HasSuffix(item.Name(), suffix) {
					files[index] = item
					index++
				}
			}
		}
		result := files[:index]
		sort.Slice(result, func(i, j int) bool { return result[i].Name() < result[j].Name() })
		return result, err
	} else {
		return dirs, err
	}
}
