package main

import (
	"encoding/json"
	"github.com/go-vgo/robotgo"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	configFilePath = "config.json"
)

// Program 定义配置文件中的程序结构
type Program struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Server struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type Config struct {
	Token    string    `json:"token"`
	Server   Server    `json:"server"`
	Programs []Program `json:"programs"`
}

var config Config

// authHandler 添加HTTP基本认证
func authHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token != config.Token {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Restricted\"")
			http.Error(w, "未授权", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

// playPauseHandler 模拟播放/暂停键
func playPauseHandler(w http.ResponseWriter, r *http.Request) {
	robotgo.KeyTap("audio_play")
	w.Write([]byte("已模拟播放/暂停"))
}

// previousHandler 模拟上一曲键
func previousHandler(w http.ResponseWriter, r *http.Request) {
	robotgo.KeyTap("audio_prev")
	w.Write([]byte("已模拟上一曲"))
}

// nextHandler 模拟下一曲键
func nextHandler(w http.ResponseWriter, r *http.Request) {
	robotgo.KeyTap("audio_next")
	w.Write([]byte("已模拟下一曲"))
}

// volumeUpHandler 模拟音量加键
func volumeUpHandler(w http.ResponseWriter, r *http.Request) {
	robotgo.KeyTap("audio_vol_up")
	w.Write([]byte("已模拟音量加"))
}

// volumeDownHandler 模拟音量减键
func volumeDownHandler(w http.ResponseWriter, r *http.Request) {
	robotgo.KeyTap("audio_vol_down")
	w.Write([]byte("已模拟音量减"))
}

// shutdownHandler 执行关机命令
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("shutdown", "/s", "/t", "0")
	err := cmd.Run()
	if err != nil {
		log.Println("关机失败:", err)
		http.Error(w, "关机失败", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("正在关机"))
}

// rebootHandler 执行重启命令
func rebootHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("shutdown", "/r", "/t", "0")
	err := cmd.Run()
	if err != nil {
		log.Println("重启失败:", err)
		http.Error(w, "重启失败", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("正在重启"))
}

// loadConfig 读取并解析配置文件
func loadConfig() {
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal("无法打开配置文件:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("配置文件解析失败:", err)
	}
}

// launchHandler 处理快速启动请求
func launchHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求中的程序名称参数
	programName := r.URL.Query().Get("name")
	if programName == "" {
		http.Error(w, "缺少程序名称参数", http.StatusBadRequest)
		return
	}

	// 在配置文件中查找匹配的程序
	for _, p := range config.Programs {
		if strings.EqualFold(p.Name, programName) { // 不区分大小写比较
			cmd := exec.Command(p.Path)
			err := cmd.Start()
			if err != nil {
				log.Printf("启动程序 %s 失败: %v", p.Name, err)
				http.Error(w, "启动程序失败", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("已启动程序: " + p.Name))
			return
		}
	}
	http.Error(w, "未找到指定程序", http.StatusNotFound)
}

func main() {
	loadConfig()
	// 设置HTTP路由并添加认证
	http.HandleFunc("/play_pause", authHandler(playPauseHandler))
	http.HandleFunc("/previous", authHandler(previousHandler))
	http.HandleFunc("/next", authHandler(nextHandler))
	http.HandleFunc("/volume_up", authHandler(volumeUpHandler))
	http.HandleFunc("/volume_down", authHandler(volumeDownHandler))
	http.HandleFunc("/shutdown", authHandler(shutdownHandler))
	http.HandleFunc("/reboot", authHandler(rebootHandler))
	http.HandleFunc("/launch", authHandler(launchHandler))

	host := config.Server.Host + ":" + config.Server.Port
	log.Println("服务器启动在 " + host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
