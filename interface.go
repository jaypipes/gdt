package gdt

type Runnable interface {
	RunResult()
}

type Fixture interface {
	Start()
	Stop()
}
