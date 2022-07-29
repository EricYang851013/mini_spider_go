package web_package


type Config struct {
	//种子文件路径
	UrlListFile string
	//输出文件路径
	OutputDirectory string
	//最大抓取深度
	MaxDepth int
	//抓取间隔 秒
	CrawlInterval int
	//抓取超时  秒
	CrawlTimeout int
	//需要存储的目标网页URL pattern
	TargetUrl string
	//抓取线程数
	TreadCount int
}