package main

import (
	sci "github.com/samuell/scipipe"
)

func main() {
	sci.InitLogInfo()

	foo := NewFooer()
	f2b := NewFoo2Barer()
	snk := sci.NewSink()

	f2b.InFoo = foo.OutFoo
	snk.In = f2b.OutBar

	pl := sci.NewPipeline()
	pl.AddProcs(foo, f2b, snk)
	pl.Run()
}

// ------------------------------------------------------------------------
// Components
// ------------------------------------------------------------------------

// Fooer

type Fooer struct {
	InnerProc *sci.SciProcess
	OutFoo    chan *sci.FileTarget
}

func NewFooer() *Fooer {
	innerFoo := sci.Shell("echo foo > {o:foo}")
	innerFoo.SetPathFormatterString("foo", "foo.txt")
	return &Fooer{
		InnerProc: innerFoo,
		OutFoo:    innerFoo.OutPorts["foo"],
	}
}

func (p *Fooer) Run() {
	p.InnerProc.OutPorts["foo"] = p.OutFoo
	p.InnerProc.Run()
}

// Foo2Barer

type Foo2Barer struct {
	InnerProc *sci.SciProcess
	InFoo     chan *sci.FileTarget
	OutBar    chan *sci.FileTarget
}

func NewFoo2Barer() *Foo2Barer {
	innerFoo2Bar := sci.Shell("sed 's/foo/bar/g' {i:foo} > {o:bar}")
	innerFoo2Bar.SetPathFormatterExtend("bar", "foo", ".bar.txt")
	return &Foo2Barer{
		InnerProc: innerFoo2Bar,
		InFoo:     innerFoo2Bar.InPorts["foo"],
		OutBar:    innerFoo2Bar.OutPorts["bar"],
	}
}

func (p *Foo2Barer) Run() {
	p.InnerProc.InPorts["foo"] = p.InFoo
	p.InnerProc.OutPorts["bar"] = p.OutBar
	p.InnerProc.Run()
}
