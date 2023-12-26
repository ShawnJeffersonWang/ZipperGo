package main

import (
	"fmt"
)

func main() {
	filename := "/home/shawn/Develop/CampusGuide/graph.txt"
	adjList, err := ReadCampusGraph(filename)
	if err != nil {
		fmt.Printf("读取文件错误: %s\n", err)
		return
	}
	var isAdmin bool
	var sourceID, targetID int

	err = RedisLogin()
	if err != nil {
		fmt.Println("登录失败", err)
		isAdmin = false
		return
	}

	// 执行不同的操作
	for {
		if isAdmin {
			// 管理员操作
			fmt.Println("请选择操作：")
			fmt.Println("1. 更新地图")
			fmt.Println("2. 更新道路")
			fmt.Println("3. 修改密码")
			fmt.Println("0. 退出")

			var choice int
			fmt.Print("请输入选项：")
			_, err = fmt.Scan(&choice)
			if err != nil {
				fmt.Println("输入错误")
				continue
			}

			switch choice {
			case 1:
				// 更新地图
				fmt.Println("更新地图")
				// 调用相应的函数进行地图更新操作
				var nodeID int
				var newName string
				_, err := fmt.Scan(&nodeID, &newName)
				if err != nil {
					return
				}

				err = adjList.UpdateNodeName(nodeID, newName)
				if err != nil {
					return
				}
			case 2:
				// 更新道路
				fmt.Println("更新道路")
				var startVex, endVex, newWeight int

				_, err := fmt.Scan(&startVex, endVex, newWeight)
				if err != nil {
					return
				}

				err = adjList.UpdateEdgeWeight(startVex, endVex, newWeight)
				if err != nil {
					return
				}
			case 3:
				fmt.Println("修改密码, 请输入密码")
				var passwd int
				_, err := fmt.Scan(&passwd)
				if err != nil {
					fmt.Println("密码输入错误")
					return
				}

				client, err := InitRedis()
				if err != nil {
					fmt.Println("Redis创建失败")
					return
				}

				client.Set("admin_password", &passwd, 0)
			case 0:
				// 退出
				fmt.Println("退出程序")
				return
			default:
				fmt.Println("无效的选项")
			}
		} else {
			// 普通用户操作
			fmt.Println("请选择操作：")
			fmt.Println("1. 查看地图")
			fmt.Println("2. 查看景点")
			fmt.Println("3. 寻找最优路径")
			fmt.Println("0. 退出")
			var choice int
			fmt.Print("请输入选项：")
			_, err = fmt.Scan(&choice)
			if err != nil {
				fmt.Println("输入错误")
				continue
			}

			switch choice {
			case 1:
				// 查看地图
				fmt.Println("校园平面图：")
				adjList.Print()
			case 2:
				// 查看景点
				fmt.Println("执行查看景点操作")
				// 调用相应的函数进行景点查看操作
			case 3:
				// 寻找最优路径
				fmt.Println("请输入要查找的sourceID和targetID：")
				_, err = fmt.Scan(&sourceID, &targetID)
				if err != nil {
					fmt.Println("输入错误")
					continue
				}
				path, distance := adjList.Dijkstra(sourceID, targetID)
				fmt.Println("最短路径：", path)
				fmt.Println("路径长度：", distance)
			case 0:
				// 退出
				fmt.Println("退出程序")
				return
			default:
				fmt.Println("无效的选项")
			}
		}
	}
}
