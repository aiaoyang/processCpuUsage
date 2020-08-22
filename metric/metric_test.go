package metric

// func Test_metric(t *testing.T) {
// 	m := NewCustomMetric(nil)
// 	m.Insert("a", "b")
// 	m.Send()
// }

/*
	Running tool: /usr/local/go/bin/go test -benchmem -run=^$ github.com/aiaoyang/processCpuUsage/metric -bench ^(Benchmark_copy)$ -v
		goos: linux
		goarch: amd64
		pkg: github.com/aiaoyang/processCpuUsage/metric
		Benchmark_copy
		Benchmark_copy-24    	  158278	      6754 ns/op	    2149 B/op	       4 allocs/op
*/
// func Benchmark_copy(b *testing.B) {
// 	m := NewCustomMetric(nil)
// 	tmp := make(map[string]interface{})
// 	tmp["1"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["2"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["3"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["4"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["5"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["6"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["7"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["8"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["a"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["b"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["c"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["d"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["e"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["f"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["g"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	tmp["h"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	m.Add(tmp)
// 	for i := 0; i < b.N; i++ {
// 		m.Copy()
// 	}
// }
