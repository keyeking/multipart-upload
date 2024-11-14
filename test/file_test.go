package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	// FileName := "Redis-x64-3.2.100.zip"
	FilePath := "D:\\test\\2024\\10\\17\\Redis-x64-3.2.10011.zip"
	// SliceCount := 40
	// 获取目录下的所有文件
	files, err := os.ReadDir(FilePath)
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		return
	}
	// 创建一个 map 来存储每个文件名及其对应的分片
	fileChunks := make(map[int]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// 提取分片编号
		chunkIndex := strings.Split(file.Name(), "-")[0]
		// 将分片编号转换为整数
		index, _ := strconv.Atoi(chunkIndex)

		fileChunks[index] = file.Name()
	}
	fmt.Println(fileChunks)
	mergeChunks("123.zip", fileChunks, "D:\\test\\2024\\10\\17\\123.zip")
}

// 合并分片
func mergeChunks(filename string, chunks map[int]string, finalFilePath string) error {
	outFile, outFileErr := os.Create(finalFilePath)
	if outFileErr != nil {
		fmt.Println("outFileErr==", outFileErr)
		return outFileErr
	}
	defer outFile.Close()

	for i := 1; i <= len(chunks); i++ {
		chunkFilePath := filepath.Join("D:\\test\\2024\\10\\17\\Redis-x64-3.2.10011.zip", chunks[i])
		chunkFile, chunkFileErr := os.Open(chunkFilePath)
		if chunkFileErr != nil {
			fmt.Println("chunkFileErr==", chunkFileErr)
			return chunkFileErr
		}
		defer chunkFile.Close()
		b, copyErr := io.Copy(outFile, chunkFile)
		if copyErr != nil {
			fmt.Println("copyErr==", copyErr)
			return copyErr
		}
		fmt.Println(b, "b")
	}

	return nil
}

// 清理临时文件
func cleanupTempFiles(filePath string) error {
	// 递归删除目录下的所有文件和子目录
	if err := removeDir(filePath); err != nil {
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

// removeDir 删除给定路径的目录及其下所有文件和子目录。
func removeDir(path string) error {
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
