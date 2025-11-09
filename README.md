
## twitter_cti

## Introduction 



Twitter IOC hunter written in golang. Which is based on the twitter-scarpper[https://github.com/n0madic/twitter-scraper] package of golang. Currently this tool parses the IP,url[Defang and fang both] and hashes.

### installtion 
Step 0: ```go mod init twitter_cti```

Step 1: ```go get -u github.com/n0madic/twitter-scraper```

Step 2: ```go build -o twitter_cti twitter_cti_v1.go```


### Usage 

To get an output as the files you need to specify the path with `-p` parameter which is writable. You can specify number of tweets to look by `-n` option

eg. 

`./twitter_cti -H "#emotet" -p ~/Desktop/go_lang/empty/ -n 100`


### TO DO 
1. Add default list of keyword to lookup 
2. Add top CTI Twitter handels
3. Add block for IOC enrichment

### Limitations 
1. Unable to parse domains 
2. Sometime false positive info is captures 
