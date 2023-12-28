package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("static/*") // 指定HTML模板目录

	// 设置登录页面路由
	router.Static("/images", "./images")
	router.GET("/", Index)
	router.GET("/startLogin", StartLogin)
	router.POST("/login", Login)
	router.GET("/startAdmin", StartAdmin)
	router.POST("/admin", Admin)
	router.POST("/user", User)
	router.POST("/updateMap", UpdateMap)
	router.POST("/updateRoad", UpdateRoad)
	router.POST("/shortestPath", ShortestPath)
	router.POST("/bfsPath", BFSPath)

	//// 定义GET请求的处理函数，用于显示表单页面
	//router.GET("/", func(c *gin.Context) {
	//	c.HTML(200, "index.html", nil)
	//})

	// 启动服务器
	router.Run(":8081")
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func StartLogin(c *gin.Context) {
	userType := c.Query("type")
	if userType == "admin" {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	} else if userType == "normal" {
		c.HTML(http.StatusOK, "normal.html", gin.H{})
	} else {
		c.JSON(400, "Invalid user type")
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	client, _ := InitRedis()
	savedUsername, _ := client.Get("username").Result()
	savedPassword, _ := client.Get("admin_password").Result()

	if username == savedUsername && password == savedPassword {
		c.HTML(http.StatusOK, "startAdmin.html", gin.H{})
	} else {
		c.String(http.StatusBadRequest, "Invalid username or password")
		return
	}
}

func StartAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "startAdmin.html", gin.H{})
}

func Admin(c *gin.Context) {
	choice := c.PostForm("choice")

	switch choice {
	case "1":
		c.HTML(http.StatusOK, "updateMap.html", gin.H{})
	case "2":
		c.HTML(http.StatusOK, "updateRoad.html", gin.H{})
	case "0":
		c.Redirect(302, "/")
	default:
		c.String(http.StatusBadRequest, "无效的选择")
	}
}

func UpdateMap(c *gin.Context) {
	filename := "/home/shawn/Develop/CampusGuide/graph.txt"
	adjList, err := ReadCampusGraph(filename)
	if err != nil {
		fmt.Printf("读取文件错误: %s\n", err)
		return
	}

	idStr := c.PostForm("nodeID")
	newName := c.PostForm("newName")
	nodeID, _ := strconv.Atoi(idStr)
	err = adjList.UpdateNodeName(nodeID, newName)
	if err != nil {
		return
	}

	err = SaveCampusGraph(adjList, filename)
	if err != nil {
		return
	}

	c.Redirect(302, "/startAdmin")
}

func UpdateRoad(c *gin.Context) {
	filename := "/home/shawn/Develop/CampusGuide/graph.txt"
	adjList, err := ReadCampusGraph(filename)
	if err != nil {
		fmt.Printf("读取文件错误: %s\n", err)
		return
	}

	_startVex := c.PostForm("startVex")
	_endVex := c.PostForm("endVex")
	_newWeight := c.PostForm("newWeight")
	startVex, _ := strconv.Atoi(_startVex)
	endVex, _ := strconv.Atoi(_endVex)
	newWeight, _ := strconv.Atoi(_newWeight)

	err = adjList.UpdateEdgeWeight(startVex, endVex, newWeight)
	if err != nil {
		return
	}

	err = SaveCampusGraph(adjList, filename)
	if err != nil {
		return
	}

	c.Redirect(302, "/startAdmin")
}

func User(c *gin.Context) {
	choice := c.PostForm("choice")

	switch choice {
	case "1":
		// 查看地图
		fmt.Println("校园平面图：")
		filename := "/home/shawn/Develop/CampusGuide/graph.txt"
		adjList, err := ReadCampusGraph(filename)
		if err != nil {
			fmt.Printf("读取文件错误: %s\n", err)
			return
		}

		c.HTML(http.StatusOK, "print.html", gin.H{
			"Nodes":     adjList.Nodes,
			"Adjacency": adjList.Adjacency,
		})
	case "2":
		// 寻找最优路径
		c.HTML(http.StatusOK, "dijkstra.html", gin.H{})
	case "3":
		// 不考虑权重
		c.HTML(http.StatusOK, "bfs.html", gin.H{})
	case "0":
		// 退出
		fmt.Println("退出程序")
		c.Redirect(302, "/")

	default:
		fmt.Println("无效的选项")
	}
}

func ShortestPath(c *gin.Context) {
	filename := "/home/shawn/Develop/CampusGuide/graph.txt"
	adjList, err := ReadCampusGraph(filename)
	source := c.PostForm("sourceID")
	target := c.PostForm("targetID")

	sourceID, err := strconv.Atoi(source)
	if err != nil {
		c.JSON(400, "Invalid sourceID")
		return
	}

	targetID, err := strconv.Atoi(target)
	if err != nil {
		c.JSON(400, "Invalid targetID")
		return
	}
	path, weight := adjList.Dijkstra(sourceID, targetID)
	if path == nil {
		c.JSON(404, "path not found")
	} else {
		// 将路径转换为节点名称的数组
		nodeNames := make([]string, len(path))
		for i, nodeID := range path {
			nodeNames[i] = adjList.Nodes[nodeID].Name
		}

		c.HTML(200, "shortestPath.html", gin.H{
			"sourceID": adjList.Nodes[sourceID].Name,
			"targetID": adjList.Nodes[targetID].Name,
			"path":     nodeNames,
			"weight":   weight,
		})
	}
}

func BFSPath(c *gin.Context) {
	filename := "/home/shawn/Develop/CampusGuide/graph.txt"
	adjList, err := ReadCampusGraph(filename)
	source := c.PostForm("sourceID")
	target := c.PostForm("targetID")

	sourceID, err := strconv.Atoi(source)
	if err != nil {
		c.JSON(400, "Invalid sourceID")
		return
	}

	targetID, err := strconv.Atoi(target)
	if err != nil {
		c.JSON(400, "Invalid targetID")
		return
	}
	path := adjList.BFS(sourceID, targetID)
	if path == nil {
		c.JSON(404, "path not found")
	} else {
		nodeNames := make([]string, len(path))
		for i, nodeID := range path {
			nodeNames[i] = adjList.Nodes[nodeID].Name
		}

		c.HTML(200, "shortestPath.html", gin.H{
			"sourceID": adjList.Nodes[sourceID].Name,
			"targetID": adjList.Nodes[targetID].Name,
			"path":     nodeNames,
			"weight":   nil,
		})
	}
}
