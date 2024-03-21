package ktest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/khan-lau/kutils/file_format/efile"
)

func Test_ParseDocumentHeader(t *testing.T) {
	//<! Entity=华东 type=测试2011-11-03 dataTime='20120411 11:12:14' !>
	str := "<! Entity=华东 type=测试2011-11-03 dataTime='20120411 11:12:14' !>// fuck off"
	pAttributes, err := efile.ParseDocumentHeader(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	for e := pAttributes.Front(); e != nil; e = e.Next() {
		fmt.Printf("%-v\n", e.Value)
	}
}

func Test_ParseNodeLine(t *testing.T) {
	//<test::华东 DDMM='华东电网' date='2012-04-11 11:12' >
	// str := "<test::华东 DDMM='华东电网' date='2012-04-11 11:12' / >  //test"
	str := "<DataBlock NameTag=\"DG\"  DateTag=\"date\">"
	// str := "<sec/>"
	name, isEnd, pAttributes, err := efile.ParseNodeLine(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	fmt.Printf("name: %s, isEnd: %v\n", name, isEnd)

	for e := pAttributes.Front(); e != nil; e = e.Next() {
		fmt.Printf("%-v\n", e.Value)
	}

	fmt.Printf("\n")

	str = "<test::华东 DDMM='华东电网' date='2012-04-11 11:12'>  //test"
	name, isEnd, pAttributes, err = efile.ParseNodeLine(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	fmt.Printf("name: %s, isEnd: %v\n", name, isEnd)

	for e := pAttributes.Front(); e != nil; e = e.Next() {
		fmt.Printf("%-v\n", e.Value)
	}
}

func Test_ParseEText(t *testing.T) {
	str := "#1 '花花电网' 花花电网 '2011-11-03 00:00:02.0' 32  ''  "
	str = "#2 发生时间 '2011-11-03 00:00:02.0'"
	// str = "#1 单位名称 花花电网 花花电网 花花电网 花花电网 "
	str = "#2 发生时间 '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' "
	// str = "#3 次数 32 32 32 32 "

	delim := " "
	if strings.LastIndex(str, "\t") != -1 {
		delim = "\t"
	}

	fields, err := efile.ParseEText(str, delim)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	idx := 1
	for e := fields.Front(); e != nil; e = e.Next() {
		fmt.Printf("index: %d, field: %s\n", idx, e.Value)
		idx++
	}
}

func Test_ParseETable(t *testing.T) {
	str := `
	<! Entity=华东 type=测试2011-11-03 dataTime='20120411 11:12:14' !>
	<test::华东 DDMM='华东电网' date='2012-04-11 11:12' >
		@顺序 单位名称 发生时间 次数 
		#1 花花电网 无 1000 
		#2 花花电网 无 1000 
		#3 花花电网 无 1000 
		#4 花花电网 无 1000 
		#5 花花电网 无 1000 
		#6 花花电网 无 1000 
		#7 花花电网 无 1000 
		#8 花花电网 无 1000 
		#9 花花电网 无 1000 
		#10 花花电网1 无 1000 
		#11 花花电网 无 无 
		#12 花花电网 无 无 
		#13 花花电网 无 无 
		#14 花花电网 无 无 
		#15 花花电网 无 无 
		#16 花花电网 无 无 
		#17 花花电网 无 无 
		#18 花花电网 无 无 
		#19 花花电网 无 无 
		#20 花花电网1 无 无 
		#21 花花电网1 无 无 
		#22 花花电网2 无 无 
		#23 花花电网33 无 无 
		#24 花花电网as 无 无 
		#25 花花电网 无 无 
		#26 花花电网 无 无 
		#27 花花电网 无 无 
		#28 花花电网 无 无 
		#29 花花电网 无 11111 
		#30 花花电网 无 无 
		#31 花花电网 无 无 
		#32 花花电网 无 无 
		#33 花花电网 无 无 
		#34 花花电网 无 无 
		#35 花花电网 无 无 
		#36 花花电网 无 无 
		#37 花花电网 无 无 
		#38 花花电网 无 无 
		#39 花花电网1 无 无 
		#40 花花电网 无 无 
		#41 花花电网 无 无 
		#42 花花电网 无 无 
		#43 花花电网 无 无 
		#44 花花电网 无 无 
		#45 花花电网 无 无 
		#46 花花电网 无 无 
		#47 花花电网 无 无 
		#48 花花电网 无 无 
		#49 花花电网1 无 无 
		#50 花花电网1 无 无 
		#51 花花电网2 无 无 
		#52 花花电网33 无 无 
		#53 花花电网as 无 无 
		#54 花花电网 无 无 
		#55 花花电网 无 无 
		#56 花花电网 无 无 
		#57 花花电网 无 无 
		#58 花花电网 无 无 
	</test::华东>	
	`
	fmt.Println("横表式:")
	root, err := efile.ParseRootString(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	records, err := efile.ParseETable(root, "test::华东")
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	for eRow := records.Front(); eRow != nil; eRow = eRow.Next() {
		row := eRow.Value
		for eCol := row.Front(); eCol != nil; eCol = eCol.Next() {
			col := eCol.Value

			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}
	fmt.Println()

	str = `
	<! Entity=铁心桥 type=测试2011-11-03 dataTime='20120423 13:30:07' !>
	<DG::铁心桥 date='2012-04-23' DDMM='达梦' >
		@#顺序 属性名 '花花电网' '花花电网' '花花电网' '花花电网' 
		#1 单位名称 花花电网 花花电网 花花电网 花花电网 
		#2 发生时间 '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' 
		#3 次数 32 32 32 32 
	</DG::铁心桥>
	`
	fmt.Println("多列式:")
	root, err = efile.ParseRootString(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	records, err = efile.ParseETable(root, "DG::铁心桥")
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	for eRow := records.Front(); eRow != nil; eRow = eRow.Next() {
		row := eRow.Value
		for eCol := row.Front(); eCol != nil; eCol = eCol.Next() {
			col := eCol.Value

			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}
	fmt.Println()

	str = `
	<! Entity=铁心桥 type=测试2011-11-03 dataTime='20120423 13:34:20' !>
	<DG::铁心桥 date='2012-04-23' DDMM='达梦' >
		@@顺序 属性名 属性值
		#1 单位名称 花花电网 
		#2 发生时间 '2011-11-03 00:00:02.0' 
		#3 次数 32 
		
		#1 单位名称 花花电网 
		#2 发生时间 '2011-11-03 00:00:02.0' 
		#3 次数 32 
		
		#1 单位名称 花花电网 
		#2 发生时间 '2011-11-03 00:00:02.0' 
		#3 次数 32 
		
		#1 单位名称 花花电网 
		#2 发生时间 '2011-11-03 00:00:02.0' 
		#3 次数 32 
	</DG::铁心桥>
	`
	fmt.Println("单列式:")
	root, err = efile.ParseRootString(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	records, err = efile.ParseETable(root, "DG::铁心桥")
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	for eRow := records.Front(); eRow != nil; eRow = eRow.Next() {
		row := eRow.Value
		for eCol := row.Front(); eCol != nil; eCol = eCol.Next() {
			col := eCol.Value

			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}
	fmt.Println()
}

func Test_ParseRootBytes(t *testing.T) {
	str := `
	<! Entity=铁心桥 type=测试2011-11-03 dataTime='20120423 13:30:07' !>
	<DG::铁心桥 date='2012-04-23' DDMM='达梦' >
	</DG::铁心桥>
	
	
	
	<DataBlock NameTag="DG"  DateTag="date">
		<Data type="@#" col3="time" col4="value@DTNXJK:HSBFC:CDQ" col5="value@DTNXJK:HSBFC:Thr">
			<Sec1>
				<Row1>
				</Row1>
				<Row>
				</Row>
			</Sec>
			<Sec>
				<Row>
				</Row>
				<Row>
				</Row>
			</Sec>
		</Data>
	</DataBlock>
	`

	buf := bytes.NewBufferString(str)
	root, err := efile.ParseRootBytes(buf)

	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	jsonStr := root.ToJson()

	// fmt.Printf("%s\n", jsonStr)

	fmt.Printf("\n")

	v := make(map[string]any)
	err = json.Unmarshal([]byte(jsonStr), &v)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	//data := v.(map[string]interface{})
	jsonStr1, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	fmt.Printf("%s\n", string(jsonStr1))
}

func Test_ParseRootEFile(t *testing.T) {
	filePath := "C:\\Private\\Test\\elanguage\\data\\横表式.txt"
	root, err := efile.ParseRootEFile(filePath)

	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	records, err := efile.ParseETable(root, "test::华东")
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	for eRow := records.Front(); eRow != nil; eRow = eRow.Next() {
		row := eRow.Value
		for eCol := row.Front(); eCol != nil; eCol = eCol.Next() {
			col := eCol.Value

			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}
	fmt.Println()
}

func Test_GetENodeByPath(t *testing.T) {
	str := `
	<! Entity=铁心桥 type=测试2011-11-03 dataTime='20120423 13:30:07' !>
	<DG::铁心桥 date='2012-04-23' DDMM='达梦' >
		@#顺序 属性名 '花花电网' '花花电网' '花花电网' '花花电网' 
		#1 单位名称 花花电网 花花电网 花花电网 花花电网 
		#2 发生时间 '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' 
		#3 次数 32 32 32 32 
	</DG::铁心桥>
	
	
	
	<DataBlock NameTag="DG"  DateTag="date">
		<Data type="@#" col3="time" col4="value@DTNXJK:HSBFC:CDQ" col5="value@DTNXJK:HSBFC:Thr">
			<Sec1>
				<Row1>
					@#顺序 属性名 '花花电网' '花花电网' '花花电网' '花花电网' 
					#1 单位名称 花花电网 花花电网 花花电网 花花电网 
					#2 发生时间 '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' 
					#3 次数 32 32 32 32 
				</Row1>
				<Row>
				</Row>
			</Sec>
			<Sec>
				<Row>
				</Row>
				<Row>
				</Row>
			</Sec>
		</Data>
	</DataBlock>
	`

	buf := bytes.NewBufferString(str)
	root, err := efile.ParseRootBytes(buf)

	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	node, err := efile.GetENodeByPath(root, "DataBlock/Data/Sec1/Row1")
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	fmt.Printf("%s\n", node.Value)
}
