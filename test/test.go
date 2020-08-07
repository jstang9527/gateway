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
func Create(list []int) (head *Node) {
	if len(list) == 0 {
		return nil
	}
	head = &Node{Value: list[0]}
	cur := head //0.cur存储head的地址
	for i := 1; i < len(list); i++ {
		//1.实例化node对象
		next := &Node{Value: list[i]}
		//2.当前指向性生成的对象
		cur.Next = next
		//3.当前对象更为最新(cur存储下一个对象)
		cur = next
	}
	return
}

// Delete 删除链表节点, 根据值找对应值的对象,并删除
func (head *Node) Delete(val int) {
	item := head
	for item.Next != nil {
		if item.Next.Value == val {
			item.Next = item.Next.Next
			return
		}
		item = item.Next
	}
}

// Deserialization 反序列化
// 思路:
// 1. 第一位的下一位指向空。
// 2.其余位的下一位指向上一位。
func (head *Node) Deserialization() *Node {
	cur := head   //传进来的是链表头对象
	var pre *Node // 上一位起初是nil
	for cur != nil {
		// 将当前对象的下一位存起来
		next := cur.Next
		// 将当前对象的下一位指向上一位
		cur.Next = pre
		// 将上一位变为当前位
		pre = cur
		// 将当前位变为下一位
		cur = next
	}
	return pre //最后一位就是链表头啦
}

func main() {
	obj := Create([]int{4, 2, 7, 5, 9, 1, 3})
	obj.testResult()

	//删除节点
	obj.Delete(5)
	obj.testResult()

	//序列化节点
	obj = obj.Deserialization()
	obj.testResult()
}

// + testResult 测试打印 ---------------------------------------------
func (head *Node) testResult() {
	if head == nil {
		fmt.Println("Null")
		return
	}
	item := head
	for item != nil {
		fmt.Print(item.Value, "->")
		item = item.Next
	}
	fmt.Printf("null\n")
}
