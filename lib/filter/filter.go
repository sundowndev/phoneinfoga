package filter

type Filter interface {
	Match(string) bool
}

type Engine struct {
	rules []string
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) AddRule(r ...string) {
	e.rules = append(e.rules, r...)
}

func (e *Engine) Match(r string) bool {
	for _, rule := range e.rules {
		if rule == r {
			return true
		}
	}
	return false
}
