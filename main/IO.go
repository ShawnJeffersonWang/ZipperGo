package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
)

// ReadCampusGraph 从 Redis 中读取数据并构建校园平面图
func ReadCampusGraph(redisClient *redis.Client) (*AdjList, error) {
	adjList := NewAdjList()

	nodeKeys := redisClient.Keys("nodes:*").Val()
	for _, nodeKey := range nodeKeys {
		nodeIDStr := strings.TrimPrefix(nodeKey, "nodes:")
		nodeID, err := strconv.Atoi(nodeIDStr)
		if err != nil {
			return nil, fmt.Errorf("节点ID解析错误：%s", err)
		}

		nodeName, err := redisClient.Get(nodeKey).Result()
		if err != nil {
			return nil, fmt.Errorf("无法获取节点名称：%s", err)
		}

		node := Node{ID: nodeID, Name: nodeName}
		adjList.Nodes = append(adjList.Nodes, node)
	}

	edgeKeys := redisClient.Keys("edges:*").Val()
	for _, edgeKey := range edgeKeys {
		edgeIDStr := strings.TrimPrefix(edgeKey, "edges:")
		elements := strings.Split(edgeIDStr, "_")
		if len(elements) != 2 {
			return nil, fmt.Errorf("边ID解析错误：%s", edgeIDStr)
		}

		startID, err := strconv.Atoi(elements[0])
		if err != nil {
			return nil, fmt.Errorf("起始节点ID解析错误：%s", err)
		}
		endID, err := strconv.Atoi(elements[1])
		if err != nil {
			return nil, fmt.Errorf("结束节点ID解析错误：%s", err)
		}

		weightStr, err := redisClient.Get(edgeKey).Result()
		if err != nil {
			return nil, fmt.Errorf("无法获取边权重：%s", err)
		}
		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			return nil, fmt.Errorf("边权重解析错误：%s", err)
		}

		edge := Edge{StartVex: startID, EndVex: endID, Weight: weight}
		adjList.AddEdge(edge)
	}

	return adjList, nil
}

// SaveCampusGraph 将校园平面图保存到 Redis 中
func SaveCampusGraph(adjList *AdjList, redisClient *redis.Client) error {
	// 保存节点信息
	for _, node := range adjList.Nodes {
		nodeKey := fmt.Sprintf("nodes:%d", node.ID)
		err := redisClient.Set(nodeKey, node.Name, 0).Err()
		if err != nil {
			return fmt.Errorf("保存节点信息失败：%s", err)
		}
	}

	// 保存边信息
	for _, node := range adjList.Nodes {
		edges := adjList.GetOutEdges(node.ID)
		for _, edge := range edges {
			edgeKey := fmt.Sprintf("edges:%d_%d", edge.StartVex, edge.EndVex)
			err := redisClient.Set(edgeKey, strconv.Itoa(edge.Weight), 0).Err()
			if err != nil {
				return fmt.Errorf("保存边信息失败：%s", err)
			}
		}
	}

	return nil
}

////ReadCampusGraph 从文件中读取数据并构建校园平面图
//func ReadCampusGraph(filename string) (*AdjList, error) {
//	file, err := os.Open(filename)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	adjList := NewAdjList()
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		line := scanner.Text()
//		if line == "" {
//			continue
//		}
//
//		elements := strings.Split(line, " ")
//		if len(elements) == 2 {
//			// 处理节点行
//			nodeID, err := strconv.Atoi(elements[0])
//			if err != nil {
//				return nil, fmt.Errorf("节点ID解析错误：%s", err)
//			}
//			nodeName := elements[1]
//
//			node := Node{ID: nodeID, Name: nodeName}
//			adjList.Nodes = append(adjList.Nodes, node)
//		} else if len(elements) == 3 {
//			// 处理边行
//			startID, err := strconv.Atoi(elements[0])
//			if err != nil {
//				return nil, fmt.Errorf("起始节点ID解析错误：%s", err)
//			}
//			endID, err := strconv.Atoi(elements[1])
//			if err != nil {
//				return nil, fmt.Errorf("结束节点ID解析错误：%s", err)
//			}
//			weight, err := strconv.Atoi(elements[2])
//			if err != nil {
//				return nil, fmt.Errorf("边权重解析错误：%s", err)
//			}
//
//			edge := Edge{StartVex: startID, EndVex: endID, Weight: weight}
//			adjList.AddEdge(edge)
//		}
//	}
//
//	if err := scanner.Err(); err != nil {
//		return nil, fmt.Errorf("读取文件错误：%s", err)
//	}
//
//	return adjList, nil
//}

//// SaveCampusGraph 将校园平面图保存到文件中
//func SaveCampusGraph(adjList *AdjList, filename string) error {
//	file, err := os.Create(filename)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	writer := bufio.NewWriter(file)
//
//	// 写入节点信息
//	for _, node := range adjList.Nodes {
//		line := fmt.Sprintf("%d %s\n", node.ID, node.Name)
//		_, err := writer.WriteString(line)
//		if err != nil {
//			return err
//		}
//	}
//
//	// 写入边信息
//	for _, node := range adjList.Nodes {
//		edges := adjList.GetOutEdges(node.ID)
//		for _, edge := range edges {
//			line := fmt.Sprintf("%d %d %d\n", edge.StartVex, edge.EndVex, edge.Weight)
//			_, err := writer.WriteString(line)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	err = writer.Flush()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
