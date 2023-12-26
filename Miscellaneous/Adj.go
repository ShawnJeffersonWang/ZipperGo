package Miscellaneous

import (
	"fmt"
	"math"
)

type ArcNode struct {
	AdjVex int
	Weight int
	Next   *ArcNode
}

type Node struct {
	ID   int
	Name string
	Head *ArcNode
}

type AdjList struct {
	ArcNum int
	VexNum int
	Nodes  []*Node
}

//func NewGraph() *AdjList {
//	newAdjList := AdjList{}
//	newNode := &Node{}
//	data, err := os.ReadFile("campus_data.json")
//	if err != nil {
//		return nil
//	}
//
//	err = json.Unmarshal(data, &newNode)
//	if err != nil {
//		return nil
//	}
//
//	newArcNode := &ArcNode{}
//	data, err = os.ReadFile("campus_data2.json")
//	if err != nil {
//		return nil
//	}
//	err = json.Unmarshal(data, &newArcNode)
//	i := 0
//	for ;i< ; i++ {
//		if i == 0 {
//			newNode.Head = newArcNode
//		} else {
//			currNode := newNode.Head
//			for ; currNode.Next != nil; currNode = currNode.Next {
//			}
//			currNode.Next = newArcNode
//		}
//	}
//	newAdjList.ArcNum = i
//}

func NewAdjList(vexNum int) *AdjList {
	return &AdjList{
		VexNum: vexNum,
		Nodes:  make([]*Node, vexNum),
	}
}

func (adjList *AdjList) AddNode(id int, name string) {
	newNode := &Node{
		ID:   id,
		Name: name,
		Head: nil,
	}
	adjList.Nodes[id] = newNode
}

func (adjList *AdjList) AddEdge(startID, endID, weight int) {
	newArcNode := &ArcNode{
		AdjVex: endID,
		Weight: weight,
		Next:   nil,
	}
	if adjList.Nodes[startID].Head == nil {
		adjList.Nodes[startID].Head = newArcNode
	} else {
		current := adjList.Nodes[startID].Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newArcNode
	}
	adjList.ArcNum++
}

func (adjList *AdjList) GetOutEdges(nodeID int) []*ArcNode {
	var edges []*ArcNode
	current := adjList.Nodes[nodeID].Head
	for current != nil {
		edges = append(edges, current)
		current = current.Next
	}
	return edges
}

func (adjList *AdjList) GetInEdges(nodeID int) []*ArcNode {
	//edges := []*ArcNode{}
	var edges []*ArcNode
	for _, node := range adjList.Nodes {
		current := node.Head
		for current != nil {
			if current.AdjVex == nodeID {
				edges = append(edges, current)
				break
			}
			current = current.Next
		}
	}
	return edges
}

// BFS performs breadth-first search starting from a given node in the adjacency list.
func (adjList *AdjList) BFS(startID int) {
	visited := make([]bool, adjList.VexNum) // 用于标记节点是否已访问过
	queue := make([]int, 0)                 // 用于存储待访问的节点

	visited[startID] = true
	queue = append(queue, startID)

	for len(queue) > 0 {
		nodeID := queue[0] // 取出队首节点
		queue = queue[1:]  // 出队

		fmt.Printf("Visited node: %d\n", nodeID)

		// 获取当前节点的所有邻接节点
		edges := adjList.GetOutEdges(nodeID)
		for _, edge := range edges {
			adjNodeID := edge.AdjVex
			if !visited[adjNodeID] {
				visited[adjNodeID] = true
				queue = append(queue, adjNodeID) // 邻接节点入队
			}
		}
	}
}

// Dijkstra finds the shortest path from a given start node to all other nodes in the adjacency list.
func (adjList *AdjList) Dijkstra(startID int) map[int]int {
	dist := make(map[int]int) // 用于记录起始节点到各节点的最短距离

	// 初始化距离为无穷大
	for _, node := range adjList.Nodes {
		dist[node.ID] = math.MaxInt32
	}

	dist[startID] = 0 // 起始节点到自身的距离为0

	visited := make(map[int]bool) // 用于标记节点是否已访问过

	for i := 0; i < adjList.VexNum-1; i++ {
		// 找到当前距离起始节点最近的未访问节点
		minDist := math.MaxInt32
		minDistNode := -1
		for _, node := range adjList.Nodes {
			if !visited[node.ID] && dist[node.ID] < minDist {
				minDist = dist[node.ID]
				minDistNode = node.ID
			}
		}

		if minDistNode == -1 {
			// 所有节点都已访问，退出循环
			break
		}

		visited[minDistNode] = true

		// 更新与当前节点邻接的节点的最短距离
		edges := adjList.GetOutEdges(minDistNode)
		for _, edge := range edges {
			adjNodeID := edge.AdjVex
			newDist := dist[minDistNode] + edge.Weight
			if newDist < dist[adjNodeID] {
				dist[adjNodeID] = newDist
			}
		}
	}

	return dist
}
