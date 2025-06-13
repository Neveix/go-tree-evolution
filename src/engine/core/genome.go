package core

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Genome struct {
	data [64]byte
}

func GenomeCreate() *Genome {
	a := Genome{}
	for i := 0; i < 64; i++ {
		a.data[i] = byte(rand.Intn(32))
	}
	return &a
}

func (g *Genome) Mutate() {
	for k := 0; k < rand.Intn(10)+1; k++ {
		x := rand.Intn(64)
		g.data[x] = (g.data[x] + byte(rand.Intn(32))) % 32
	}
}

func (g *Genome) Mutated() *Genome {
	newgValue := *g
	newg := &newgValue
	newg.Mutate()
	return newg
}

func (g *Genome) ToString(malMode bool) string {
	var buf strings.Builder

	for i := 0; i < 64; i += 4 {
		if malMode {
			buf.WriteString(
				fmt.Sprintf(
					"%c: %c %c %c %c\n",
					ByteTo32mal(byte(i/4)),
					ByteTo32mal(g.data[i+0]), ByteTo32mal(g.data[i+1]),
					ByteTo32mal(g.data[i+2]), ByteTo32mal(g.data[i+3])))
		} else {
			buf.WriteString(
				fmt.Sprintf(
					"%d: %d %d %d %d\n",
					byte(i/4),
					g.data[i+0], g.data[i+1],
					g.data[i+2], g.data[i+3]))
		}
	}

	return buf.String()
}

func (g *Genome) Print() {
	fmt.Print(g.ToString(true))
}

func (g *Genome) GetGen(n byte) [4]byte {
	return [4]byte(g.data[n*4 : n*4+4])
}

func (g *Genome) Save() {
	os.WriteFile("genome.txt", []byte(g.ToString(false)), 0644)
}

func (g *Genome) Load() {
	data, _ := os.ReadFile("genome.txt")
	datastring := string(data)
	lines := strings.Split(datastring, "\n")[:16]
	for i, line := range lines {
		afterColon := strings.Split(line, ":")[1]
		// fmt.Printf("%d: ", i)
		gens := strings.Split(afterColon, " ")[1:]
		for j := 0; j < 4; j++ {
			trimmedGen := strings.TrimSpace(gens[j])
			intGen, _ := strconv.Atoi(trimmedGen)
			// fmt.Printf("%c(%d+%d) ", ByteTo32mal(byte(intGen)), i, j)
			g.data[4*i+j] = byte(intGen)
		}
		// fmt.Println()
	}
}
