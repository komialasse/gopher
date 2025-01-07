#!/bin/bash

go test -v clustering-index_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
go test -v count-index_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
go test -v points-index_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
go test -v counters_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
go test -v point_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
go test -v queue_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
go test -v sets_test.go ./point.go ./sets.go ./clustering-index.go ./geo-index.go ./queue.go ./points-index.go ./count-index.go ./counters.go ./testing.go ./benchmarks.go
