package main

import (
	"fmt"

	"github.com/pkgng/structor"
)

type Human struct {
	Name  string
	Role  string
	Age   *int32
	Tel   string
	Notes []string
	Flags string
}

type AddressT struct {
	Address string
	Adcode  string
	Gps     string
}

type WifeT struct {
	structor.BaseStructor `structor:"CopyByName,Wife"`
	Name                  string
	NickName              string `structor:"Wife.Name.toLocaleLowerCase()"`
	Age                   int
	Age3                  int32 `structor:"self.Age + 3"`
}

type Farmer struct {
	structor.BaseStructor `structor:"CopyByName,Human,address"`
	Name                  string
	Age                   int64
	Nickname              string   `structor:"self.Name.toLocaleLowerCase()"`
	DoubleAge             int32    `structor:"Human.Age * 2"`
	SuperRole             string   `structor:"'Farmer-' + Human.Role"`
	Notes                 []string `structor:"Human.Notes.reverse()"`
	Flags                 []string `structor:"Human.Flags.split(',')"`
	Contact               struct {
		Tel     string `structor:"Human.Tel"`
		Address string `structor:"address.Address"`
		Adcode  string `structor:"address.Adcode"`
	}
	Wife WifeT
}

func main() {
	var age int32 = 23
	var age2 int32 = 22
	man := Human{Name: "LiLei", Age: &age, Tel: "18611009988", Role: "Farmer", Notes: []string{"hello", "world"}, Flags: "a,b,c"}
	address := AddressT{Adcode: "110108", Address: "北京海淀区五道口优盛大厦D座", Gps: "116.328115,40.054629"}
	wife := Human{Name: "HanMeiMei", Age: &age2, Role: "Wife", Notes: []string{"hello", "world"}, Flags: "e,f,g"}

	farmer := Farmer{}
	structor.New().Set("Human", &man).Set("address", address).Set("Wife", &wife).Construct(&farmer)

	fmt.Printf("%#v\n", farmer)
}
