package grapeTimer

import (
	"fmt"
	"testing"
)

func CallBackOnce(arg1 string, arg2 int, arg3 float32) {
	fmt.Println(arg1, arg2, arg3)
}

func Test_CallReflect(t *testing.T) {
	cb, err := reflectFunc(CallBackOnce, "this arg1", 2000, float32(300.5))
	if err != nil {
		t.Error(err)
		return
	}

	callFunc(cb)
}

func Benchmark_Parallel(b *testing.B) {
	cb, err := reflectFunc(CallBackOnce, "this arg1", 2000, float32(300.5))
	if err != nil {
		b.Error(err)
		return
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			callFunc(cb)
		}
	})
}
