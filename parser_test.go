package parser

import (
	"strings"
	"testing"

	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

func σ(a *G.Node) *G.Node {
	return G.Must(G.Sigmoid(a))
}

const lstm = `
xₜ∈ℝ⁶⁵
fₜ∈ℝ¹⁰⁰
iₜ∈ℝ¹⁰⁰
oₜ∈ℝ¹⁰⁰
hₜ∈ℝ¹⁰⁰
hₜ₋₁∈ℝ¹⁰⁰
cₜ₋₁∈ℝ¹⁰⁰
cₜ∈ℝ¹⁰⁰
xᵢ∈ℝ¹⁰⁰ˣ⁶⁵
Wᵢ∈ℝ¹⁰⁰ˣ⁶⁵
Uᵢ∈ℝ¹⁰⁰ˣ¹⁰⁰
Bᵢ∈ℝ¹⁰⁰
Wₒ∈ℝ¹⁰⁰ˣ⁶⁵
Uₒ∈ℝ¹⁰⁰ˣ¹⁰⁰
Bₒ∈ℝ¹⁰⁰
Wf∈ℝ¹⁰⁰ˣ⁶⁵
Uf∈ℝ¹⁰⁰ˣ¹⁰⁰
Bf∈ℝ¹⁰⁰
Wc∈ℝ¹⁰⁰ˣ⁶⁵
Uc∈ℝ¹⁰⁰ˣ¹⁰⁰
Bc∈ℝ¹⁰⁰
Wy∈ℝ⁶⁵ˣ¹⁰⁰
By∈ℝ⁶⁵

iₜ=(Wᵢ·xₜ+Uᵢ·hₜ₋₁+Bᵢ)
fₜ=σ(Wf·xₜ+Uf·hₜ₋₁+Bf)
oₜ=σ(Wₒ·xₜ+Uₒ·hₜ₋₁+Bₒ)
ĉₜ=tanh(Wc·xₜ+Uc·hₜ₋₁+Bc)
cₜ=fₜ*cₜ₋₁+iₜ*ĉₜ
hₜ=oₜ*tanh(cₜ)
y=Wy·hₜ+By
`

func TestProcess(t *testing.T) {
	g := G.NewGraph()
	p := NewParser(g)
	_, err := p.Process(strings.NewReader(lstm), nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParse(t *testing.T) {
	g := G.NewGraph()
	wfT := tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float32{1, 1, 1, 1}))
	wf := G.NewMatrix(g, tensor.Float32, G.WithName("wf"), G.WithShape(2, 2), G.WithValue(wfT))
	htprevT := tensor.New(tensor.WithBacking([]float32{1, 1}), tensor.WithShape(2))
	htprev := G.NewVector(g, tensor.Float32, G.WithName("ht-1"), G.WithShape(2), G.WithValue(htprevT))
	xtT := tensor.New(tensor.WithBacking([]float32{1, 1}), tensor.WithShape(2))
	xt := G.NewVector(g, tensor.Float32, G.WithName("xt"), G.WithShape(2), G.WithValue(xtT))
	bfT := tensor.New(tensor.WithBacking([]float32{1, 1}), tensor.WithShape(2))
	bf := G.NewVector(g, tensor.Float32, G.WithName("bf"), G.WithShape(2), G.WithValue(bfT))

	p := NewParser(g)
	p.Set(`Wf`, wf)
	p.Set(`hₜ₋₁`, htprev)
	p.Set(`xₜ`, xt)
	p.Set(`bf`, bf)
	//result, err := p.Parse(`σ(1*Wf·hₜ₋₁+ Wf·xₜ+ bf)`)
	type test struct {
		equation string
		expected []float32
	}
	for _, test := range []test{
		{
			`1*Wf·hₜ₋₁+ Wf·xₜ+ bf`,
			[]float32{5, 5},
		},
		{
			`σ(1*Wf·hₜ₋₁+ Wf·xₜ+ bf)`,
			[]float32{0.9933072, 0.9933072},
		},
		{
			`tanh(1*Wf·hₜ₋₁+ Wf·xₜ+ bf)`,
			[]float32{0.9999092, 0.9999092},
		},
	} {
		result, err := p.Parse(test.equation)
		if err != nil {
			t.Fatal(err)
		}
		machine := G.NewLispMachine(g, G.ExecuteFwdOnly())
		if err := machine.RunAll(); err != nil {
			t.Fatal(err)
		}
		res := result.Value().Data().([]float32)
		if len(res) != 2 {
			t.Fail()
		}
		if res[0] != test.expected[0] || res[1] != test.expected[1] {
			t.Fail()
		}
	}
}
