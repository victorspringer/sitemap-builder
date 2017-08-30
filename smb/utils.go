package smb

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

func (b *Builder) println(a ...interface{}) {
	if b.verbose {
		fmt.Printf("%s - %s\n", time.Now().Format(time.RFC3339), a)
	}
}

func (b *Builder) buildIndex(filesLength int) error {
	lastmod := time.Now().Format(time.RFC3339)
	sContent := ""

	for i := 1; i <= filesLength; i++ {
		loc := fmt.Sprintf("%s/sitemaps/%s%v.xml", b.host, b.filename, i)
		if b.compress {
			loc += ".gz"
		}

		smIndex := index{
			Location:         loc,
			LastModification: lastmod,
		}

		b, err := xml.Marshal(smIndex)
		if err != nil {
			return err
		}

		sContent += string(b)
	}

	content := []byte(xml.Header + sitemapindexOpenTag + sContent + sitemapindexCloseTag)
	filePath := fmt.Sprintf("%s/%s.xml", b.publicPath, b.filename)
	if err := b.createFile(content, filePath); err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildURLSet(urlset []*URL) (int, error) {
	sContent := ""
	urlsetLength := len(urlset)
	totalIndex := int(math.Ceil(float64(urlsetLength) / maxItemsLength))
	urlIndex := 0
	limit := 0

	for index := 1; index <= totalIndex; index++ {
		if index*maxItemsLength > urlsetLength {
			urlIndex = limit
			limit = urlsetLength
		} else {
			limit = maxItemsLength * index
			urlIndex = limit - maxItemsLength
		}
		for _, url := range urlset[urlIndex:limit] {
			b, err := xml.Marshal(&url)
			if err != nil {
				return 0, err
			}
			toString := string(b)
			if url.Mobile {
				toString = strings.Replace(toString, "</loc>", "</loc><mobile:mobile/>", 1)
			}
			sContent += toString
		}
		content := []byte(xml.Header + urlsetOpenTag + sContent + urlsetCloseTag)
		sContent = ""

		filePath := fmt.Sprintf("%s/%s%v.xml", b.publicPath, b.filename, index)
		if err := b.createFile(content, filePath); err != nil {
			return 0, err
		}
	}

	return totalIndex, nil
}

func (b *Builder) createFile(content []byte, path string) error {
	if err := ioutil.WriteFile(path, content, 0644); err != nil {
		return err
	}

	b.println("created file: " + path)

	if b.compress {
		if err := b.compressFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) compressFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	info, _ := file.Stat()
	size := info.Size()
	fileBytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(fileBytes)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(fileBytes)
	writer.Close()

	err = ioutil.WriteFile(path+".gz", buf.Bytes(), info.Mode())
	if err != nil {
		return err
	}

	file.Close()
	err = os.Remove(path)
	if err != nil {
		return err
	}

	b.println("compressed ", path)

	return nil
}
