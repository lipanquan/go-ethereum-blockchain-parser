package main

import (
  "bytes"
//  "fmt"
  "log"
  "strconv"
  "strings"
  "os/exec" 
  "math"
  "time"
  "runtime"
  "encoding/json"
  "os"
  "./lib"
//  "./go-ethereum/common"
//  "./go-ethereum/common/hexutil"
  //"./golang-set"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/common/hexutil"
//  "github.com/deckarep/golang-set"
)

var MAX_NUM_SHARE = 10
var MAX_NUM_RETRY = 3
var APPEND_FILE_NUM = 100
var c = make(chan string, 100)

type ResultOfGetBlockByNumber struct {
    Jsonrpc    string
    Id    int
    Result	Block
}
type RPCTransaction struct {
        BlockHash        common.Hash     `json:"blockHash"`
        BlockNumber      *hexutil.Big    `json:"blockNumber"`
        From             *common.Address  `json:"from"`
        Gas              hexutil.Uint64  `json:"gas"`
        GasPrice         *hexutil.Big    `json:"gasPrice"`
        Hash             common.Hash     `json:"hash"`
        Input            hexutil.Bytes   `json:"input"`
        Nonce            hexutil.Uint64  `json:"nonce"`
        To               *common.Address `json:"to"`
        TransactionIndex hexutil.Uint    `json:"transactionIndex"`
        Value            *hexutil.Big    `json:"value"`
        V                *hexutil.Big    `json:"v"`
        R                *hexutil.Big    `json:"r"`
        S                *hexutil.Big    `json:"s"`
}

type Block struct {
        UncleHashes  []common.Hash    `json:"uncles"`
        Hash         common.Hash      `json:"hash"`
        Timestamp    string      `json:"timestamp"`
        Transactions []RPCTransaction `json:"transactions"`
}

var tenToAny map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}
 
 
func decimalToAny(num, n int) string {
 new_num_str := ""
 var remainder int
 var remainder_string string
 for num != 0 {
  remainder = num % n
  if 76 > remainder && remainder > 9 {
   remainder_string = tenToAny[remainder]
  } else {
   remainder_string = strconv.Itoa(remainder)
  }
  new_num_str = remainder_string + new_num_str
  num = num / n
 }
 return new_num_str
}

func findkey(in string) int {
 result := -1
 for k, v := range tenToAny {
  if in == v {
   result = k
  }
 }
 return result
}
func anyToDecimal(num string, n int) float64 {
 var new_num float64
 new_num = 0.0
 nNum := len(strings.Split(num, "")) - 1
 for _, value := range strings.Split(num, "") {
  tmp := float64(findkey(value))
  if tmp != -1 {
   new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
   nNum = nNum - 1
  } else {
   break
  }
 }
 return float64(new_num)
}

func exec_shell(s string, blockNumber int, retry int, fileContent *bytes.Buffer) {
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer

    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        //log.Fatal("cmd="+s+" run fail!", err)
	retry--
        log.Print("Fatal!!!retry=",retry," cmd="+s+" run fail!", err)
	if(retry > 0) {
		exec_shell(s, blockNumber, retry, fileContent)
    	}
    }
    
    stb := &ResultOfGetBlockByNumber{}
    //fmt.Printf("%s\n", out.String())
    err = json.Unmarshal([]byte(out.String()), &stb)
    
    if err != nil {
	retry--;
        log.Print("umarshal fail!!!retry=",retry," block=",blockNumber," cmd=", s, " input=", out.String()," ", err)
	if(retry > 0) {
		exec_shell(s, blockNumber, retry, fileContent)
    	}
    } 
    //fmt.Printf("----result=%s:%d\n", stb.Jsonrpc, stb.Id)
    blockInfo := stb.Result
    if len(blockInfo.Transactions) <= 0 {
	//log.Print("block=", blockNumber, " transaction len <= 0")
	return
    }
    timeLayout := "2006-01-02-15-04-05"
    timestamp := blockInfo.Timestamp
    timestamp = timestamp[2 : len(timestamp)] // get rid of 0x prefix
    txTime := time.Unix(int64(anyToDecimal(timestamp, 16)), 0).Format(timeLayout)
    for _, value := range blockInfo.Transactions {
	    fileContent.WriteString(txTime + "\t")
	    fileContent.WriteString(value.Hash.String() + "\t")
	    if value.From != nil {
		fileContent.WriteString(value.From.String() + "\t")
	    } else {
		fileContent.WriteString("null\t")
	    }
		
	    if value.To != nil {
		fileContent.WriteString(value.To.String() + "\t")
	    } else {
		fileContent.WriteString("null\t")
	    }
	    txValue := value.Value.String()
	    txValue = txValue[2 : len(txValue)] // get rid of 0x prefix
            price := anyToDecimal(txValue, 16)/1000000000000000000
	    fileContent.WriteString(strconv.FormatFloat(price, 'f', -1, 64) + "\n")
	    //fmt.Printf("----block=%d transactionIndex=%s fr=%s, to=%s\n", blockNumber, value.TransactionIndex.String(), value.From.String(), value.To.String())
    }
    //log.Print(accountSet)
}

func pathExists(path string) (bool) {
        _, err := os.Stat(path)
        if err == nil {
                return true
        }
        return false
}

func getTxByBlock(fromBlockNumber int, toBlockNumber int, fileName string) {
	taskId := strconv.Itoa(fromBlockNumber) + "-" + strconv.Itoa(toBlockNumber)
	log.Print(taskId," begin: file=", fileName)
        fd,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
        if err != nil {
		log.Fatal(fileName, " open fail! ", err)
		return
        }
        fileContent := bytes.Buffer{} 
	i := 0
	for blockNumber := fromBlockNumber; blockNumber <= toBlockNumber; blockNumber++ {
  		blockNumberHex := "0x" + decimalToAny(blockNumber, 16);
    		command := "curl -X POST --data '{\"jsonrpc\":\"3.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"" + blockNumberHex + "\", true],\"id\":1}' -H \"Content-type: application/json;charset=UTF-8\"  localhost:8545";

    	 	exec_shell(command, blockNumber, MAX_NUM_RETRY, &fileContent)
		if i%APPEND_FILE_NUM == 0 {
			if len(fileContent.String()) > 0 {
                        	buf:=[]byte(fileContent.String())
                        	fd.Write(buf)
                        	fileContent = bytes.Buffer{}
				log.Print("write accounts into file=", fileName, " taskId=", taskId, ", i=",i)
			} else {
                        	fileContent = bytes.Buffer{}
				log.Print("fileContent len <= 0, not need to write accounts into file=", fileName, " taskId=", taskId, ", i=",i)
			}
                }
		i++
	}
	if len(fileContent.String()) > 0 {
                buf:=[]byte(fileContent.String())
                fd.Write(buf)
		log.Print("write accounts into file=", fileName, " taskId=", taskId, ", i=",i)
        }
        fd.Close() 
	c <- taskId
}

func main() {

  if len(os.Args) < 3 {
        log.Fatal("Param Invalid!!! go run getTxByTime.go [timeFrom] [timeTo], eg. go run getTxByTime.go 2018-01-01-00-00-00 2018-02-01-00-00-00")
  }
  log.Print("====getTxByTime begin==================");
  timeBegin := time.Now().Unix()
  MULTICORE := runtime.NumCPU()
  runtime.GOMAXPROCS(MULTICORE)
  //blockNumber := 4927600;
  timeFrom := os.Args[1]
  timeTo := os.Args[2]
  if timeFrom >= timeTo {
        log.Fatal("timeFrom=", timeFrom, " >= timeTo=", timeTo)
  }
  blockNumberBegin, blockNumberEnd := lib.GetBlockNumberByTime(timeFrom, timeTo)
  totalBlockNum := blockNumberEnd - blockNumberBegin + 1;
  share := totalBlockNum / MAX_NUM_SHARE + 1
  loopCount := 0
  dir := lib.GetAndCheckDir("tx")
  filePrefix := timeFrom + "-" + timeTo 
  files := dir + "/" + filePrefix + "-*"
  resultFile := dir + "/" + filePrefix
  resultFileTmp := dir + "/"+strconv.FormatInt(timeBegin, 10)+"tmp"
  lib.ExecCmd("rm " + files, false)
  lib.ExecCmd("rm " + resultFileTmp, false)
  for i := blockNumberBegin; i <= blockNumberEnd; i++ {
	from := i
	if ((share+i) <= blockNumberEnd) {
		i += share 
	} else {
		i = blockNumberEnd
	}
	to := i
	loopCount++
        fileName := dir + "/" + filePrefix  + "-"+ strconv.Itoa(loopCount)
	go getTxByBlock(from, to, fileName);
  }
  for i := 0; i < loopCount; i++ {
	taskId := <- c
	log.Print(taskId, " finish");
  }

  lib.ExecCmd( "cat "+files + " >> " + resultFileTmp, true)
  lib.ExecCmd("rm " + files, true)
  lib.ExecCmd("mv " + resultFileTmp + " " + resultFile, true)
  lib.ExecCmd("sort -k 3 " + resultFile + " > " + resultFile + "-from-sort", true)
  lib.ExecCmd("sort -k 4 " + resultFile + " > " + resultFile + "-to-sort", true)


  timeEnd := time.Now().Unix()  
  log.Print("getTxByBlock finish, cost=", (timeEnd - timeBegin), "s")
}
