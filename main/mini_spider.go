package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"mini_spider/scheduler"
	"mini_spider/web_package"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)


type CommandParams struct {
	configPath  string
	logPath     []string
	showVersion bool
	showHelp    bool
}


func getConfig(configPath string) web_package.Config {
	var config web_package.Config

	if len(configPath) == 0 {
		return config
	}
	file, _ := os.Open(configPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	configMap := make(map[string]string)
	count := 0
	for scanner.Scan() {
		count++
		// 读取当前行内容
		line := scanner.Text()
		strArray := strings.Split(line, "=")
		if len(strArray) > 1 {
			var key  = strings.TrimSpace(strArray[0]);
			var value = strings.TrimSpace(strArray[1]);
			configMap[key] = value
		}
		fmt.Printf("%d %s\n", count, line)
	}

	if len(configMap) == 0 {
		return config
	}
	test := configMap ["urlListFile"]
	println(test);
	config.UrlListFile = configMap ["urlListFile"]
	config.OutputDirectory = configMap["outputDirectory"]
	config.MaxDepth, _ = strconv.Atoi(configMap["maxDepth"])
	config.CrawlInterval, _ = strconv.Atoi(configMap["crawlInterval"])
	config.CrawlTimeout, _ = strconv.Atoi(configMap["crawlTimeout"])
	config.TargetUrl = configMap["targetUrl"]
	config.TreadCount, _ = strconv.Atoi(configMap["threadCount"])
	return config
}

func getCrawlUrls(path string) []string {

	content, error := os.ReadFile(path)
	if error != nil {
		return nil
	}

	var m []string
	jsonConvertError := json.Unmarshal(content, &m)
	if jsonConvertError != nil {
		return m
	}
	return m
}

func handleCommandLine() CommandParams {
	var commandParams CommandParams
	for i := 0; i < len(os.Args); i++ {
		str := os.Args[i]
		switch str {
		case "-v":
			commandParams.showVersion = true
		case "-h":
			commandParams.showHelp = true
		case "-c":
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				commandParams.configPath = os.Args[i+1]
			}

		case "-l":
			var array []string
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				array = append(array, os.Args[i+1])
			}
			if i+2 < len(os.Args) && !strings.HasPrefix(os.Args[i+2], "-") {
				array = append(array, os.Args[i+2])
			}
			commandParams.logPath = array

		}
	}
	return commandParams
}

func main() {
	commandParams := handleCommandLine()
	if commandParams.showHelp {
		//TODO:yangxu 显示help
	} else if commandParams.showVersion {
		//TODO:yangxu 显示版本号
	}
	if len(commandParams.configPath) == 0 {
		return
	}

	var config = getConfig(commandParams.configPath)
	var urls = getCrawlUrls(config.UrlListFile)
	var schedulerInstance, _ = scheduler.NewScheduler(&config, urls)
	schedulerInstance.Run()

	go func() {
		for {
			if schedulerInstance.GetUnfinished() == 0 {
				os.Exit(0)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	//捕获程序终止，
	ch := make(chan os.Signal,1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	os.Exit(0)
}
