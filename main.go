package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	Path        string
	Size        int64
	Hash        string
	ProcessedAt time.Time
}

type History struct {
}

func main() {
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Ошибка получения информации о файле: %v\n", err)
		return
	}
	fileSize := fileInfo.Size()

	var lineCount, wordCount int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		wordCount += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении: %v\n", err)
		return
	}

	fmt.Printf("Файл: %s\n", filepath)
	fmt.Printf("Размер: %d байт\n", fileSize)
	fmt.Printf("Количество строк: %d\n", lineCount)
	fmt.Printf("Количество слов: %d\n", wordCount)
}
