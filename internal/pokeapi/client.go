package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nordluma/pokedexcli/internal/pokecache"
)

type Client struct {
	cache      internal.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheTtl time.Duration) Client {
	return Client{
		cache: internal.NewCache(cacheTtl),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Get(url string) ([]byte, error) {
	var areaResponse LocationAreaResponse

	// check cache
	entry, found := c.cache.Get(url)
	if found {
		if err := json.Unmarshal(entry, &areaResponse); err != nil {
			return nil, err
		}

		return entry, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok response: %d", res.StatusCode)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, data)

	return data, nil
}
