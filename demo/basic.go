package main

import (
	"fmt"

	"github.com/pkgng/structor"
)

type Human struct {
	Name  string
	Role  string
	Age   *int32
	Notes []string
	Flags []byte
}

type Farmer struct {
	Name      string
	Age       int64
	Nickname  string   `structor:"self.Name.toLocaleLowerCase()"`
	DoubleAge int32    `structor:"Human.Age * 2"`
	SuperRule string   `structor:"'Farmer-' + Human.Role"`
	Notes     []string `structor:"Human.Notes.reverse()"`
	FlagCnt   int      `structor:"Human.Flags.length"`
}

func main() {
	var age int32 = 18
	man := Human{Name: "ZhangSan", Age: &age, Role: "Admin", Notes: []string{"hello", "world"}, Flags: []byte{'x', 'y', 'z'}}
	farmer := Farmer{}

	structor.NewStructor(&farmer).Set("Human", &man).CopyByName().Construct()

	fmt.Printf("%#v\n", farmer)
}
