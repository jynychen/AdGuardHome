package home

import (
	"fmt"
	"time"

	"github.com/AdguardTeam/dnsproxy/proxy"
	"github.com/AdguardTeam/golibs/stringutil"
	"github.com/jynychen/AdGuardHome/pkg/client"
	"github.com/jynychen/AdGuardHome/pkg/filtering"
	"github.com/jynychen/AdGuardHome/pkg/filtering/safesearch"
	"github.com/jynychen/AdGuardHome/pkg/whois"
)

// Client contains information about persistent clients.
type Client struct {
	// upstreamConfig is the custom upstream config for this client.  If
	// it's nil, it has not been initialized yet.  If it's non-nil and
	// empty, there are no valid upstreams.  If it's non-nil and non-empty,
	// these upstream must be used.
	upstreamConfig *proxy.UpstreamConfig

	safeSearchConf filtering.SafeSearchConfig
	SafeSearch     filtering.SafeSearch

	// BlockedServices is the configuration of blocked services of a client.
	BlockedServices *filtering.BlockedServices

	Name string

	IDs       []string
	Tags      []string
	Upstreams []string

	UseOwnSettings        bool
	FilteringEnabled      bool
	SafeBrowsingEnabled   bool
	ParentalEnabled       bool
	UseOwnBlockedServices bool
	IgnoreQueryLog        bool
	IgnoreStatistics      bool
}

// ShallowClone returns a deep copy of the client, except upstreamConfig,
// safeSearchConf, SafeSearch fields, because it's difficult to copy them.
func (c *Client) ShallowClone() (sh *Client) {
	clone := *c

	clone.BlockedServices = c.BlockedServices.Clone()
	clone.IDs = stringutil.CloneSlice(c.IDs)
	clone.Tags = stringutil.CloneSlice(c.Tags)
	clone.Upstreams = stringutil.CloneSlice(c.Upstreams)

	return &clone
}

// closeUpstreams closes the client-specific upstream config of c if any.
func (c *Client) closeUpstreams() (err error) {
	if c.upstreamConfig != nil {
		err = c.upstreamConfig.Close()
		if err != nil {
			return fmt.Errorf("closing upstreams of client %q: %w", c.Name, err)
		}
	}

	return nil
}

// setSafeSearch initializes and sets the safe search filter for this client.
func (c *Client) setSafeSearch(
	conf filtering.SafeSearchConfig,
	cacheSize uint,
	cacheTTL time.Duration,
) (err error) {
	ss, err := safesearch.NewDefault(conf, fmt.Sprintf("client %q", c.Name), cacheSize, cacheTTL)
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return err
	}

	c.SafeSearch = ss

	return nil
}

// RuntimeClient is a client information about which has been obtained using the
// source described in the Source field.
type RuntimeClient struct {
	// WHOIS is the filtered WHOIS data of a client.
	WHOIS *whois.Info

	// Host is the host name of a client.
	Host string

	// Source is the source from which the information about the client has
	// been obtained.
	Source client.Source
}
