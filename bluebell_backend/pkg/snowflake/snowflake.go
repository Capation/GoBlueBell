package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

// 雪花算法

var node *sf.Node

// 需传入当前的机器ID
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	//fmt.Println(st)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID) // 指定一下机器的ID
	return
}

// GetID生成ID
func GenID() int64 {
	return node.Generate().Int64()
}

//func main() {
//	if err := Init("2020-07-01", 1); err != nil {
//		fmt.Printf("Init failed, err:%v\n", err)
//		return
//	}
//	id := GetID()
//	fmt.Println(id)
//}
