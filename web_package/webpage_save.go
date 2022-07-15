package web_package


import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

const (
	OutputFileMode = 0644
)

func pathExit(rootPath string) Bool {
	if len(rootPath) {
		return false;
	}
	_ , error := os.Stat(rootPath)
	if !error  {
		return true;
	}
	return false;
}


/*
generate file path for given url
Params:
	- url: url to crawl
	- rootPath: root path for saving file
Returns:
	- file path
*/
func genFilePath(urlStr, rootPath string) string {
	filePath := url.QueryEscape(urlStr)
	filePath = path.Join(rootPath, filePath)
	return filePath
}


func SaveWebPage(rootPath string, url string, data []byte) error {
	if len(rootPath) {
		return fmt.Errorf("路径为空");
	}
		// create root dir, if not exist
	if !pathExit(rootPath) {
		os.MkdirAll(rootPath, 0777)
	}
	// generate full file path
	filePath := genFilePath(url, rootPath)

	// save to file
	err := ioutil.WriteFile(filePath, data, OutputFileMode)
	if err != nil {
		return fmt.Errorf("ioutil.WriteFile(%s):%s", filePath, err.Error())
	}

	return nil
}