package main

import "runtime"

// from https://github.com/hellofresh/health-go/blob/b37d36e420e29217c2d25614268b00b655bb318c/health.go#L57-L69
type SystemMetrics struct {
	// Version is the go version.
	Version string `json:"version"`
	// GoroutinesCount is the number of the current goroutines.
	GoroutinesCount int `json:"goroutines_count"`
	// TotalAllocBytes is the total bytes allocated.
	TotalAllocBytes int `json:"total_alloc_bytes"`
	// HeapObjectsCount is the number of objects in the go heap.
	HeapObjectsCount int `json:"heap_objects_count"`
	// TotalAllocBytes is the bytes allocated and not yet freed.
	AllocBytes int `json:"alloc_bytes"`
}

// from https://github.com/hellofresh/health-go/blob/b37d36e420e29217c2d25614268b00b655bb318c/health.go#L243-L254
func newSystemMetrics() SystemMetrics {
	s := runtime.MemStats{}
	runtime.ReadMemStats(&s)

	return SystemMetrics{
		Version:          runtime.Version(),
		GoroutinesCount:  runtime.NumGoroutine(),
		TotalAllocBytes:  int(s.TotalAlloc),
		HeapObjectsCount: int(s.HeapObjects),
		AllocBytes:       int(s.Alloc),
	}
}
