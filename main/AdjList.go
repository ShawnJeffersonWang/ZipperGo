package main

import (
	"fmt"
	"math"
)

// Node 表示校园平面图中的一个节点
type Node struct {
	ID   int
	Name string
}

// Edge 表示校园平面图中的一条边
type Edge struct {
	StartVex int
	EndVex   int
	Weight   int
}

// AdjList 表示校园平面图的邻接表
type AdjList struct {
	Nodes     []Node
	Adjacency map[int][]Edge
}

// NewAdjList 创建一个新的校园平面图的邻接表
func NewAdjList() *AdjList {
	return &AdjList{
		Nodes:     []Node{},
		Adjacency: make(map[int][]Edge),
	}
}

// AddNode 添加一个节点到校园平面图中
func (g *AdjList) AddNode(node Node) {
	g.Nodes = append(g.Nodes, node)
}

// AddEdge 向邻接表中添加一条边
func (g *AdjList) AddEdge(edge Edge) {
	g.Adjacency[edge.StartVex] = append(g.Adjacency[edge.StartVex], edge)
}

// GetOutEdges 获取指定节点的出边列表
func (g *AdjList) GetOutEdges(nodeID int) []Edge {
	return g.Adjacency[nodeID]
}

// UpdateNodeName 更新节点的名称
func (g *AdjList) UpdateNodeName(nodeID int, newName string) error {
	for i, node := range g.Nodes {
		if node.ID == nodeID {
			g.Nodes[i].Name = newName
			return nil
		}
	}
	return fmt.Errorf("节点 %d 不存在", nodeID)
}

// UpdateEdgeWeight 更新边的权重
func (g *AdjList) UpdateEdgeWeight(startVex, endVex, newWeight int) error {
	edges, ok := g.Adjacency[startVex]
	if !ok {
		return fmt.Errorf("节点 %d 不存在", startVex)
	}

	for i, edge := range edges {
		if edge.EndVex == endVex {
			g.Adjacency[startVex][i].Weight = newWeight
			return nil
		}
	}
	return fmt.Errorf("边 (%d, %d) 不存在", startVex, endVex)
}

// Print 打印校园平面图的内容
func (g *AdjList) Print() {
	fmt.Println("节点信息：")
	for _, node := range g.Nodes {
		fmt.Printf("节点ID: %d, 节点名称: %s\n", node.ID, node.Name)
	}

	fmt.Println("边信息：")
	for _, node := range g.Nodes {
		edges := g.GetOutEdges(node.ID)
		for _, edge := range edges {
			fmt.Printf("起始节点ID: %d, 结束节点ID: %d, 边权重: %d\n", edge.StartVex, edge.EndVex, edge.Weight)
		}
	}
}

// Dijkstra 计算校园平面图中两个节点之间的最短路径
func (g *AdjList) Dijkstra(sourceID int, targetID int) ([]int, int) {
	// 初始化距离和前驱节点
	dist := make(map[int]int)
	prev := make(map[int]int)
	for _, node := range g.Nodes {
		dist[node.ID] = math.MaxInt32
		prev[node.ID] = -1
	}
	dist[sourceID] = 0

	// 创建一个保存已访问节点的集合
	visited := make(map[int]bool)

	// 执行 Dijkstra 算法
	for len(visited) < len(g.Nodes) {
		// 选择距离最小的未访问节点
		minDist := math.MaxInt32
		minNode := -1
		for _, node := range g.Nodes {
			if !visited[node.ID] && dist[node.ID] < minDist {
				minDist = dist[node.ID]
				minNode = node.ID
			}
		}

		// 标记该节点为已访问
		visited[minNode] = true

		// 更新与当前节点相邻节点的距离
		edges := g.GetOutEdges(minNode)
		for _, edge := range edges {
			newDist := dist[minNode] + edge.Weight
			if newDist < dist[edge.EndVex] {
				dist[edge.EndVex] = newDist
				prev[edge.EndVex] = minNode
			}
		}
	}

	// 构造最短路径
	var path []int
	current := targetID
	for current != -1 {
		path = append([]int{current}, path...)
		current = prev[current]
	}

	return path, dist[targetID]
}

// BFS 使用广度优先搜索算法查找最短路径
func (g *AdjList) BFS(startID, targetID int) []int {
	visited := make(map[int]bool) // 记录节点是否已访问
	queue := [][]int{{startID}}   // 使用队列保存路径，初始路径只包含起始节点
	for len(queue) > 0 {
		path := queue[0] // 取出队列中的第一个路径
		queue = queue[1:]
		nodeID := path[len(path)-1] // 获取当前路径的最后一个节点
		if nodeID == targetID {     // 如果找到目标节点，返回路径
			return path
		}
		if !visited[nodeID] { // 如果当前节点未被访问，则继续扩展路径
			visited[nodeID] = true          // 标记当前节点为已访问
			adjEdges := g.Adjacency[nodeID] // 获取当前节点的相邻边
			for _, edge := range adjEdges {
				newPath := append(path, edge.EndVex) // 将当前边的终点添加到路径中
				queue = append(queue, newPath)       // 将新路径加入队列
			}
		}
	}
	return nil // 如果未找到最短路径，返回空路径
}

// DFS1 使用深度优先搜索算法查找最短路径
func (g *AdjList) DFS1(startID, targetID int) []int {
	visited := make(map[int]bool) // 记录节点是否已访问
	var shortestPath []int        // 最短路径
	path := []int{startID}        // 当前路径
	g.dfsHelper(startID, targetID, visited, path, &shortestPath)
	return shortestPath
}

// dfsHelper 是DFS1的辅助函数，用于递归搜索最短路径
func (g *AdjList) dfsHelper(currentID, targetID int, visited map[int]bool, path []int, shortestPath *[]int) {
	if currentID == targetID { // 如果当前节点是目标节点
		if len(*shortestPath) == 0 || len(path) < len(*shortestPath) {
			*shortestPath = append([]int(nil), path...) // 更新最短路径
		}
		return
	}
	visited[currentID] = true // 标记当前节点为已访问
	adjEdges := g.Adjacency[currentID]
	for _, edge := range adjEdges {
		if !visited[edge.EndVex] {
			newPath := append(path, edge.EndVex) // 添加当前边的终点到路径中
			g.dfsHelper(edge.EndVex, targetID, visited, newPath, shortestPath)
		}
	}
	visited[currentID] = false // 回溯时取消当前节点的访问标记
}

// DFS 使用深度优先搜索算法查找最短路径
func (g *AdjList) DFS(startID, targetID int) []int {
	visited := make(map[int]bool)     // 记录节点是否已访问
	stack := []int{startID}           // 使用栈模拟DFS
	path := make(map[int][]int)       // 记录路径
	path[startID] = []int{startID}     // 起始节点的路径只包含自身

	for len(stack) > 0 {
		currentID := stack[len(stack)-1] // 获取栈顶节点
		stack = stack[:len(stack)-1]     // 出栈
		if currentID == targetID {       // 如果找到目标节点，返回路径
			return path[currentID]
		}
		if !visited[currentID] {        // 如果当前节点未被访问，则继续搜索
			visited[currentID] = true    // 标记当前节点为已访问
			adjEdges := g.Adjacency[currentID]
			for _, edge := range adjEdges {
				if !visited[edge.EndVex] {
					stack = append(stack, edge.EndVex) // 将相邻节点入栈
					newPath := append(path[currentID], edge.EndVex) // 更新路径
					path[edge.EndVex] = newPath
				}
			}
		}
	}

	return nil // 如果未找到最短路径，返回空路径
}
