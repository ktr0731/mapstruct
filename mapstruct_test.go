package mapstruct

import (
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
		_, _, err := checkPrecondition(a, Foo{})
		require.Error(t, err)
	})

	t.Run("interface type", func(t *testing.T) {
		var a, b interface{}
		a = "foo"
		b = "bar"
		_, _, err := checkPrecondition(a, b)
		require.Error(t, err)
	})

	t.Run("not same types", func(t *testing.T) {
		_, _, err := checkPrecondition("foo", 0)
		require.Error(t, err)
	})

	t.Run("not struct", func(t *testing.T) {
		_, _, err := checkPrecondition("foo", "bar")
		require.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		_, _, err := checkPrecondition(Foo{}, Foo{})
		require.NoError(t, err)
	})
}

func TestMapStruct(t *testing.T) {
	t.Run("left has value", func(t *testing.T) {
		res, err := MapStruct(Foo{Hoge: "HOGE"}, Foo{})
		require.NoError(t, err)
		require.Exactly(t, Foo{Hoge: "HOGE"}, res.(Foo))
	})

	t.Run("right has value", func(t *testing.T) {
		res, err := MapStruct(Foo{}, Foo{Hoge: "HOGE"})
		require.NoError(t, err)
		require.Exactly(t, Foo{Hoge: "HOGE"}, res.(Foo))
	})

	t.Run("left overwritten by right value", func(t *testing.T) {
		res, err := MapStruct(Foo{Hoge: "dummy"}, Foo{Hoge: "HOGE"})
		require.NoError(t, err)
		require.Exactly(t, Foo{Hoge: "HOGE"}, res.(Foo))
	})
}
