# mcurl

<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Installation](#installation)
- [How to use - help](#how-to-use-help)
- [Example](#example)
    - [Command](#command)
    - [Output](#output)
- [TODO](#todo)

<!-- /TOC -->

## Installation

```shell script
go install -v github.com/ninjawule/mcurl
```

[Top](#mcurl)

## How to use - help

```shell script
mcurl -h
```

[Top](#mcurl)

## Example

### Command

```shell script
mcurl -v -X POST --url https://example.com -H Content-Type:application/json -H Authorization:Bearer\ xMLdsjhryq67snqa089FX -r 3 \{\"key\":\"value\"\} -n 5
```

Explanation :

- we're activating the verbose mode with `-v`
- we're requesting with HTTP method `POST`
- working with url `https://example.com`
- setting 2 header parameters with `-H` used twice here
- launching 5 identical requests in parallel with `-n 5`
- sending data with `-d` (or from a file with flag `-i`)
- doing all this **3 times**, for better average times computation, with `-r 3`

[Top](#mcurl)

### Output

```
2020/11/05 09:31:04 -----------------------------------------------------------------------------------------------
2020/11/05 09:31:04 --- STARTING the testing of 5 'parallel' requests @ https://example.com (3 times)
2020/11/05 09:31:04 -----------------------------------------------------------------------------------------------
2020/11/05 09:31:04 Running options =
method  : POST
url     : https://example.com
nbReqs  : 5
repeat  : 3
data    : {"key":"value"}
headers : Content-Type: application/json, Authorization: Bearer K0i79im-QAmeRe9wASjzpA
2020/11/05 09:31:04 --- Doing try n°1
2020/11/05 09:31:04 Starting request 5
2020/11/05 09:31:04 Starting request 2
2020/11/05 09:31:04 Starting request 3
2020/11/05 09:31:04 Starting request 1
2020/11/05 09:31:04 Starting request 4
2020/11/05 09:31:04 Request 1 done in 504.988041ms
2020/11/05 09:31:04 Request 5 done in 505.236564ms
2020/11/05 09:31:04 Request 3 done in 505.31119ms
2020/11/05 09:31:04 Request 2 done in 505.933135ms
2020/11/05 09:31:04 Request 4 done in 506.504041ms
2020/11/05 09:31:04 --- --> Done in 506.704702ms. Average / Max request time: 505.594594ms / 506.504041ms
2020/11/05 09:31:05 --- Doing try n°2
2020/11/05 09:31:05 Starting request 5
2020/11/05 09:31:05 Starting request 2
2020/11/05 09:31:05 Starting request 3
2020/11/05 09:31:05 Starting request 1
2020/11/05 09:31:05 Starting request 4
2020/11/05 09:31:06 Request 5 done in 149.238699ms
2020/11/05 09:31:06 Request 1 done in 149.201297ms
2020/11/05 09:31:06 Request 2 done in 149.523666ms
2020/11/05 09:31:06 Request 4 done in 150.656459ms
2020/11/05 09:31:06 Request 3 done in 151.896139ms
2020/11/05 09:31:06 --- --> Done in 152.061727ms. Average / Max request time: 150.103252ms / 151.896139ms
2020/11/05 09:31:07 --- Doing try n°3
2020/11/05 09:31:07 Starting request 5
2020/11/05 09:31:07 Starting request 2
2020/11/05 09:31:07 Starting request 3
2020/11/05 09:31:07 Starting request 1
2020/11/05 09:31:07 Starting request 4
2020/11/05 09:31:07 Request 5 done in 142.497641ms
2020/11/05 09:31:07 Request 3 done in 143.873313ms
2020/11/05 09:31:07 Request 1 done in 143.868609ms
2020/11/05 09:31:07 Request 2 done in 144.226347ms
2020/11/05 09:31:07 Request 4 done in 144.551647ms
2020/11/05 09:31:07 --- --> Done in 144.74459ms. Average / Max request time: 143.803511ms / 144.551647ms
2020/11/05 09:31:08 -----------------------------------------------------------------------------------------------
2020/11/05 09:31:08 --- Done in 3.812915071s. Average / Max request time: 266.500452ms / 506.504041ms
2020/11/05 09:31:08 -----------------------------------------------------------------------------------------------
```

[Top](#mcurl)

## TODO

- being able to change the body for each request

[Top](#mcurl)
