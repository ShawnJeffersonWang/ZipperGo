package Miscellaneous

import (
	"fmt"
)

// 图的邻接表节点
type AdjacencyListNode struct {
	Vertex int                // 顶点
	Next   *AdjacencyListNode // 指向下一个邻接节点的指针
}

// 图的邻接表
type AdjacencyList struct {
	Vertices int                  // 顶点数
	Nodes    []*AdjacencyListNode // 邻接表节点数组
}

// 添加边
func (g *AdjacencyList) AddEdge(src, dest int) {
	// 创建新的邻接节点
	newNode := &AdjacencyListNode{
		Vertex: dest,
		Next:   nil,
	}

	// 将新节点插入到源节点的邻接表中
	if g.Nodes[src] == nil {
		g.Nodes[src] = newNode
	} else {
		currNode := g.Nodes[src]
		for currNode.Next != nil {
			currNode = currNode.Next
		}
		currNode.Next = newNode
	}

	// 如果是无向图，需要添加反向边
	// 创建新的反向邻接节点
	newReverseNode := &AdjacencyListNode{
		Vertex: src,
		Next:   nil,
	}

	// 将新反向节点插入到目标节点的邻接表中
	if g.Nodes[dest] == nil {
		g.Nodes[dest] = newReverseNode
	} else {
		currReverseNode := g.Nodes[dest]
		for currReverseNode.Next != nil {
			currReverseNode = currReverseNode.Next
		}
		currReverseNode.Next = newReverseNode
	}
}

// 打印邻接表
func (g *AdjacencyList) Print() {
	for i := 0; i < g.Vertices; i++ {
		fmt.Printf("顶点 %d 的邻接表：", i)
		currNode := g.Nodes[i]
		for currNode != nil {
			fmt.Printf(" -> %d", currNode.Vertex)
			currNode = currNode.Next
		}
		fmt.Println()
	}
}

func main() {
	// 创建一个包含5个顶点的图
	graph := &AdjacencyList{
		Vertices: 5,
		Nodes:    make([]*AdjacencyListNode, 5),
	}

	// 添加边
	graph.AddEdge(0, 1)
	graph.AddEdge(0, 4)
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(1, 4)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)

	// 打印邻接表
	graph.Print()
}
