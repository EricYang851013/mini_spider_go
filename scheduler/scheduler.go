package scheduler
import(
	"mini_spider/web_package"
)

type Scheduler struct {
	config *web_package.Config //配置信息
	urlTable *UrlTable  //记录已经抓取的url
	queue    Queue      //任务队列
	crawlers []*Crawler //抓取线程
}

// crawl task
type CrawlTask struct {
	Url    string            // url to crawl
	Depth  int               // depth of the url
	Header map[string]string // http header
}

// create new mini-spider
func NewScheduler(conf *web_package.Config, seeds []string) (*Scheduler, error) {
	scheduler := new(Scheduler)
	scheduler.config = conf

	// create url table
	scheduler.urlTable = NewUrlTable()

	// initialize queue
	scheduler.queue.Init()

	// add seeds to queue
	for _, seed := range seeds {
		task := &CrawlTask{Url: seed, Depth: 1, Header: make(map[string]string)}
		scheduler.queue.Add(task)
	}

	// create crawlers, thread count was defined in conf
	scheduler.crawlers = make([]*Crawler, 0)
	for i := 0; i < conf.TreadCount; i++ {
		crawler := NewCrawler(scheduler.urlTable, scheduler.config, &scheduler.queue)
		scheduler.crawlers = append(scheduler.crawlers, crawler)
	}

	return scheduler, nil
}

// run mini spider
func (ms *Scheduler) Run() {
	// start all crawlers
	// //TODO:yangxu
	// crawler := ms.crawlers[0];
	// crawler.Run()


	for _, crawler := range ms.crawlers {
		go crawler.Run()
	}
}

// get number of unfinished task
func (ms *Scheduler) GetUnfinished() int {

	return ms.queue.GetUnfinished()
}
