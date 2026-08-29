package serialize

import (
	"testing"

	"github.com/gongt/go/assets/serialize/pkg1"
	"github.com/gongt/go/assets/serialize/pkg2"
	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/serialize"
	"github.com/stretchr/testify/require"
)

func TestSomeTypeRoundTrip(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	src := &SomeType{
		Field1: pkg1.PublicType1{Field2: "hello"},
		Field2: &pkg2.PublicType2{},
		Field3: privateType{
			Field1: pkg1.PublicType1{Field2: "world"},
			Field2: &pkg2.PublicType2{},
		},
	}

	sCtx, err := serialize.NewSerializeContext("")
	require.NoError(t, err)
	data, err := src.Marshal(sCtx)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	dst := &SomeType{}

	dCtx, err := serialize.NewDeserializeContext()
	require.NoError(t, err)
	require.NoError(t, dst.Unmarshal(data, dCtx))

	require.Equal(t, src.Field1.Field2, dst.Field1.Field2)
	require.Equal(t, src.Field3.Field1.Field2, dst.Field3.Field1.Field2)
}

// func TestPublicType2ChanFuncID(t *testing.T) {
// 	myenv.RedirectDebugTesting(t)
//
// 	src := &pkg2.PublicType2{}
//
// 	sCtx, err := serialize.NewSerializeContext("")
// 	require.NoError(t, err)
// 	data, err := src.Marshal(sCtx)
// 	require.NoError(t, err)
//
// 	dCtx, err := serialize.NewDeserializeContext()
// 	require.NoError(t, err)
// 	dst := &pkg2.PublicType2{}
// 	require.NoError(t, dst.Unmarshal(data, dCtx))
// }
