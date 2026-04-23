package radio

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetadataConcurrent(t *testing.T) {
	r := NewRadio("http://example.com/stream")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 100 {
			r.mu.Lock()
			r.metadata = Metadata{Artist: "A", Title: string(rune('0' + i%10))}
			r.mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for range 100 {
			m := r.GetMetadata()
			assert.NotNil(t, m)
		}
	}()

	wg.Wait()
}
