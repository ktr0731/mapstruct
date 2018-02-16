package mapstruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	Hoge string
	Fuga int
	piyo rune
}

func Test_checkPrecondition(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		var a interface{}
		_, err := checkPrecondition(a, Foo{})
		require.Error(t, err)
	})

	t.Run("interface type", func(t *testing.T) {
		var a, b interface{}
		a = "foo"
		b = "bar"
		_, err := checkPrecondition(a, b)
		require.Error(t, err)
	})

	t.Run("not same types", func(t *testing.T) {
		_, err := checkPrecondition("foo", 0)
		require.Error(t, err)
	})

	t.Run("not struct", func(t *testing.T) {
		_, err := checkPrecondition("foo", "bar")
		require.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		_, err := checkPrecondition(Foo{}, Foo{})
		require.NoError(t, err)
	})
}

func Test_obtainConcrete(t *testing.T) {
	t.Run("non pointer value", func(t *testing.T) {
		v := obtainConcrete(reflect.ValueOf("foo"))
		require.Equal(t, "foo", v.Interface().(string))
	})

	t.Run("pointer value", func(t *testing.T) {
		foo := "foo"
		v := obtainConcrete(reflect.ValueOf(&foo))
		require.Equal(t, foo, v.Interface().(string))
	})

	t.Run("nested pointer value", func(t *testing.T) {
		foo := "foo"
		foo2 := &foo
		v := obtainConcrete(reflect.ValueOf(&foo2))
		require.Equal(t, foo, v.Interface().(string))
	})
}

func TestMap(t *testing.T) {
	t.Run("left has value", func(t *testing.T) {
		res, err := Map(Foo{Hoge: "HOGE"}, Foo{})
		require.NoError(t, err)
		require.Exactly(t, Foo{Hoge: "HOGE"}, res.(Foo))
	})

	t.Run("right has value", func(t *testing.T) {
		res, err := Map(Foo{}, Foo{Hoge: "HOGE"})
		require.NoError(t, err)
		require.Exactly(t, Foo{Hoge: "HOGE"}, res.(Foo))
	})

	t.Run("left overwritten by right value", func(t *testing.T) {
		res, err := Map(Foo{Hoge: "dummy"}, Foo{Hoge: "HOGE"})
		require.NoError(t, err)
		require.Exactly(t, Foo{Hoge: "HOGE"}, res.(Foo))
	})

	t.Run("pointers", func(t *testing.T) {
		res, err := Map(&Foo{Hoge: "dummy"}, &Foo{Hoge: "HOGE"})
		require.NoError(t, err)
		require.Exactly(t, &Foo{Hoge: "HOGE"}, res.(*Foo))
	})

	t.Run("nested", func(t *testing.T) {
		type Bar struct {
			Foo Foo
		}
		res, err := Map(&Bar{Foo{Hoge: "dummy"}}, &Bar{Foo{Hoge: "HOGE"}})
		require.NoError(t, err)
		require.Exactly(t, &Bar{Foo{Hoge: "HOGE"}}, res.(*Bar))
	})

	t.Run("nested2", func(t *testing.T) {
		type Bar struct {
			Foo *Foo
		}
		res, err := Map(&Bar{&Foo{Hoge: "dummy"}}, &Bar{&Foo{Hoge: "HOGE"}})
		require.NoError(t, err)
		require.Exactly(t, &Bar{&Foo{Hoge: "HOGE"}}, res.(*Bar))
	})
}
