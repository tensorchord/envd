package ir

import "fmt"

func (g Graph) CacheID(filename string) string {
	gpu := g.CUDA == nil || g.CUDNN == nil
	if gpu {
		return fmt.Sprintf("/%s-gpu-%s", g.CachePrefix, filename)
	}
	return fmt.Sprintf("/%s-cpu-%s", g.CachePrefix, filename)
}
