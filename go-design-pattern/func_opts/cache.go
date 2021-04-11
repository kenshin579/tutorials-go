package func_opts

import "time"

type Config struct {
	MaxAge  time.Duration
	MaxSize int
}

type Cache struct {
	Config
}

type CacheOption func(c *Cache) error

func NewCache(opts ...CacheOption) (*Cache, error) {
	c := &Cache{}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func CacheMaxAge(maxAge time.Duration) CacheOption {
	return func(c *Cache) error {
		// set max age
		c.MaxAge = maxAge
		return nil
	}
}

func CacheMaxEntries(maxSize int) CacheOption {
	return func(c *Cache) error {
		// set max entries
		c.MaxSize = maxSize
		return nil
	}
}
