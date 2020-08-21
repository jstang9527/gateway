package selm2

import (
	"fmt"
	"regexp"
	"strings"
)

// Action ...
type Action struct {
	Info   string   //关键词
	Path   []string //关键词原始绝对路径
	First  []string //*[id="root"]/div/ul/li/button/div/div ==> //*[id="root"]/div/ul/li/button
	Link   []string //*[id="root"]/div/ul/li/a/div  ==>//*[id="root"]/div/ul/li/a
	Button []string //*[id="root"]/div/ul/li/button/div/div ==> //*[id="root"]/div/ul/li/button
	Input  []string //*[@id="root"]/div/input
	Table  []string //*[@id="root"]/div/div/div[3]/div[2]/div[3]/div/div/div/div/div[1]/div/div/div/div/div/table/thead/tr/th[9]

}

// ActionDTO 输出结构体
type ActionDTO struct {
	Prefix string
	Info   string
	Path   string
	First  string
	Link   string
	Button string
	Input  string
	Table  string
}

// Reverse 数组倒序函数
func (a *Action) Reverse() *Action {
	arr := a.Path
	var length = len(arr)
	var temp string
	for i := 0; i < length/2; i++ {
		temp = (arr)[i]
		(arr)[i] = (arr)[length-1-i]
		(arr)[length-1-i] = temp
	}
	a.Path = arr
	return a
}

// Selector 利用Path初始化其他标签数组
func (a *Action) Selector() *Action {
	if len(a.Path) > 1 {
		a.First = a.Path[:len(a.Path)-1] //前一个标签，如果当前为li,那么前一个就是ul,就可以对其他li进行操作。比如批量删除
	} else {
		a.First = a.Path
	}

	for index, item := range a.Path {
		if is, _ := regexp.MatchString(`a\[\d+]$`, item); is { //a标签
			a.Link = a.Path[:index+1]
			continue
		}
		if is, _ := regexp.MatchString(`button`, item); is { //button标签
			a.Button = a.Path[:index+1]
			continue
		}
		if is, _ := regexp.MatchString(`input`, item); is { //input标签
			a.Input = a.Path[:index+1]
			continue
		}
		if is, _ := regexp.MatchString(`table`, item); is { //table标签
			a.Table = a.Path[:index+1]
			continue
		}
	}
	return a
}

// ToString 将数组转换成字符串
func (a *Action) ToString() *ActionDTO {
	prefix := `//*[@id="root"]/`
	return &ActionDTO{
		Prefix: prefix,
		Info:   a.Info,
		Path:   fmt.Sprintf("%s%v", prefix, strings.Join(a.Path, "/")),
		First:  fmt.Sprintf("%s%v", prefix, strings.Join(a.First, "/")),
		Link:   fmt.Sprintf("%s%v", prefix, strings.Join(a.Link, "/")),
		Button: fmt.Sprintf("%s%v", prefix, strings.Join(a.Button, "/")),
		Input:  fmt.Sprintf("%s%v", prefix, strings.Join(a.Input, "/")),
		Table:  fmt.Sprintf("%s%v", prefix, strings.Join(a.Table, "/")),
	}
}

// *****************不在append路径，而是append对象，拿attr属性和name属性************************

// Search ...
func (a *Action) Search(nodes []*H2j) (ok bool) {
	temp := make(map[string]int, 10)
	for _, node := range nodes {
		temp[node.Name]++
		// 1.看看本节点是否是target
		if strings.Contains(node.Text, a.Info) {
			// fmt.Println(node.Type, node.Text)
			return true
		}

		// 2.不是的话,看本节点的孩子
		if node.Children != nil {
			c1 := node.Children
			if ok := a.Search(c1); ok {
				a.Path = append(a.Path, fmt.Sprintf("%v[%v]", node.Name, temp[node.Name]))
				return true
			}
		}
	}
	return false
}
