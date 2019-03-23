package fetcher

type serviceFetch struct{}

func (s *serviceFetch) Get() (string, error) {
	return "get", nil
}

func (s *serviceFetch) List() ([]string, error) {
	return []string{"get", "list"}, nil
}

var defaultFetcher Fetcher = newPoolFetcher(&serviceFetch{}, 2)

// NewFetcher - return singletone instance of fetcher
func NewFetcher() Fetcher {
	return defaultFetcher
}

// Fetcher - call external API
type Fetcher interface {
	Get() (string, error)
	List() ([]string, error)
}

type poolFetcher struct {
	f    Fetcher
	pool chan struct{}
}

var _ Fetcher = (*poolFetcher)(nil)

func newPoolFetcher(f Fetcher, workers int) Fetcher {
	return &poolFetcher{
		f:    f,
		pool: make(chan struct{}, workers),
	}
}

func (p *poolFetcher) Get() (string, error) {
	p.pool <- struct{}{}
	out, err := p.f.Get()
	<-p.pool
	return out, err
}

func (p *poolFetcher) List() ([]string, error) {
	p.pool <- struct{}{}
	out, err := p.f.List()
	<-p.pool
	return out, err
}
