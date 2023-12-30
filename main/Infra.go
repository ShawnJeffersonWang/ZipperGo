package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

func readPassword() (string, error) {
	fmt.Print("请输入密码: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin)) // 从终端读取密码
	fmt.Println()                                                  // 打印换行
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}

func InitRedis() (*redis.Client, error) {
	// 加载配置文件
	viper.SetConfigFile("/home/shawn/config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// 从配置文件中获取Redis密码
	password := viper.GetString("redis_password")
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: password,
		DB:       0,
	})

	// 检查与Redis的连接是否成功
	pong, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("无法连接到Redis: %v", err)
	}
	fmt.Println("已连接到Redis", pong)

	return client, nil
}

func RedisLogin() error {
	client, err := InitRedis()
	if err != nil {
		return fmt.Errorf("无法连接到Redis: %v", err)
	}

	// 检查Redis中是否已经存在密码
	// :=左边得有新变量才能赋值成功，至少一个，例如passwd, err := client.Get("admin_password").Result()
	_, err = client.Get("admin_password").Result()
	// if err == redis.Nil {
	if errors.Is(err, redis.Nil) {
		// Redis中不存在密码，提示用户输入密码
		fmt.Println("请设置密码: ")
		//_, err = fmt.Scan(&passwd)
		password, err := readPassword()
		if err != nil {
			return fmt.Errorf("密码输入错误: %v", err)
		}

		// 将密码保存到Redis
		err = client.Set("admin_password", password, 0).Err()
		if err != nil {
			return fmt.Errorf("无法存储密码: %v", err)
		}
		println("密码设置成功")
	}

	// 用户登录
	fmt.Print("请输入密码: ")
	//_, err = fmt.Scan(&inputPassword)
	password, err := readPassword()
	if err != nil {
		fmt.Println("密码输入错误:", err)
		os.Exit(1)
	}

	if err != nil {
		return fmt.Errorf("密码输入错误: %v", err)
	}

	// 从Redis获取密码并与用户输入的密码进行比较
	savedPassword, err := client.Get("admin_password").Result()
	if err != nil {
		return fmt.Errorf("无法获取密码：%v", err)
	}

	if password == savedPassword {
		fmt.Println("密码正确，登录成功")
	} else {
		return fmt.Errorf("密码错误，登录失败")
	}

	// 关闭Redis连接
	err = client.Close()
	if err != nil {
		return fmt.Errorf("关闭Redis连接失败: %v", err)
	}

	return nil
}
