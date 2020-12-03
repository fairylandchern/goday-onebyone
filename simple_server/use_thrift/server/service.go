package main

import (
	"context"
	"goday-onebyone/simple_server/use_thrift/gen-go/compute"
	"math/big"
)

// computeThrift 实现service中定义的方法
type divmodThrift struct {
}

// 每个方法除了定义的返回值之外还要返回一个error，包括定义成void的方法。自定义类型会在名字之后加一条下划线// 暂时用不到context，所以忽略
func (d *divmodThrift) DoDivMod(_ context.Context, arg1, arg2 int64) (*compute.Result_, error) {
	divRes := int64(arg1 / arg2)
	modRes := int64(arg1 % arg2)
	// 生成的用于生成自定义数据对象的函数
	res := compute.NewResult_()
	res.Div = divRes
	res.Mod = modRes

	return res, nil
}

// 尽量一个struct对应一个service
type mulrangeThrift struct {
}

func (m *mulrangeThrift) BigRange(_ context.Context, max int64) (string, error) {
	result := new(big.Int)
	result.SetString("1", 10)
	result = result.MulRange(2, max)

	return result.String(), nil
}
