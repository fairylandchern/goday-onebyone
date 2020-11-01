package sklist

import "math/rand"

type Node struct {
	key      interface{}
	level    int64
	score    int64
	forwards []*Node
}

type SL struct {
	head   *Node
	level  int64
	length int64
}

var (
	maxlevel = int64(16)
)

func newSkipListNode(v interface{}, score int64, level int64) *Node {
	return &Node{
		key: v,
		level: level,
		score: score,
		forwards: make([]*Node,level,level),
	}
}

// 返回int，1：key为空，2：key存在，0：插入成功
func (l *SL) Insert(key interface{}, score int64) int {
	if key == nil {
		return 1
	}

	// 主要比较关注点相关数据，记录每层位置数据
	cur := l.head
	update := make([]*Node, maxlevel)
	i := maxlevel - 1
	for ; i >= 0; i-- {
		for cur.forwards[i] != nil {
			if cur.forwards[i].key == key {
				return 2
			}

			if cur.forwards[i].score > score {
				update[i] = cur
				break
			}
			cur = cur.forwards[i]
		}
		if cur.forwards[i] == nil {
			update[i] = cur
		}
	}
	// 随机算法获取层数信息
	level := int64(1)
	for i := int64(0); i < maxlevel; i++ {
		if rand.Int31()%7 == 1 {
			level++
		}
	}
	// 生成新的链表节点信息
	//newNode := &Node{score: score, key: key, level: level, forwards: make([]*Node, level, level)}
	newNode:=newSkipListNode(key,score,level)
	// 原有节点连接
	for i := int64(0); i < level; i++ {
		next := update[i].forwards[i]
		update[i].forwards[i] = newNode
		newNode.forwards[i] = next
	}

	// 更新跳表层数
	if level > l.level {
		l.level = level
	}

	// 更新跳表长度
	l.length++

	return 0
}

// search data
func (l *SL) Find(v interface{}, score int64) *Node {
	if v == nil || l.length == 0 {
		return nil
	}

	// 遍历节点，查找符合条件的节点信息，并返回
	cur := l.head
	for i := l.level - 1; i >= 0; i-- {
		for cur.forwards[i] != nil {
			if cur.forwards[i].key == v && cur.forwards[i].score == score {
				return cur.forwards[i]
			}
			if cur.forwards[i].score > score {
				break
			}
			cur = cur.forwards[i]
		}
	}

	return nil
}

// delete data
func (l *SL) Delete(v interface{}, score int64) int {
	if v == nil {
		return 1
	}

	// 查找将删除节点的前驱节点
	update := make([]*Node, maxlevel, maxlevel)
	cur := l.head
	for i := cur.level - 1; i >= 0; i-- {
		update[i] = l.head
		for cur.forwards[i] != nil {
			if cur.forwards[i].key==v &&cur.score==score{
				update[i]=cur
				break
			}
			cur=cur.forwards[i]
		}
	}

	cur=update[0].forwards[0]
	// 查找并删除节点相关信息
	for i:=cur.level-1;i>=0;i-- {
		if update[i]==l.head&&cur.forwards[i]==nil{
			l.level=i
		}
		if update[i].forwards[i]!=nil{
			update[i].forwards[i]=update[i].forwards[i].forwards[i]
		}
	}
	l.length--
	return 0
}
