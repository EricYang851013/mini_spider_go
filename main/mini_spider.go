package main

import(
	"bufio"
    "fmt"
	"os"
	"os/signal"
	"encoding/json"
)


import (
	"scheduler"
)

type CommandParams struct {
	configPath string
	logPath [...]string
	showVersion bool
	showHelp    bool
}

type Config struct {
	//种子文件路径
	urlListFile string
	//输出文件路径
	outputDirectory string
	//最大抓取深度
	maxDepth int
	//抓取间隔 秒
	crawlInterval int 
	//抓取超时  秒
	crawlTimeout int
	//需要存储的目标网页URL pattern
	targetUrl  string
	//抓取线程数
	treadCount int
}

func getConfig(configPath string) Config {
	if len(configPath) == 0 {
		return nil;
	}
	file, _ := os.Open(configPath)
	defer file.Close()
    scanner := bufio.NewScanner(file)
	var config Config;
     var configMap := make(map[string]string)
	for scanner.Scan() {
        count++
        // 读取当前行内容
		line := scanner.Text()
		strArray = strings.Split(line, "=");
		if len(strArray)  > 1{
			configMap[strArray[0]] = strArray[1]
		}
        fmt.Printf("%d %s\n", count, line)
	}
	
	if len(configMap)  == 0{
		return nil;
	}
	config.urlListFile = configMapp["urlListFile"]
	config.outputDirectory = configMapp["outputDirectory"]
	config.maxDepth = configMapp["maxDepth"]
	config.crawlInterval = configMapp["crawlInterval"]
	config.crawlTimeout = configMapp["crawlTimeout"]
	config.targetUrl = configMapp["targetUrl"]
	config.treadCount = configMapp["treadCount"]
	return config;
}

func getCrawlUrls(path string) []string {
	
	content, error := os.ReadFile(path)
	if error != nil {
		return nil
	}

	m := []string()
	jsonConvertError := json.Unmarshal(content, &m)
	if jsonConvertError != nil {
		return nil;
	}
	return m;
}

func  handleCommandLine() CommandParams {
	var commandParams CommandParams;
	for i := 0; i < len(os.Args); i++ {
		var str := os.Args[i];
		switch str  {
		case "-v":
			commandParams.showVersion = true;
		case "-h":
			commandParams.showHelp = true;
		case "-c":
			if i + 1 < len(os.Args) && !string.HasPrefix(os.Args[i + 1], "-"){
				commandParams.configPath = os.Args[i + 1];
			}

		case "-h":
			var array  [2] string
			if i + 1 < len(os.Args) && !string.HasPrefix(os.Args[i + 1], "-"){
				array[0] := os.Args[i + 1]
			}
			if i + 2 < len(os.Args) && !string.HasPrefix(os.Args[i + 2], "-"){
				array[1] := os.Args[i + 2]
			}
			commandParams.logPath = array

		}
	}
	return commandParams;
}


func main()  {
	var  commandParams := handleCommandLine();
	if commandParams.showHelp {
		//TODO:yangxu 显示help
	}
	else if commandParams.showVersion {
		//TODO:yangxu 显示版本号
	}
	if len(commandParams.configPath) == 0 {
		return;
	}

	var config = getConfig(commandParams.configPath);
	var urls = getCrawlUrls(config.urlListFile);
	var schedulerInstance = scheduler.NewScheduler(config,urls);
	schedulerInstance.Run();

	go func ()  {
		for  {
			if schedulerInstance.GetUnfinished() == 0 {
				Exit(0)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	//捕获程序终止，
	ch := make(chan os.Signal);
	signal.Notify(ch,syscall.SIGINT, syscall.SIGTERM);
	<-ch


	Exit(0)
}