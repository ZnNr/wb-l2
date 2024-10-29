package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

/*
L2.10 «Утилита wget»
Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

// Функция для сохранения содержимого в файл
func saveToFile(filename string, content []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0644)
}

// Функция для загрузки страницы
func downloadPage(url string, baseDir string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// Check if the response status is OK
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при загрузке страницы: %s", res.Status)
	}

	// Creating a file path from the URL
	parsedUrl := strings.TrimPrefix(url, "http://")
	parsedUrl = strings.TrimPrefix(parsedUrl, "https://")
	filePath := filepath.Join(baseDir, parsedUrl)
	filePath += "index.html" // Save as index.html for the main page

	content := make([]byte, 0)
	for {
		buffer := make([]byte, 4096)
		n, err := res.Body.Read(buffer)
		if n > 0 {
			content = append(content, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	if err := saveToFile(filePath, content); err != nil {
		return err
	}

	fmt.Printf("Скачано: %s\n", filePath)
	return nil
}

// Функция для извлечения ссылок из HTML
func extractLinks(url string, body []byte) []string {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

// Основная функция
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <URL>")
		return
	}

	url := os.Args[1]
	baseDir := "downloads"

	if err := downloadPage(url, baseDir); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Здесь необходимо добавить логику для извлечения и загрузки дополнительных ресурсов и страниц
}
