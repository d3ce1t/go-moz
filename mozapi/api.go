package mozapi

import "fmt"

const (
	// WaitTimeBetweenRequests is the number of seconds to wait
	// between each requests in order to not overcome the quotas
	// of the free access
	WaitTimeBetweenRequests = 10

	// MaxRequestsPerSecond is the number of URLs that can be retrieved each second
	MaxRequestsPerSecond = 10

	// ExpireTimeInSeconds of the request
	ExpireTimeInSeconds = 300
)

// Columns to retrieve for a given URL
// https://moz.com/help/guides/moz-api/mozscape/api-reference/url-metrics
const (
	Title               int64 = 1
	CanonicalURL        int64 = 4
	ExternalEquityLinks int64 = 32
	Links               int64 = 2048
	MozRankForURL       int64 = 16384
	MozRankForSubdomain int64 = 32768
	HTTPStatusCode      int64 = 536870912
	PageAuthority       int64 = 34359738368
	DomainAuthority     int64 = 68719476736
	TimeLastCrawled     int64 = 144115188075855872
)

type MozAPI interface {
	MetricsForURL(url string, columns int64) (*URLMetrics, error)
	MetricsForURLBatch(urls []string, columns int64) ([]*URLMetrics, error)
}

// URLMetrics retrieved by Moz Free API
type URLMetrics struct {

	// The title of the page, if available
	Title string `json:"ut,omitempty"`

	// The canonical form of the URL
	CanonicalURL string `json:"uu,omitempty"`

	// The number of external equity links to the URL
	ExternalEquityLinks int `json:"ueid,omitempty"`

	// The number of links (equity or nonequity or not,
	// internal or external) to the URL
	Links int `json:"uid,omitempty"`

	// The MozRank of the URL, in both the normalized 10-point
	// score (umrp) and the raw score (umrr)
	MozRankForURLNormalized float32 `json:"umrp,omitempty"`
	MozRankForURLRaw        float32 `json:"umrr,omitempty"`

	// The MozRank of the URL's subdomain, in both the normalized
	// 10-point score (fmrp) and the raw score (fmrr)
	MozRankForSubdomainNormalized float32 `json:"fmrp,omitempty"`
	MozRankForSubdomainRaw        float32 `json:"fmrr,omitempty"`

	// The HTTP status code recorded by Mozscape for this URL, if available
	HTTPStatusCode int `json:"us,omitempty"`

	// A normalized 100-point score representing the likelihood of
	// a page to rank well in search engine results
	PageAuthority float32 `json:"upa,omitempty"`

	// A normalized 100-point score representing the likelihood of
	// a domain to rank well in search engine results
	DomainAuthority float32 `json:"pda,omitempty"`

	// The time and date on which Mozscape last crawled the URL,
	// returned in Unix epoch format
	TimeLastCrawled int64 `json:"ulc,omitempty"`
}

func (m URLMetrics) String() string {
	return fmt.Sprintf("URL: %v, PA: %v, DA: %v, MRU: %v, MRS: %v",
		m.CanonicalURL, m.PageAuthority, m.DomainAuthority,
		m.MozRankForURLNormalized, m.MozRankForSubdomainNormalized)
}
