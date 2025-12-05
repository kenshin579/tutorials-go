package error_handling

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
)

func Test(t *testing.T) {
	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < 5; i++ {
		i := i
		g.Go(func() error {
			if err := printIndex(i); err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Error(err)
	}
}

func printIndex(n int) error {
	if n == 3 {
		return errors.New("invalid index")
	}

	fmt.Println("goroutine", n)
	return nil
}
