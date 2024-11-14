package utils

import (
	"fmt"
	"io"
	"multipart-upload/global"
	"multipart-upload/models"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MkMultiDir 创建多级文件夹
func MkMultiDir(file string) (string, error) {
	logPath := fmt.Sprintf("%s/%d/%d/%d/%s", global.Config.File.Path, time.Now().Year(), time.Now().Month(), time.Now().Day(), file)
	//判断目录是否存在
	_, err := os.Stat(logPath)
	//目录不存在就创建
	if err != nil {
		if mkdirAllErr := os.MkdirAll(logPath, os.ModePerm); mkdirAllErr != nil {
			return "", mkdirAllErr
		}
	}
	return logPath, nil
}

// IsExist 判断是否存在指定目录
func IsExist(targetName string) (bool, string) {
	var found bool = false
	var foundPath string = ""
	err := filepath.WalkDir(global.Config.File.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == targetName {
			found = true
			foundPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		fmt.Printf("遍历目录时发生错误: %v", err)
	}
	return found, foundPath
}

// 合并文件
func MergeFile(file models.FileModel) error {
	// 创建一个 map 来存储每个文件名及其对应的分片
	fileChunks := make(map[int]string)
	for _, f := range file.FileChunkModel {
		chunkIndex := f.Index
		fileChunks[chunkIndex] = f.FilePath
	}
	if err := MergeChunks(fileChunks, file); err != nil {
		fmt.Println("合并失败")
		return err
	}

	return nil
}

// 合并分片
func MergeChunks(chunks map[int]string, file models.FileModel) error {

	filePath := fmt.Sprintf("%s_%v%s", strings.Split(strings.TrimSuffix(file.FilePath, filepath.Ext(file.FilePath)), "_")[0], time.Now().UnixNano(), filepath.Ext(file.FileName))
	// filePath := fmt.Sprintf("%s_%v%s", strings.TrimSuffix(file.FilePath, filepath.Ext(file.FileName)), time.Now().UnixNano(), filepath.Ext(file.FileName))
	outFile, outFileErr := os.Create(filePath)
	if outFileErr != nil {
		fmt.Println("outFileErr==", outFileErr)
		return outFileErr
	}
	defer outFile.Close()
	for i := 1; i <= len(chunks); i++ {
		chunkFilePath := filepath.Join(chunks[i])
		chunkFile, chunkFileErr := os.Open(chunkFilePath)
		if chunkFileErr != nil {
			fmt.Println("chunkFileErr==", chunkFileErr)
			return chunkFileErr
		}
		defer chunkFile.Close()
		_, copyErr := io.Copy(outFile, chunkFile)
		if copyErr != nil {
			fmt.Println("copyErr==", copyErr)
			return copyErr
		}
	}

	return nil
}

// 清理临时文件
func CleanupTempFiles(filePath string) error {
	// 递归删除目录下的所有文件和子目录
	if err := RemoveDir(filePath); err != nil {
		fmt.Printf("无法删除目录内容: %v\n", err)
		return err
	}
	// 最后删除空目录
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("无法删除空目录: %v\n", err)
		return err
	} else {
		fmt.Println("目录及其所有内容已成功删除")
	}
	return nil
}

// RemoveDir 删除给定路径的目录及其下所有文件和子目录。
func RemoveDir(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 先处理文件和子目录，确保在删除父目录前所有内容都已删除
		if !info.IsDir() {
			// 如果不是目录，则删除文件
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}
