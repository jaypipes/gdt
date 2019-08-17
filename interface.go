package gdt

type Runnable interface {
	func Run() RunResult
}

type Fixture interface {
    func Start()
    func Stop()
}
