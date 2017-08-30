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
	"time"
)

func (b *Builder) println(a ...interface{}) {
	if b.verbose {
		fmt.Printf("%s - %s\n", time.Now().Format(time.RFC3339), a)
	}
}

func (b *Builder) buildIndex(filesLength int) error {
	lastmod := time.Now().Format(time.RFC3339)
	filePath := fmt.Sprintf("%s/%s.xml", b.publicPath, b.filename)
	if err := b.createFile(xml.Header+sitemapindexOpenTag, filePath); err != nil {
		return err
	}

	for i := 1; i <= filesLength; i++ {
		loc := fmt.Sprintf("%s/sitemaps/%s%v.xml", b.host, b.filename, i)
		if b.compress {
			loc += ".gz"
		}

		smIndex := index{
			Location:         loc,
			LastModification: lastmod,
		}

		bt, err := xml.Marshal(smIndex)
		if err != nil {
			return err
		}

		if err = writeContent(&bt, filePath); err != nil {
			os.Remove(filePath)
			return err
		}
	}

	if err := b.closeFile(sitemapindexCloseTag, filePath); err != nil {
		os.Remove(filePath)
		return err
	}

	return nil
}

func (b *Builder) buildURLSet(urlset []*URL) (int, error) {
	urlsetLength := len(urlset)
	totalIndex := int(math.Ceil(float64(urlsetLength) / maxItemsLength))
	limit := maxItemsLength

	for index := 1; index <= totalIndex; index++ {
		filePath := fmt.Sprintf("%s/%s%v.xml", b.publicPath, b.filename, index)
		if err := b.createFile(xml.Header+urlsetOpenTag, filePath); err != nil {
			return 0, err
		}

		if index == totalIndex {
			limit = urlsetLength - ((index - 1) * maxItemsLength)
		}
		for _, url := range urlset[:limit] {
			var v interface{}
			if url.Mobile {
				v = buildURLMobile(url)
			} else {
				v = &url
			}
			bt, err := xml.Marshal(v)
			if err != nil {
				return 0, err
			}

			if err = writeContent(&bt, filePath); err != nil {
				os.Remove(filePath)
				return 0, err
			}
		}

		if err := b.closeFile(urlsetCloseTag, filePath); err != nil {
			os.Remove(filePath)
			return 0, err
		}
		urlset = removeFromURLSet(urlset, limit)
	}

	return totalIndex, nil
}

func (b *Builder) createFile(content, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	b.println("creating file: " + path)

	return nil
}

func (b *Builder) closeFile(content, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	file.Close()

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
	os.Remove(path)

	b.println("compressed ", path)

	return nil
}

func writeContent(content *[]byte, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(*content)
	if err != nil {
		return err
	}

	content = nil

	return nil
}

func removeFromURLSet(urlset []*URL, n int) []*URL {
	copy(urlset[0:], urlset[n:])
	for i, j := len(urlset)-n, len(urlset); i < j; i++ {
		urlset[i] = nil
	}
	return urlset[:len(urlset)-n]
}

func buildURLMobile(url *URL) interface{} {
	type urlMobile struct {
		XMLName          xml.Name `xml:"url"`
		Location         string   `xml:"loc"`
		Mobile           string   `xml:"mobile:mobile,allowempty"`
		ChangeFrequency  string   `xml:"changefreq"`
		Priority         float32  `xml:"priority"`
		LastModification string   `xml:"lastmod"`
	}

	return &urlMobile{
		Location:         url.Location,
		Mobile:           "",
		ChangeFrequency:  url.ChangeFrequency,
		Priority:         url.Priority,
		LastModification: url.LastModification,
	}
}
