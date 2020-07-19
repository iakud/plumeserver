package table

import (
	"testing"
)

var testDatas = []*testData{
	&testData{1, 12, 1, 0},
	&testData{2, 12, 2, 50},
	&testData{3, 12, 3, 100},
	&testData{4, 12, 4, 150},
	&testData{5, 11, 1, 0},
	&testData{6, 11, 2, 50},
	&testData{7, 11, 3, 100},
	&testData{8, 11, 4, 150},
	&testData{9, 15, 1, 0},
	&testData{10, 15, 2, 50},
	&testData{11, 15, 3, 100},
	&testData{12, 15, 4, 150},
	&testData{13, 13, 1, 0},
	&testData{14, 13, 2, 50},
	&testData{15, 13, 3, 100},
	&testData{16, 13, 4, 150},
	&testData{17, 14, 1, 0},
	&testData{18, 14, 2, 50},
	&testData{19, 14, 3, 100},
	&testData{20, 14, 4, 150},
	&testData{21, 22, 1, 0},
	&testData{22, 22, 2, 50},
	&testData{23, 22, 3, 100},
	&testData{24, 22, 4, 150},
	&testData{25, 21, 1, 0},
	&testData{26, 21, 2, 50},
	&testData{27, 21, 3, 100},
	&testData{28, 21, 4, 150},
	&testData{29, 25, 1, 0},
	&testData{30, 25, 2, 50},
	&testData{31, 25, 3, 100},
	&testData{32, 25, 4, 150},
	&testData{33, 23, 1, 0},
	&testData{34, 23, 2, 50},
	&testData{35, 23, 3, 100},
	&testData{36, 23, 4, 150},
	&testData{37, 24, 1, 0},
	&testData{38, 24, 2, 50},
	&testData{39, 24, 3, 100},
	&testData{40, 24, 4, 150},
	&testData{41, 32, 1, 0},
	&testData{42, 32, 2, 50},
	&testData{43, 32, 3, 100},
	&testData{44, 32, 4, 150},
	&testData{45, 31, 1, 0},
	&testData{46, 31, 2, 50},
	&testData{47, 31, 3, 100},
	&testData{48, 31, 4, 150},
	&testData{49, 35, 1, 0},
	&testData{50, 35, 2, 50},
	&testData{51, 35, 3, 100},
	&testData{52, 35, 4, 150},
	&testData{53, 33, 1, 0},
	&testData{54, 33, 2, 50},
	&testData{55, 33, 3, 100},
	&testData{56, 33, 4, 150},
	&testData{57, 34, 1, 0},
	&testData{58, 34, 2, 50},
	&testData{59, 34, 3, 100},
	&testData{60, 34, 4, 150},
}

type testData struct {
	Id    int32
	Line  int32
	Level int32
	Param int32
}

func (t *testData) GetId() int32 {
	return t.Id
}

// 通过slice和指定key函数生成map
func TestKey(t *testing.T) {
	dataMap := make(map[int32]*testData)
	err := SliceToMap(&testDatas, &dataMap, WithKey((*testData).GetId))
	if err != nil {
		t.Fatal(err)
	}

	data := dataMap[4]
	t.Logf("TestKey: [id=4]: line=%v, level=%v, param=%v\n", data.Line, data.Level, data.Param)
}

// 通过slice和指定key函数生成具有组合键的map
func TestMultiKey(t *testing.T) {
	type Key struct {
		Line  int32
		Level int32
	}
	dataMap := make(map[Key]*testData)
	// 获取key的函数
	keyFunc := func(t *testData) Key {
		return Key{t.Line, t.Level}
	}
	err := SliceToMap(&testDatas, &dataMap, WithKey(keyFunc))
	if err != nil {
		t.Fatal(err)
	}

	data := dataMap[Key{14, 3}]
	t.Logf("TestMultiKey: [line=4,level=3], id=%v, param=%v\n", data.Id, data.Param)
}

// testData的Key方法
func (t *testData) Key() int32 {
	return t.Id
}

// 实现Key方法通过slice和生成map
func TestMethodKey(t *testing.T) {
	// map的key类型必须和Key()方法返回的实际类型一致
	dataMap := make(map[int32]*testData)

	err := SliceToMap(&testDatas, &dataMap)
	if err != nil {
		t.Fatal(err)
	}

	data := dataMap[17]
	t.Logf("TestInterfaceKey: [id=17]: line=%v, level=%v, param=%v\n", data.Line, data.Level, data.Param)
}
