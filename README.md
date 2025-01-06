# URLExtractor
 A Go library for extracting and parsing URL components.

 Returns result only if valid ICANN TLD.

 a-z0-9\-\_

 

## Usage

```
package main

import (
	"fmt"
	"os"

	"github.com/knapish/urlextractor"
)

func main() {
	testInput := os.Args[1]
	fmt.Println("Input:", testInput)
	url, err := urlextractor.Extract(testInput)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Scheme:", url.Scheme)
		fmt.Println("UserInfo:", url.UserInfo)
		fmt.Println("Subdomain:", url.SubDomain)
		fmt.Println("Apex:", url.ApexDomain)
		fmt.Println("TLD:", url.TLD)
		fmt.Println("Port:", url.Port)
		fmt.Println("Path:", url.Path)
		fmt.Println("Query:", url.Query)
		fmt.Println("Fragment:", url.Fragment)
	}
}
```

### Input: `https://user:pass@test.www.example.com:444/example/path?query=data#fragment`
```
Scheme: https
UserInfo: user:pass
Subdomain: test.www
Apex: example
TLD: com
Port: 444
Path: /example/path
Query: query=data
Fragment: fragment
```
