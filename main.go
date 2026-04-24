package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: program <путь_к_файлу>")
		return
	}

	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	// Получаем размер файла (может понадобиться для FileInfo)
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Ошибка получения информации о файле: %v\n", err)
		return
	}
	fileSize := fileInfo.Size()

	// Вычисляем SHA-256 хэш содержимого
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Printf("Ошибка вычисления хэша: %v\n", err)
		return
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Возвращаем указатель чтения в начало файла
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("Ошибка перемещения по файлу: %v\n", err)
		return
	}

	// Инициализируем историю (в памяти на время выполнения)
	history := &History{
		files:  make([]FileInfo, 0),
		hashes: make(map[string]bool),
	}

	// Проверяем, обрабатывался ли файл ранее
	if history.Exists(hash) {
		fmt.Printf("Файл %s уже был обработан ранее (хэш: %s). Пропускаем.\n", filepath, hash)
		return
	}

	// Подсчёт строк и слов
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

	// Добавляем запись в историю
	history.Add(FileInfo{
		Path:        filepath,
		Size:        fileSize,
		Hash:        hash,
		ProcessedAt: time.Now(),
	})

	// Выводим статистику
	fmt.Printf("Файл: %s\n", filepath)
	fmt.Printf("Размер: %d байт\n", fileSize)
	fmt.Printf("Количество строк: %d\n", lineCount)
	fmt.Printf("Количество слов: %d\n", wordCount)
}
