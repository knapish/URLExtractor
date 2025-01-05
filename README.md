# URLExtractor
 A Go library for extracting and parsing URL components.


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

### Input: `http://ö.example.com:80/path`
```
Scheme: http
UserInfo: 
Subdomain: xn--nda
Apex: example
TLD: com
Port: 80
Path: /path
Query: 
Fragment: 
```

### Input: `http://xn--nda.test.example.com:81`
```
Scheme: http
UserInfo: 
Subdomain: xn--nda.test
Apex: example
TLD: com
Port: 81
Path: 
Query: 
Fragment: 
```

### Input: `www.example.com`
```
Scheme: 
UserInfo: 
Subdomain: www
Apex: example
TLD: com
Port: 0
Path: 
Query: 
Fragment: 
```

### Input: `_valid.example.com`
```
Scheme: 
UserInfo: 
Subdomain: _valid
Apex: example
TLD: com
Port: 0
Path: 
Query: 
Fragment: 
```

### Input: `x.com`
```
Scheme: 
UserInfo: 
Subdomain: 
Apex: x
TLD: com
Port: 0
Path: 
Query: 
Fragment: 
```

### Input: `xcom`
```
Error: invalid domain - missing TLD in string: xcom
```

### Input: `-invalid.example.com`
```
Error: invalid domain - invalid sub domain: -invalid
```

### Input: `invalid._example.com`
```
Error: invalid domain - invalid apex domain: _example
```

### Input: `invalid`
```
Error: invalid domain - missing TLD in string: invalid
```

### Input: `ööö`
```
Error: invalid domain - missing TLD in string: invalid
```

### Input: `example.example`
```
Input: example.example
Error: invalid domain - missing valid TLD in string: example.example
```
