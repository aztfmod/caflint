package lint

//Wrapper for OS.Exit to make testing easier.

type Exiter struct {
	statusCode int
	exiterFunc func(int)
}

func (e *Exiter) Exit(statusCode int) {
	e.statusCode = statusCode
	e.exiterFunc(statusCode)
}

func NewExiter(exiterFunc func(int)) *Exiter {
	exiter := new(Exiter)
	exiter.exiterFunc = exiterFunc
	return exiter
}
