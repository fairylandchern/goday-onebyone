package sklist

import "testing"

// 跳表练习
func TestSl1(t *testing.T) {
	sl := &SL{
		head:   &Node{key: 0, level: maxlevel, score: 0, forwards: make([]*Node, maxlevel, maxlevel)},
		level:  1,
		length: 0,
	}
	// 遍历头节点,预先占了15级节点
	//for k := range sl.head.forwards {
	//	t.Log("now index:", k)
	//}
	t.Log(sl.Insert(3,888))
	n:=sl.Find(3,888)
	t.Log("node info:",n," sl:",sl)
	for _,v:=range sl.head.forwards{
		t.Log("after insert forwards info:",v)
	}
	t.Log("delete:",sl.Delete(3,888)," sl:",sl)
	for _,v:=range sl.head.forwards{
		t.Log("after delete forwards info:",v)
	}
}

