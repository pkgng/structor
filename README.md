# Structor

  I am a struct initializer, init struct from others 


## Features

* Copy from struct to struct with same name
* Calc field from JavaScript expression
* Calc field from multi Object

## Demo Usage

* source

```go
package main

import (
	"encoding/json"
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

	// fmt.Printf("%#v\n", farmer)

	b, err := json.Marshal(farmer)
	if err != nil {
		fmt.Println("JSON ERR:", err)
	}
	fmt.Println(string(b))
}
```

* output

```output
	{
		"Name":"LiLei",
		"Age":23,
		"Nickname":"lilei",
		"DoubleAge":46,
		"SuperRole":"Farmer-Farmer",
		"Notes":["world","hello"],
		"Flags":["a","b","c"],
		"Contact":{
			"Tel":"18611009988",
			"Address":"北京海淀区五道口优盛大厦D座",
			"Adcode":"110108"
		},
		"Wife":{
			"Name":"HanMeiMei",
			"NickName":"hanmeimei",
			"Age":22,
			"Age3":25
		}
	}
```

## 基础介绍

本工具用于对象的构建，可以指定相同名字字段的copy，

也可以使用其他对象的成员变量计算生成，tag内文本语法为javascrip表达式。

---

* 基类 structor.BaseStructor

基类的 structor Tag 为逗号分割的字符串，用于指定构建对象需要用到的对象名字；

成员字段计算Tag中可以直接使用这些对象的名字来引用对象；


内置功能字符串为：

|内置功能字符串 | 功能描述 |
| --- | --- |
| CopyByName | copy 相同名字的字段到目标对象上 |

demo 中 Farmer 的基类Tag 含义为： 

从 Human 和 address 两个对象来构建自己，并将有相同名字的字段 copy过来。

---

* 成员字段计算Tag

tag内文本语法为javascrip表达式，可以使用的变量为基类中指定的对象； 

另外 self变量指向本对象，可以使用self来引用CopyByName的字段，因为Copy逻辑在计算逻辑之前执行。

---

* Structor 对象 Set(name string, value interface{})

添加构建目标对象的资源, 需要指定名字，名字要与BaseStructor Tag 中的名字相同

---

* Structor  Construct(target interface{})

构建目标方法，调用后立即执行，输入参数为需要构建的对象的地址

---

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

# Author

**pkgng**

* <http://github.com/pkgng>


## License

Released under the [MIT License](https://github.com/pkgng/structor/blob/master/License).
