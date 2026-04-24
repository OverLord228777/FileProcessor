package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// FileInfo содержит идентификационные данные обработанного файла.
type FileInfo struct {
	Path        string
	Size        int64
	Hash        string
	ProcessedAt time.Time
}

// History хранит записи о ранее обработанных файлах.
type History struct {
	files  []FileInfo
	hashes map[string]bool
}

// Add добавляет запись в историю, если файла с таким хэшем ещё нет.
func (h *History) Add(f FileInfo) {
	if h.hashes == nil {
		h.hashes = make(map[string]bool)
	}
	if !h.Exists(f.Hash) {
		h.files = append(h.files, f)
		h.hashes[f.Hash] = true
	}
}

// Exists проверяет, присутствует ли в истории файл с заданным хэшем.
func (h *History) Exists(hash string) bool {
	if h.hashes == nil {
		return false
	}
	return h.hashes[hash]
}

// Processor описывает контракт для плагинов обработки содержимого.
type Processor interface {
	Process(fileData []byte) (string, error)
}

// LineCounter подсчитывает количество строк в данных.
type LineCounter struct{}

func (l LineCounter) Process(data []byte) (string, error) {
	count := bytes.Count(data, []byte("\n"))
	// Если последняя строка не заканчивается \n, она всё равно считается строкой
	if len(data) > 0 && data[len(data)-1] != '\n' {
		count++
	}
	return fmt.Sprintf("Number of lines: %d", count), nil
}

// WordCounter подсчитывает количество слов в данных.
type WordCounter struct{}

func (w WordCounter) Process(data []byte) (string, error) {
	words := strings.Fields(string(data))
	return fmt.Sprintf("Number of words: %d", len(words)), nil
}

// Checksummer вычисляет SHA-256 хэш содержимого.
type Checksummer struct{}

func (c Checksummer) Process(data []byte) (string, error) {
	hash := sha256.Sum256(data)
	return "SHA256: " + hex.EncodeToString(hash[:]), nil
}

// GzipCompressor сжимает данные и возвращает размер после сжатия.
type GzipCompressor struct{}

func (g GzipCompressor) Process(data []byte) (string, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return "", fmt.Errorf("gzip write error: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("gzip close error: %w", err)
	}
	return fmt.Sprintf("Compressed size: %d bytes", buf.Len()), nil
}

// Report генерирует отчёт о результате обработки файла.
type Report struct{}

func (r Report) Generate(fileName string, result string) {
	fmt.Printf("===== Report =====\n")
	fmt.Printf("File: %s\n", fileName)
	fmt.Printf("Result:\n%s\n", result)
	fmt.Printf("==================\n")
}

func main() {
	// Определяем флаг режима
	mode := flag.String("mode", "", "режим обработки: lines, words, hash, compress")
	flag.Parse()

	// Проверка, что режим указан
	if *mode == "" {
		fmt.Println("Укажите режим через -mode (lines, words, hash, compress)")
		flag.Usage()
		os.Exit(1)
	}

	// Получаем путь к файлу из не-флаговых аргументов
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Укажите путь к файлу для обработки")
		os.Exit(1)
	}
	filePath := args[0]

	// Читаем всё содержимое файла в память (для удобства передачи в процессор)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	// Вычисляем хэш всего содержимого для проверки истории
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// История (пока только в памяти на время выполнения)
	history := &History{
		files:  make([]FileInfo, 0),
		hashes: make(map[string]bool),
	}

	if history.Exists(hashStr) {
		fmt.Printf("Файл %s уже был обработан ранее (хэш: %s). Пропускаем.\n", filePath, hashStr)
		return
	}

	// Хранилище процессоров
	processors := map[string]Processor{
		"lines":    LineCounter{},
		"words":    WordCounter{},
		"hash":     Checksummer{},
		"compress": GzipCompressor{},
	}

	processor, ok := processors[*mode]
	if !ok {
		fmt.Printf("Неизвестный режим: %s. Доступны: lines, words, hash, compress\n", *mode)
		os.Exit(1)
	}

	// Выполняем обработку
	result, err := processor.Process(data)
	if err != nil {
		fmt.Printf("Ошибка обработки: %v\n", err)
		os.Exit(1)
	}

	// Добавляем файл в историю (запоминаем факт обработки)
	fileInfo, _ := os.Stat(filePath)
	history.Add(FileInfo{
		Path:        filePath,
		Size:        fileInfo.Size(),
		Hash:        hashStr,
		ProcessedAt: time.Now(),
	})

	// Генерируем отчёт
	var report Report
	report.Generate(filePath, result)
}
