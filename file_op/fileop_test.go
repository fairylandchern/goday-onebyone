package file_op

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"os"
	"testing"
)

// 参考网络io中半包和协议的方式，处理文件读写相关操作demo
// lenth(8byte):data
func TestWriteAndReadLogFile(t *testing.T) {
	f, err := os.OpenFile("./test_log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatal("open file err:", err)
		return
	}
	defer f.Close()
	_ = f

	// write file,write file and set the file meta data.
	//w := bufio.NewWriter(f)
	//for i := 0; i < 100; i++ {
	//	d := &WData{fmt.Sprintf("%v%d", "helloworld", i)}
	//	data, err := json.Marshal(d)
	//	if err != nil {
	//		t.Log("json marshal data err:", err)
	//		continue
	//	}
	//	lenth := len(data)
	//	dataLen := setDataLenth(uint64(lenth))
	//	t.Log("lenth:",dataLen)
	//	dataLen=append(dataLen,data...)
	//	count,err:=w.Write(dataLen)
	//	t.Logf("count:%d,err:%v",count,err)
	//	w.Flush()
	//}
	// read file, read file according to meta data.
	r := bufio.NewReader(f)
	//for r.Size()>0{
	lenth := make([]byte, 8, 8)
	count, err := r.Read(lenth)
	if err != nil {
		t.Fatal("err read file:", err)
		//break
	}
	// 25或者26对应的ascii码为em，或者sub，无法打印出来
	// 因此在日志文件中无法展示出来，但是可以在读文件的时候固定大小
	// 从而确定将要读取的文件行数
	datalen := binary.LittleEndian.Uint64(lenth)
	t.Logf("count:%v,lenth:%v,lenth:%v", count, binary.LittleEndian.Uint64(lenth), lenth)
	data := make([]byte, datalen)
	count, err = r.Read(data)
	d := &WData{}
	err = json.Unmarshal(data, d)
	t.Logf("data:%v,err:%v,wdata:%v", data, err, d)
	//}
}

type WData struct {
	Content string
}

func TestDataLenth(t *testing.T) {
	lenth := 888
	data := setDataLenth(uint64(lenth))
	t.Log("lenth:", lenth, " data:", data)
	nl := binary.LittleEndian.Uint64(data)
	t.Log("new lenth:", nl)
}

func setDataLenth(lenth uint64) []byte {
	data := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(data, lenth)
	return data
}
