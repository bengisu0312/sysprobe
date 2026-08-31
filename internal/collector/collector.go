package collector

// Metric toplanan her bir verinin standart yapısını temsil eder.
type Metric struct {
	Name  string
	Value float64
	Unit  string
}

// Collector Name ve Collect metodlarına sahip olan her yapının (struct) 
// otomatik olarak bir collector olmasını sağlayan davranış tarifidir.
type Collector interface {
	Name() string
	Collect() ([]Metric, error)
}
func DefaultCollectors() []Collector {
	return []Collector{
		CPUCollector{},
		MemoryCollector{},
		DiskCollector{Path: "/"},
		LoadCollector{},
	}
}
