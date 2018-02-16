# mapstruct

``` go
type Foo struct {
  Hoge string
  Fuga string
  piyo string
}

v1 := Foo{Hoge: "dummy", Fuga: "FUGA", piyo: "PIYO"}
v2 := Foo{Hoge: "HOGE"}
ires, _ := MapStruct(v1, v2)
res := ires.(Foo)

fmt.Println("%#v", res) // main.Foo{Hoge:"HOGE", Fuga:"FUGA", piyo:""}

```
