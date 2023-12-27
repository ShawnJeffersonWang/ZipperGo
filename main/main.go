package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"syscall"
)

func main() {

	filename := "/home/shawn/Develop/CampusGuide/graph.txt"
	adjList, err := ReadCampusGraph(filename)
	if err != nil {
		fmt.Printf("读取文件错误: %s\n", err)
		return
	}
	shortestPath := adjList.BFS(1, 4)
	fmt.Println("Shortest path", shortestPath)
	syscall.Pause()
	var isAdmin bool
	var sourceID, targetID int

	router := gin.Default()
	router.LoadHTMLGlob("main/*") // 指定HTML模板目录

	// 设置登录页面路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})

	// 处理登录请求
	router.POST("/login", func(c *gin.Context) {
		// 获取登录表单数据
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 进行登录校验
		if username == "a" && password == "123" {
			// 登录成功，跳转到导航页面
			c.Redirect(302, "/navigation")
		} else {
			// 登录失败，返回错误提示信息
			c.HTML(200, "login.html", gin.H{
				"Error": "Invalid username or password",
			})
		}
	})

	// 设置导航页面路由
	router.GET("/navigation", func(c *gin.Context) {
		// 在这里渲染导航页面的HTML模板
		c.HTML(200, "navigation.html", gin.H{})
	})

	//// 定义GET请求的处理函数，用于显示表单页面
	//router.GET("/", func(c *gin.Context) {
	//	c.HTML(200, "index.html", nil)
	//})

	router.POST("/shortestPath", func(c *gin.Context) {
		sourceID := c.PostForm("sourceID")
		targetID := c.PostForm("targetID")

		// 将sourceID和targetID转换为整数
		source, err := strconv.Atoi(sourceID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid sourceID"})
			return
		}

		target, err := strconv.Atoi(targetID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid targetID"})
			return
		}

		// 使用Dijkstra算法查找最短路径
		path, weight := adjList.Dijkstra(source, target)

		if path == nil {
			c.JSON(404, gin.H{"error": "Path not found"})
		} else {
			c.HTML(200, "result.html", gin.H{
				"sourceID": sourceID,
				"targetID": targetID,
				"path":     path,
				"weight":   weight,
			})
		}
	})

	// 启动服务器
	router.Run(":8081")

	//router.POST("/upload", func(c *gin.Context) {
	//	name := c.PostForm("name")
	//	email := c.PostForm("email")
	//
	//	file, err := c.FormFile("file")
	//	if err != nil {
	//		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
	//		return
	//	}
	//
	//	filename := filepath.Base(file.Filename)
	//	if err := c.SaveUploadedFile(file, filename); err != nil {
	//		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
	//		return
	//	}
	//
	//	c.String(http.StatusOK, "File %s uploaded successfullu with fields name=%s and email=%s.", file.Filename, name, email)
	//})
	//router.Run(":8080")

	err = RedisLogin()
	if err != nil {
		fmt.Println("登录失败", err)
		isAdmin = false
		return
	}
	isAdmin = true

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
