package accept

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/sesehai/qqwry"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//专门放些杂七杂八的东西

var filterIp = []string{"172.31.", "127.0.0.1"}

const (
	CONNECT_ = "CONNECT"
	//qqwry这个IP库需要定时更新
	QQWRYFILE = "/home/atlaswu/go/src/generproxy/tools/qqwry.dat"
)

//过滤链接
//返回值: true 需要过滤
func filterReq(ipAddr string) bool {
	var bRet = false

	for _, v := range filterIp {
		if strings.Contains(ipAddr, v) {
			bRet = true
			break
		}
	}

	return bRet
}

//获取用户地理位置
//返回值: 位置信息 运营商
func getUserGeolocation(address net.Addr) (string, string) {
	var ip = address.String()
	var qqwryfile = QQWRYFILE

	ip_port := strings.Split(ip, ":")
	file, _ := qqwry.Getqqdata(qqwryfile)
	country, location := qqwry.Getlocation(file, ip_port[0])

	return country, location
}

//编码转换器 GBK ->UTF8
//返回值:
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func writeFile(s string) {
	filePath := "./generTest/golang.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(s)
	write.Flush()
}
