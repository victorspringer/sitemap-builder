package smb

import "encoding/xml"

const (
	maxItemsLength      = 50000
	sitemapindexOpenTag = `<sitemapindex
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
		xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
		http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd">`
	sitemapindexCloseTag = "</sitemapindex>"
	urlsetOpenTag        = `<urlset
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
		xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"
		xmlns:video="http://www.google.com/schemas/sitemap-video/1.1"
		xmlns:geo="http://www.google.com/geo/schemas/sitemap/1.0"
		xmlns:news="http://www.google.com/schemas/sitemap-news/0.9"
		xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0"
		xmlns:pagemap="http://www.google.com/schemas/sitemap-pagemap/1.0"
		xmlns:xhtml="http://www.w3.org/1999/xhtml"
		xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
		http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">`
	urlsetCloseTag = "</urlset>"
)

type index struct {
	XMLName          xml.Name `xml:"sitemap"`
	Location         string   `xml:"loc"`
	LastModification string   `xml:"lastmod"`
}

// URL is the block model contained in the sitemap urlset.
type URL struct {
	XMLName          xml.Name `xml:"url"`
	Location         string   `xml:"loc"`
	Mobile           bool     `xml:"-"`
	ChangeFrequency  string   `xml:"changefreq"`
	Priority         float32  `xml:"priority"`
	LastModification string   `xml:"lastmod"`
}
