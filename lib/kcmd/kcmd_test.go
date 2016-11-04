
package kcmd

import (
	"../ginput"
	"testing"
	"fmt"
)

func Test1(t *testing.T) {
	cmd1 := []ginput.InputBits{
		ginput.Down,
		ginput.Down | ginput.Right,
		ginput.Right,
		ginput.Key1,
	}

	x := New()

	x.Update(ginput.Down)
	x.Update(ginput.Up)
	x.Update(ginput.Up)
	x.Update(ginput.Down | ginput.Right)
	x.Update(ginput.Right)
	x.Update(ginput.Key1)
	/*
	x.Update(ginput.Down | ginput.Right)
	x.Update(ginput.Right)
	x.Update(ginput.Key1)
	*/

	fmt.Println("x.Check()", x.Check(cmd1))
	fmt.Println("x.Buffer", x.Buffer)
}

func Test2(t *testing.T) {
	cmd1 := []ginput.InputBits{
		ginput.Nkey1,
		ginput.Key1,
	}

	x := New()

	/*
	x.Update(ginput.Down)
	x.Update(ginput.Up)
	x.Update(ginput.Up)
	x.Update(ginput.Down | ginput.Right)
	x.Update(ginput.Right)
	*/
	x.Update(ginput.Key1)
	x.Update(ginput.Key1)
	x.Update(ginput.Key1)
	x.Update(ginput.Key1)
	/*
	x.Update(ginput.Down | ginput.Right)
	x.Update(ginput.Right)
	x.Update(ginput.Key1)
	*/

	fmt.Println("x.Check()", x.Check(cmd1))
	fmt.Println("x.Buffer", x.Buffer)
}
