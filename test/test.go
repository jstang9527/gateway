package main

import (
	"fmt"
)

// Node ...
type Node struct {
	Value int
	Next  *Node
}

// Create n就是0位,指向第一位
func (n *Node) Create(list []int) {
	cur := n
	for i := 0; i < len(list); i++ {
		//1.实例化node对象
		next := &Node{Value: list[i]}
		//2.当前指向性生成的对象
		cur.Next = next
		//3.当前对象更为最新
		cur = next
	}
}

// Delete 删除链表节点, 根据值找对应值的对象,并删除
func (n *Node) Delete(val int) {
	item := n
	for item.Next != nil {
		if item.Next.Value == val {
			item.Next = item.Next.Next
			return
		}
		item = item.Next
	}
}

func main() {
	obj := &Node{}
	obj.Create([]int{4, 2, 7, 1, 9, 3})
	obj.testResult()

	// 删除节点
	obj.Delete(5)
	obj.testResult()
}

// + 测试打印 ---------------------------------------------
func (n *Node) testResult() {
	item := n.Next
	fmt.Print(n.Value)
	for item != nil {
		fmt.Print("->", item.Value)
		item = item.Next
	}
	fmt.Println()
}
