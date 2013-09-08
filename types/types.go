package types

type InputURLs struct {
	Inputs struct {
		URLs []string `json:"urls"`
	} `json:"inputs"`
}

type OutputLatencies struct {
	Outputs []*Stopwatch `json:"outputs"`
}

type Stopwatch struct {
	URL     string `json:"url"`
	Latency int64  `json:"latency"`  // Milliseconds
}
