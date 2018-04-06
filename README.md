[![](https://godoc.org/gorgonia.org/parser?status.svg)](http://godoc.org/gorgonia.org/parser)

#About

The goal of this project is to parse a mathematical formulae written in unicode and to transpile is into an `ExprGraph` of the Gorgonia project.

## Principle

This parser is generated from a `yacc` file from the subdirectory `src`.
If you want to contribute or add some new functionalities, you may need the `goyacc` tool and then run `go generate` from the `src` dubdirectory.
For a more complete explanation, you can refer to this [blog post](https://blog.owulveryck.info/2017/12/18/parsing-mathematical-equation-to-generate-computation-graphs---first-step-from-software-1.0-to-2.0-in-go.html)

## Available operations

As of today, the parser understands the following opreations on node objects (tensor based):

| Operation          | Gorgonia Operation                                                   | Symbol  | Unicode character |
|--------------------|----------------------------------------------------------------------|---------|-------------------|
| Multiplication     | [Mul](https://godoc.org/gorgonia.org/gorgonia#Mul)                   | ·       | U+00B7            |
| Hadamard Product   | [HadamardProd](https://godoc.org/gorgonia.org/gorgonia#HadamardProd) | *       |                   |
| Addition           | [Add](https://godoc.org/gorgonia.org/gorgonia#Add)                   | +       |                   |
| Substraction       | [Sub](https://godoc.org/gorgonia.org/gorgonia#Sub)                   | -       |                   |
| Pointwise Negation | [Neg]([Add](https://godoc.org/gorgonia.org/gorgonia#Neg)             | -       |                   |
| Sigmoid            | [Sigmiod](https://godoc.org/gorgonia.org/gorgonia#Sigmoid)           | σ       | U+03C3            |
| Tanh               | [Tanh](https://godoc.org/gorgonia.org/gorgonia#Tanh)                 | tanh    |                   |
| Softmax            | [Softmax](https://godoc.org/gorgonia.org/gorgonia#Softmax)           | softmax |                   |

## Example

```go
import (
	G "gorgonia.org/gorgonia"
	"gorgonia.org/parser"
	"gorgonia.org/tensor"
)

func main(){
	g := G.NewGraph()
	wfT := tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float32{1, 1, 1, 1}))
	wf := G.NewMatrix(g, tensor.Float32, G.WithName("wf"), G.WithShape(2, 2), G.WithValue(wfT))
	htprevT := tensor.New(tensor.WithBacking([]float32{1, 1}), tensor.WithShape(2))
	htprev := G.NewVector(g, tensor.Float32, G.WithName("ht-1"), G.WithShape(2), G.WithValue(htprevT))
	xtT := tensor.New(tensor.WithBacking([]float32{1, 1}), tensor.WithShape(2))
	xt := G.NewVector(g, tensor.Float32, G.WithName("xt"), G.WithShape(2), G.WithValue(xtT))
	bfT := tensor.New(tensor.WithBacking([]float32{1, 1}), tensor.WithShape(2))
	bf := G.NewVector(g, tensor.Float32, G.WithName("bf"), G.WithShape(2), G.WithValue(bfT))

	p := parser.NewParser(g)
	p.Set(`Wf`, wf)
	p.Set(`h`, htprev)
	p.Set(`x`, xt)
	p.Set(`bf`, bf)
        result, _ := p.Parse(`σ(1*Wf·h+ Wf·x+ bf)`)
        machine := G.NewLispMachine(g, G.ExecuteFwdOnly())
        if err := machine.RunAll(); err != nil {
              t.Fatal(err)
        }
        res := result.Value().Data().([]float32)
}
```

# Caution

* The parser is internally using a `map` and is not concurrent safe.
* The errors are not handle correctly
