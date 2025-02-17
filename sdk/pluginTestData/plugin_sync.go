package pluginTestData

import (
	"github.com/brokercap/bifrost-core/sdk/plugin/driver"
	"math/rand"
	"time"
	"fmt"
	"log"
)

type Plugin struct {
	pluginObj driver.ConnFun
	param map[string]interface{}
	pluginParamObj interface{}
	err error
	eventType EventType
	debug bool
}

func NewPlugin(name string,url string) *Plugin {
	return &Plugin{
		pluginObj:driver.Open(name,url),
		param:nil,
		pluginParamObj:nil,
		eventType:-1,
		debug:true,
	}
}

func (This *Plugin) SetParam(m map[string]interface{}) error {
	This.param = m
	This.pluginParamObj,This.err = This.pluginObj.SetParam(This.param)
	if This.err != nil{
		return This.err
	}
	return nil
}

func (This *Plugin) SetEventType(eventType EventType) {
	This.eventType = eventType
}

func (This *Plugin) SetDebug(b bool) {
	This.debug = b
}

func (This *Plugin) DoTestStart(n uint)  error{
	if This.param == nil{
		return fmt.Errorf("SetParam please!")
	}
	var i uint = 0
	e := NewEvent()
	e.SetSaveHistory(false)
	var intN EventType
	var startTime = time.Now().UnixNano()

	log.Println("startTime:",startTime)

	defer func() {
		NowTime := time.Now().UnixNano()
		log.Println("startTime:",startTime," overTime:",NowTime, " use:",(NowTime-startTime)/1000000,"ms")
	}()
	for{
		if n > 0 && n <= i{
			break
		}
		i++

		if This.pluginParamObj == nil{
			This.pluginParamObj,This.err = This.pluginObj.SetParam(This.param)
		}else{
			This.pluginParamObj,This.err = This.pluginObj.SetParam(This.pluginParamObj)
		}

		if This.err != nil{
			return This.err
		}

		var data *driver.PluginDataType
		if This.eventType > -1{
			intN = This.eventType
		}else{
			rand.Seed(time.Now().UnixNano()+int64(i))
			intN = EventType(rand.Intn(5))
		}

		switch intN {
		case 0:
			data = e.GetTestInsertData()
			break
		case 1:
			data = e.GetTestUpdateData()
			break
		case 2:
			data = e.GetTestDeleteData()
			break
		case 3:
			data = e.GetTestQueryData()
			break
		default:
			break
		}

		var binlog *driver.PluginBinlog
		var err error
		var n0 int = 0
		var opName string
		for {
			switch intN {
			case 0:
				opName = "insert"
				binlog, err = This.pluginObj.Insert(data)
				break
			case 1:
				opName = "update"
				binlog, err = This.pluginObj.Update(data)
				break
			case 2:
				opName = "delete"
				binlog, err = This.pluginObj.Del(data)
				break
			case 3:
				opName = "sql"
				binlog, err = This.pluginObj.Query(data)
				break
			default:
				opName = "commit"
				if n == 0 {
					time.Sleep(1)
				}
				data = nil
				binlog, err = This.pluginObj.Commit()
				break
			}

			if err == nil{
				if This.debug {
					log.Println("success(", i, ") ", opName, " binlog:", binlog, " data:", data)
				}
				break
			}

			log.Println("err(",n0,") ",opName,":",err," data:",data)
			n0++

			if n0 == 60{
				return err
			}
		}
	}

	This.pluginObj.Commit()

	return nil
}

//用于性能测试。必须指定eventType,不支持debug
func (This *Plugin) DoTestStartForSpeed(n uint)  error{
	if This.param == nil{
		return fmt.Errorf("SetParam please!")
	}
	var i uint = 0
	e := NewEvent()
	e.SetSaveHistory(false)

	switch This.eventType {
	case INSERT,UPDATE,DELETE,SQLTYPE:
		break
	default:
		return fmt.Errorf("evnentType error,must be 0,1,2,3,SetEventType(n)")
	}

	var intN EventType

	intN = This.eventType
	var data *driver.PluginDataType
	switch intN {
	case 0:
		data = e.GetTestInsertData()
		break
	case 1:
		data = e.GetTestUpdateData()
		break
	case 2:
		data = e.GetTestDeleteData()
		break
	case 3:
		data = e.GetTestQueryData()
		break
	default:
		break
	}

	var startTime = time.Now().UnixNano()
	var err error

	log.Println("startTime:",startTime)

	defer func() {
		NowTime := time.Now().UnixNano()
		log.Println("startTime:",startTime," overTime:",NowTime, " use:",(NowTime-startTime)/1000000,"ms")
	}()
	for{
		if n > 0 && n <= i{
			break
		}
		i++

		if This.pluginParamObj == nil{
			This.pluginParamObj,This.err = This.pluginObj.SetParam(This.param)
		}else{
			This.pluginParamObj,This.err = This.pluginObj.SetParam(This.pluginParamObj)
		}

		if This.err != nil{
			return This.err
		}

		switch intN {
		case INSERT:
			_, err = This.pluginObj.Insert(data)
			break
		case UPDATE:
			_, err = This.pluginObj.Update(data)
			break
		case DELETE:
			_, err = This.pluginObj.Del(data)
			break
		case SQLTYPE:
			_, err = This.pluginObj.Query(data)
			break
		default:
			break
		}

		if err != nil{
			return fmt.Errorf("err:"+fmt.Sprint(err)+" data:"+fmt.Sprint(data))
		}
	}

	This.pluginObj.Commit()

	return nil
}