package smb

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Builder is the main client and contains all sitemap builder configuration options.
type Builder struct {
	host       string
	publicPath string
	filename   string
	compress   bool
	verbose    bool
}

// NewBuilder initializes a new Builder client.
func NewBuilder(host, path, filename string, compress, verbose bool) (*Builder, error) {
	if host == "" {
		return nil, errors.New("invalid host")
	}
	if path == "" {
		return nil, errors.New("invalid public path")
	}
	if filename == "" {
		return nil, errors.New("invalid filename")
	}

	return &Builder{
		host:       host,
		publicPath: path,
		filename:   filename,
		compress:   compress,
		verbose:    verbose,
	}, nil
}

// BuildFiles creates the sitemap index and all other files containing the urlset.
func (b *Builder) BuildFiles(urlset []*URL) error {
	b.println("start creating files process")

	filesLength, err := b.buildURLSet(urlset)
	if err != nil {
		return err
	}

	if err = b.buildIndex(filesLength); err != nil {
		return err
	}

	b.println("end creating files process")

	return nil
}

// PingSearchEngines sends a GET request to search engines ping URLs passing the sitemap index file.
func (b *Builder) PingSearchEngines(seURLs ...string) {
	seURLs = append(seURLs, []string{
		"http://www.google.com/webmasters/tools/ping?sitemap=%s",
		"http://www.bing.com/webmaster/ping.aspx?siteMap=%s",
	}...)
	filePath := fmt.Sprintf("%s/sitemaps/%s.xml", b.host, b.filename)
	if b.compress {
		filePath += ".gz"
	}

	buf := len(seURLs)
	ch := make(chan string, buf)
	httpClient := http.Client{Timeout: time.Duration(5 * time.Second)}
	for _, seURL := range seURLs {
		go func(URL string) {
			URL = fmt.Sprintf(URL, filePath)
			resp, err := httpClient.Get(URL)
			if err != nil {
				resp.Body.Close()
				ch <- fmt.Sprintf("ping failed for url: %s", URL)
				return
			}
			resp.Body.Close()
			ch <- fmt.Sprintf("ping succeeded: %s", URL)
		}(seURL)
	}

	for i := 0; i < buf; i++ {
		b.println(<-ch)
	}
}
