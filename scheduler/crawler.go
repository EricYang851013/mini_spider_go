package scheduler

import (
	"fmt"
	"regexp"
	"time"
	"net/http" 
	"io/ioutil" 
	"log"
)





import (
	"mini_spider/web_package"
)

// 请求方法
const (
	GET  = "GET"
	POST = "POST"
	HEAD = "HEAD"
	PUT  = "PUT"
)

type Crawler struct {
	urlTable   *UrlTable  //已抓取url
	config     *web_package.Config
	queue      *Queue
	urlPattern *regexp.Regexp
	stop       bool
}


// create new crawler
func NewCrawler(urlTable *UrlTable, config *web_package.Config, queue *Queue) *Crawler {
	c := new(Crawler)
	c.urlTable = urlTable
	c.config = config
	c.queue = queue

	// TargetUrl has been checked in conf load
	c.urlPattern, _ = regexp.Compile(c.config.TargetUrl)
	

	c.stop = false

	return c
}

func ReadGet(url string, timeout int, header map[string]string) ([]byte,error)  {
	request, err := http.NewRequest("GET",url,nil) 

	if err != nil {
		return nil, err;
	}

	
	params := request.URL.Query();
	for k, v := range header {
		params.Add(k,v);
	}
	request.URL.RawQuery = params.Encode();

	client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
    }
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("request error: url=%s, error=%d", url, err)
		return nil, err;
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body,nil;

}

func ReadPost(url string, timeout int, header map[string]string) ([]byte,error) {
	client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
    }
	request, err := http.NewRequest("POST",url,nil) 
	if err != nil {
        return nil, err
	}
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	for k, v := range header {
		request.Header.Set(k,v);
	}
	
	response,err := client.Do(request) 
	if err != nil {
        return nil, err;
	}

	resByte,err := ioutil.ReadAll(response.Body)
	return resByte, err;
}

// start crawler
func (c *Crawler) Run() {
	for !c.stop {
		// get new task from queue
		task := c.queue.Pop()
		log.Printf("from queue: url=%s, depth=%d", task.Url, task.Depth)

		// read data from given task
		data, err := ReadGet(task.Url, c.config.CrawlTimeout, task.Header)
		if err != nil {
			log.Printf("http_util.Read(%s):%s", task.Url, err.Error())

			c.queue.FinishOneTask()
			continue
		}

		// save data to file
			err = web_package.SaveWebPage(c.config.OutputDirectory, task.Url, data)
			if err != nil {
				log.Printf("web_package.SaveWebPage(%s):%s", task.Url, err.Error())
			} else {
				log.Printf("save to file: %s", task.Url)
			}
		

		// add to url table
		c.urlTable.Add(task.Url)

		// continue crawling until max depth
		if task.Depth < c.config.MaxDepth {
			err = c.crawlChild(data, task)
			if(err != nil){
				log.Printf("crawlChild(%s):%s in depth of %d", task.Url, err.Error(), task.Depth)

			}
		}

		// confirm to remove task from queue
		c.queue.FinishOneTask()

		// sleep for a while
	    time.Sleep(time.Duration(c.config.CrawlInterval) * time.Second)
	}
}


// stop crawler
func (c *Crawler) Stop() {
	c.stop = true
}


// crawl child url
func (c *Crawler) crawlChild(data []byte, task *CrawlTask ) error {
	// parse url from web page
	links, err := ParseWebPage(data, task.Url,c.urlPattern)
	if err != nil {
		return fmt.Errorf("web_package.ParseWebPage():%s", err.Error())
	}

	// add child task to queue
	for _, link := range links {
		// check whether url match the pattern, or url exists already
		if c.urlTable.Exist(link) {
			continue
		}

		taskNew := &CrawlTask{Url: link, Depth: task.Depth + 1, Header: make(map[string]string)}
		c.queue.Add(taskNew)
	}

	return nil
}