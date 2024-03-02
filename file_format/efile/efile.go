package efile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"khan/kutil/container/klists"
	"khan/kutil/container/kmaps"
	"khan/kutil/data"
)

var seqGen data.Generator = *data.NewGenerator(1)

type EAttribute struct {
	Name  string
	Value string
}

type ENode struct {
	Id     uint32
	parent *ENode
	Name   string
	Value  *bytes.Buffer

	Attribes map[string]string     //  []*EAttribute
	Children *klists.KList[*ENode] //  []*ENode

	isEnd bool
}

func NewNode(id uint32, parent *ENode, name string, attributes map[string]string, pChildren *klists.KList[*ENode]) *ENode {
	return &ENode{Id: id, parent: parent, Name: name, Attribes: attributes, Children: pChildren}
}

func (node *ENode) GetAttribute(name string) (string, error) {
	if node.Attribes == nil {
		return "", fmt.Errorf("attribute %s not found", name)
	}
	return node.Attribes[name], nil
}

func (node *ENode) AddAttribute(name string, value string) {
	if node.Attribes == nil {
		node.Attribes = make(map[string]string)
	}
	node.Attribes[name] = value
}

func (node *ENode) DelAttribute(name string) {
	if node.Attribes == nil {
		return
	}

	delete(node.Attribes, name)
}

func (node *ENode) AddChildren(children *ENode) {
	if node.Children == nil {
		node.Children = klists.New[*ENode]()
	}
	node.Children.PushBack(children)
}

func (node *ENode) DelChildrenById(id uint32) {
	if node.Children == nil || node.Children.Len() == 0 {
		return
	}

	// 遍历删除
	var next *klists.KElement[*ENode]
	for e := node.Children.Front(); e != nil; e = next {
		next = e.Next()

		if e.Value.Id == id {
			node.Children.Remove(e)
			break
		}
	}
}

func (node *ENode) hasChild() bool {
	if node.Children == nil || node.Children.Len() == 0 {
		return false
	}

	return true
}

func (node *ENode) DelChildrenByName(name string) {
	if node.Children == nil || node.Children.Len() == 0 {
		return
	}

	// 遍历删除
	var next *klists.KElement[*ENode]
	for e := node.Children.Front(); e != nil; e = next {
		next = e.Next()

		if e.Value.Name == name {
			node.Children.Remove(e)
			break
		}
	}
}

func (node *ENode) GetENodeByName(name string) *ENode {
	if node.Children == nil || node.Children.Len() == 0 {
		return nil
	}
	for e := node.Children.Front(); e != nil; e = e.Next() {
		child := e.Value
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (node *ENode) GetENodeById(id uint32) *ENode {
	if node.Children == nil || node.Children.Len() == 0 {
		return nil
	}
	for e := node.Children.Front(); e != nil; e = e.Next() {
		child := e.Value
		if child.Id == id {
			return child
		}
	}
	return nil
}

// 按name从tree结构中获取指定node
func (node *ENode) FindFirstNodeByName(name string) *ENode {
	if node.Children == nil || node.Children.Len() == 0 {
		return nil
	}

	for e := node.Children.Front(); e != nil; e = e.Next() {
		child := e.Value
		if child.Name == name {
			return child
		} else {
			if child.hasChild() {
				subChild := child.FindFirstNodeByName(name)
				if nil != subChild {
					return subChild
				}
			}
		}
	}
	return nil
}

// 按id从tree结构中获取指定node
func (node *ENode) FindFirstNodeById(id uint32) *ENode {
	if node.Children == nil || node.Children.Len() == 0 {
		return nil
	}

	for e := node.Children.Front(); e != nil; e = e.Next() {
		child := e.Value
		if child.Id == id {
			return child
		} else {
			if child.hasChild() {
				subChild := child.FindFirstNodeById(id)
				if nil != subChild {
					return subChild
				}
			}
		}
	}
	return nil
}

func (node *ENode) ToJson() string {
	buf := bytes.NewBufferString("")

	buf.WriteByte('{')

	buf.WriteString("\"id\": ")
	buf.WriteString(strconv.FormatUint(uint64(node.Id), 10))
	buf.WriteString(", \n")

	if len(node.Name) > 0 {
		buf.WriteString("\"name\": \"")
		buf.WriteString(node.Name)
		buf.WriteString("\", \n")
	}

	attribesSize := len(node.Attribes)
	if attribesSize > 0 {
		for k, v := range node.Attribes {
			buf.WriteString("\"")
			buf.WriteString(k)
			buf.WriteString("\": \"")

			buf.WriteString(v)
			buf.WriteString("\", \n")
		}
	}

	if node.Value != nil {
		buf.WriteString("\"value\": \"")
		buf.WriteString(node.Value.String())
		buf.WriteString("\", \n")
	}

	if nil != node.Children {
		buf.WriteString("\"children\": [")
		idx := 0
		for e := node.Children.Front(); e != nil; e = e.Next() {
			child := e.Value
			buf.WriteString(child.ToJson())
			if idx < node.Children.Len()-1 {
				buf.WriteString(", \n")
			}
			idx++
		}

		buf.WriteString("], \n")
	}

	buf.WriteString("\"isEnd\": ")
	buf.WriteString(strconv.FormatBool(node.isEnd))
	buf.WriteString("\n")

	buf.WriteByte('}')

	return buf.String()
}

////////////////////////////////////////////////////////////////////

func GetENodeByPath(root *ENode, path string) (*ENode, error) {
	path = strings.TrimSpace(path)

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	strArray := strings.Split(path, "/")
	node := root
	for _, name := range strArray {
		node = node.GetENodeByName(name)
		if nil == node {
			return nil, fmt.Errorf("node not found")
		}
	}
	if node.Id == root.Id {
		return nil, fmt.Errorf("node not found")
	}

	return node, nil
}

func ParseRootEFile(path string) (*ENode, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("file %s not found", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	length := fileInfo.Size()
	reader := bufio.NewReader(file)

	if length < 1 {
		return nil, fmt.Errorf("parse etext error, document is empty")
	}

	var firstErr error = nil
	var root *ENode = nil
	var parentNode *ENode = nil

	for {
		firstErr = nil
		line, err := reader.ReadString('\n')
		if nil != err {
			if io.EOF != err {
				firstErr = err
			}
			break
		}
		line = strings.TrimSpace(line)
		//忽略行注释 与 空行
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}

		// 处理 header
		if strings.HasPrefix(line, "<!") {
			// 一个document只允许一个 header
			if nil == root {

				//解析 header
				pAttributes, err := ParseDocumentHeader(line)
				if nil != err {
					firstErr = err
					break
				}

				attrs := attrListToMap(pAttributes)
				if nil == attrs {
					firstErr = fmt.Errorf("parse etext error, document header error")
					break
				}

				root = &ENode{Id: 0, parent: nil, Name: "", Value: nil, Attribes: attrs, Children: nil, isEnd: true}
				continue
			} else {
				firstErr = fmt.Errorf("parse etext error, document has many header")
				break
			}
		}

		// d
		if nil == root {
			firstErr = fmt.Errorf("parse etext error, document header not at first")
			break
		}

		// element 开始
		if strings.HasPrefix(line, "<") {
			// element 结束
			if strings.HasPrefix(line, "</") {
				if nil == parentNode {
					firstErr = fmt.Errorf("parse etext error, document element end tag error: %s", line)
					break
				}

				if parentNode.isEnd {
					firstErr = fmt.Errorf("parse etext error, document element many end tag error: name: %s, end tag: %s", parentNode.Name, line)
					break
				}

				parentNode.isEnd = true
				parentNode = parentNode.parent

			} else {
				name, isEnd, pAttributes, err := ParseNodeLine(line)
				if nil != err {
					firstErr = err
					break
				}
				attrs := attrListToMap(pAttributes)
				node := &ENode{Id: seqGen.Generate(), parent: parentNode, Name: name, Value: nil, Attribes: attrs, isEnd: isEnd}
				if nil == parentNode {
					parentNode = root
				}
				parentNode.AddChildren(node)

				if !isEnd {
					parentNode = node
				}
			}

		} else { // content
			if nil == parentNode {
				firstErr = fmt.Errorf("parse etext error, document element end tag error, content not in element: %s", line)
				break
			}

			if parentNode.isEnd {
				firstErr = fmt.Errorf("parse etext error, document element many end tag error, content not in element, name: %s, end tag: %s", parentNode.Name, line)
				break
			}

			if parentNode.Value == nil {
				parentNode.Value = bytes.NewBufferString("")
			}
			parentNode.Value.WriteString(line)
			parentNode.Value.WriteByte('\n')
		}
	}

	if nil != firstErr {
		return nil, firstErr
	}

	return root, nil
}

func ParseRootBytes(buf *bytes.Buffer) (*ENode, error) {
	if buf.Len() < 1 {
		return nil, fmt.Errorf("parse etext error, document is empty")
	}

	var firstErr error = nil
	var root *ENode = nil
	var parentNode *ENode = nil

	for {
		firstErr = nil
		line, err := buf.ReadString('\n')
		if nil != err {
			if io.EOF != err {
				firstErr = err
			}
			break
		}

		line = strings.TrimSpace(line)
		//忽略行注释 与 空行
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}

		// 处理 header
		if strings.HasPrefix(line, "<!") {
			// 一个document只允许一个 header
			if nil == root {

				//解析 header
				pAttributes, err := ParseDocumentHeader(line)
				if nil != err {
					firstErr = err
					break
				}

				attrs := attrListToMap(pAttributes)
				if nil == attrs {
					firstErr = fmt.Errorf("parse etext error, document header error")
					break
				}

				root = &ENode{Id: 0, parent: nil, Name: "", Value: nil, Attribes: attrs, Children: nil, isEnd: true}
				continue
			} else {
				firstErr = fmt.Errorf("parse etext error, document has many header")
				break
			}
		}

		// d
		if nil == root {
			firstErr = fmt.Errorf("parse etext error, document header not at first")
			break
		}

		// element 开始
		if strings.HasPrefix(line, "<") {
			// element 结束
			if strings.HasPrefix(line, "</") {
				if nil == parentNode {
					firstErr = fmt.Errorf("parse etext error, document element end tag error: %s", line)
					break
				}

				if parentNode.isEnd {
					firstErr = fmt.Errorf("parse etext error, document element many end tag error: name: %s, end tag: %s", parentNode.Name, line)
					break
				}

				parentNode.isEnd = true
				parentNode = parentNode.parent

			} else {
				name, isEnd, pAttributes, err := ParseNodeLine(line)
				if nil != err {
					firstErr = err
					break
				}
				attrs := attrListToMap(pAttributes)
				node := &ENode{Id: seqGen.Generate(), parent: parentNode, Name: name, Value: nil, Attribes: attrs, isEnd: isEnd}
				if nil == parentNode {
					parentNode = root
				}
				parentNode.AddChildren(node)

				if !isEnd {
					parentNode = node
				}
			}

		} else { // content
			if nil == parentNode {
				firstErr = fmt.Errorf("parse etext error, document element end tag error, content not in element: %s", line)
				break
			}

			if parentNode.isEnd {
				firstErr = fmt.Errorf("parse etext error, document element many end tag error, content not in element, name: %s, end tag: %s", parentNode.Name, line)
				break
			}

			if parentNode.Value == nil {
				parentNode.Value = bytes.NewBufferString("")
			}
			parentNode.Value.WriteString(line)
			parentNode.Value.WriteByte('\n')
		}
	}

	if nil != firstErr {
		return nil, firstErr
	}

	return root, nil
}

func ParseRootString(content string) (*ENode, error) {
	buf := bytes.NewBufferString(content)
	return ParseRootBytes(buf)
}

func attrListToMap(attrList *klists.KList[*EAttribute]) map[string]string {
	attrs := make(map[string]string)
	for e := attrList.Front(); e != nil; e = e.Next() {
		attrs[e.Value.Name] = e.Value.Value
	}
	if len(attrs) == 0 {
		return nil
	}
	return attrs
}

func ParseDocumentHeader(line string) (*klists.KList[*EAttribute], error) {
	//<! Entity=华东 type=测试2011-11-03 dataTime='20120411 11:12:14' !>
	line = strings.TrimSpace(line)

	// header 格式不对
	if !strings.HasPrefix(line, "<!") {
		return nil, fmt.Errorf("error header: %s", line)
	}

	if !strings.HasSuffix(line, "!>") {
		endPos := strings.LastIndex(line, "!>")
		if endPos < 0 {
			return nil, fmt.Errorf("error header: %s", line)
		}

		// 去掉行尾注释, 如果注释结尾为`!>`, 会导致解析异常
		commentPos := strings.LastIndex(line, "//")
		if commentPos > endPos {
			fmt.Printf("delete comments %s\n", line[commentPos:])
			line = line[:commentPos]
			line = strings.TrimSpace(line)
		}
	}

	// 去掉 <!!>
	line = line[2 : len(line)-2]
	line = strings.TrimSpace(line)

	attributes, err := parseLineKeyWithVal(line)
	if nil != err {
		return nil, err
	}
	return attributes, nil
}

// @bref 解析一行 xml element
//
// @param `line` `string`  一行单个 xml element, 例如 `<a f=1 h=2 />`, 假如是 `<a f=1 h=2 /></a>` 会解析失败
//
// @return
//
//			string: 解析出来的name字段
//			bool:   是否为自闭合语句; 例如: <a f=1 h=2 />
//		 *klists.KList[*EAttribute]: 解析出来的属性列表
//	  error: 错误信息
func ParseNodeLine(line string) (string, bool, *klists.KList[*EAttribute], error) {
	//<test::华东 DDMM='华东电网' date='2012-04-11 11:12' >
	line = strings.TrimSpace(line)

	// header 格式不对
	if !strings.HasPrefix(line, "<") {
		return "", false, nil, fmt.Errorf("parse etext error, error header: %s", line)
	}

	if !strings.HasSuffix(line, ">") {
		endPos := strings.LastIndex(line, ">")
		if endPos < 0 {
			return "", false, nil, fmt.Errorf("parse etext error, error header: %s", line)
		}

		// 去掉行尾注释, 如果注释结尾有`>`字符, 可能会导致解析异常
		commentPos := strings.LastIndex(line, "//")
		if commentPos > endPos {
			//fmt.Printf("delete comments %s\n", line[commentPos:])
			line = line[:commentPos]
			line = strings.TrimSpace(line)
		}
	}

	// 去掉 <>
	line = line[1 : len(line)-1]
	line = strings.TrimSpace(line)
	isEnd := false
	// 如果 xml 单行闭合: <a f0=0 f1=1/>
	if strings.HasSuffix(line, "/") {
		isEnd = true
		line = line[:len(line)-1]
		line = strings.TrimSpace(line)
	}

	pos := strings.Index(line, " ")
	name := ""
	if pos > -1 {
		name = line[:pos]
		line = line[pos+1:]
		line = strings.TrimSpace(line)
	} else {
		name = line
	}
	attributes, err := parseLineKeyWithVal(line)
	if nil != err {
		return "", false, nil, err
	}

	return name, isEnd, attributes, nil
}

func parseLineKeyWithVal(line string) (*klists.KList[*EAttribute], error) {
	//Entity=华东 type=测试2011-11-03 dataTime='20120411 11:12:14'
	length := len(line)
	if length < 1 {
		return nil, fmt.Errorf("parse etext error, error Key-Value line: %s", line)
	}

	attributes := klists.New[*EAttribute]()
	tmp := line
	endPos := length

	startPos := strings.LastIndex(tmp, "=")
	for startPos > -1 {
		value := tmp[startPos+1 : endPos]
		value = strings.TrimFunc(value, func(r rune) bool {
			return r == '\'' || r == '"' || r == ' ' || r == '\t'
		})
		tmp = tmp[0:startPos]
		endPos = startPos

		startPos = strings.LastIndex(tmp, " ")
		if startPos == -1 {
			startPos = -1
		}

		key := tmp[startPos+1 : endPos]
		// fmt.Printf("%s, key: %s, value: %s, startPos: %d, endPos: %d\n", tmp, key, value, startPos, endPos)
		attributes.PushBack(&EAttribute{Name: strings.TrimSpace(key), Value: strings.TrimSpace(value)})
		if startPos == -1 {
			break
		}

		tmp = tmp[:startPos]
		endPos = startPos
		startPos = strings.LastIndex(tmp, "=")
	}
	return attributes, nil
}

// @bref 解析一行 e文本数据
func ParseEText(line string, sep string) (*klists.KList[string], error) {
	//#1 '花花电网' 花花电网 '2011-11-03 00:00:02.0' 32  ''
	line = strings.TrimSpace(line)

	length := len(line)
	if length < 1 {
		return nil, fmt.Errorf("parse etext error, error Key-Value line: %s", line)
	}

	delim := sep
	fields := klists.New[string]()
	tmp := line

	if strings.HasSuffix(tmp, "\"") {
		delim = "\""
		tmp = tmp[:length-1]
	} else if strings.HasSuffix(tmp, "'") {
		delim = "'"
		tmp = tmp[:length-1]
	}

	endPos := len(tmp)

	startPos := strings.LastIndex(tmp, delim)
	for startPos > -1 {
		value := tmp[startPos+1 : endPos]
		value = strings.TrimFunc(value, func(r rune) bool {
			return r == '\'' || r == '"' || r == ' ' || r == '\t'
		})

		tmp = tmp[0:startPos]
		tmp = strings.TrimSpace(tmp)
		length = len(tmp)

		if strings.HasSuffix(tmp, "\"") {
			delim = "\""
			tmp = tmp[:length-1]
		} else if strings.HasSuffix(tmp, "'") {
			delim = "'"
			tmp = tmp[:length-1]
		} else {
			delim = sep
		}
		endPos = len(tmp)
		startPos = strings.LastIndex(tmp, delim)

		fields.PushFront(value)
	}
	fields.PushFront(tmp)
	return fields, nil
}

func ParseETable(root *ENode, path string) (*klists.KList[*klists.KList[string]], error) {
	node, err := GetENodeByPath(root, path)
	if nil != err {
		return nil, fmt.Errorf("node not found, path: %s", path)
	}

	if node.Value == nil {
		return nil, fmt.Errorf("node %s is empty, path: %s", node.Name, path)
	}

	var firstErr error = nil
	var records *klists.KList[*klists.KList[string]] = nil
	header := ""
	delim := " "

	buf := node.Value

	// 获取header
	for {
		line, err := buf.ReadString('\n')
		if nil != err {
			if io.EOF != err {
				firstErr = err
			}
			break
		}

		line = strings.TrimSpace(line)
		//忽略行注释 与 空行
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}

		// 处理 header
		if strings.HasPrefix(line, "@") {
			header = line
			break
		}
	}

	if nil != firstErr {
		return nil, firstErr
	}

	if strings.LastIndex(header, "\t") != -1 {
		delim = "\t"
	}

	if strings.HasPrefix(header, "@#") { // 多列式
		records, firstErr = parseMultColTable(buf, delim)
	} else if strings.HasPrefix(header, "@@") { // 单列式
		records, firstErr = parseSigleColTable(buf, delim)
	} else { //横表式
		records, firstErr = parseTable(buf, delim)

		items, err := ParseEText(header, delim)
		if err != nil {
			firstErr = fmt.Errorf("etable parse error, header: %s", header)
		} else {
			row := klists.New[string]()
			for e := items.Front(); e != nil; e = e.Next() {
				row.PushBack(e.Value)
			}
			records.PushFront(row)
		}
	}

	if nil != firstErr {
		return nil, firstErr
	}

	return records, nil
}

// @bref 单列式表格解析
//
// 例如:
//
// @@顺序 属性名 属性值
//
// #1 单位名称 花花电网
//
// #2 发生时间 '2011-11-03 00:00:02.0'
//
// #3 次数 32
//
// -------------------------------------
//
// #1 单位名称 花花电网
//
// #2 发生时间 '2011-11-03 00:00:02.0'
//
// #3 次数 32
func parseSigleColTable(buf *bytes.Buffer, delim string) (*klists.KList[*klists.KList[string]], error) {
	var firstErr error = nil
	var preRecord *klists.KList[string] = nil
	record := make(map[string]string)
	records := klists.New[*klists.KList[string]]()
	for {
		firstErr = nil
		line, err := buf.ReadString('\n')
		if nil != err {
			if io.EOF != err {
				firstErr = err
			}
			break
		}

		line = strings.TrimSpace(line)
		//忽略行注释 与 空行
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}

		// 多余的 header
		if strings.HasPrefix(line, "@") {
			firstErr = fmt.Errorf("etable has many header")
			break
		}

		items, err := ParseEText(line, delim)
		if err != nil {
			firstErr = fmt.Errorf("etable parse error, line: %s", line)
			break
		}

		if items.Len() != 3 {
			firstErr = fmt.Errorf("etable parse error, sigle colum table fields count != 3 line: %s", line)
			break
		}
		key := *items.At(1)
		val := *items.At(2)

		if kmaps.HasKey[string, string](record, key) { //凑足了一条记录
			header := make([]string, 0, len(record))
			for k := range record {
				header = append(header, k)
			}
			sort.Strings(header)
			row := klists.New[string]()

			for _, k := range header {
				row.PushBack(record[k])
			}
			if nil != preRecord && preRecord.Len() != row.Len() {
				firstErr = fmt.Errorf("etable parse error, fields error, pre at line: %s", line)
				break
			}
			records.PushBack(row)
			preRecord = row
			kmaps.Clear[string, string](record)
		}

		record[key] = val
	}
	if nil != firstErr {
		return nil, firstErr
	}

	header := make([]string, 0, len(record))
	row := klists.New[string]()
	for k := range record {
		header = append(header, k)
	}
	sort.Strings(header)
	for _, k := range header {
		row.PushBack(k)
	}
	records.PushFront(row)

	return records, nil
}

// @bref 多列式表格解析
//
// 例如:
//
// @#顺序 属性名 '花花电网' '花花电网' '花花电网' '花花电网'
//
// #1 单位名称 花花电网 花花电网 花花电网 花花电网
//
// #2 发生时间 '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0' '2011-11-03 00:00:02.0'
//
// #3 次数 32 32 32 32
func parseMultColTable(buf *bytes.Buffer, delim string) (*klists.KList[*klists.KList[string]], error) {
	var firstErr error = nil

	records := klists.New[*klists.KList[string]]()
	for {
		firstErr = nil
		line, err := buf.ReadString('\n')
		if nil != err {
			if io.EOF != err {
				firstErr = err
			}
			break
		}

		line = strings.TrimSpace(line)
		//忽略行注释 与 空行
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}

		// 多余的 header
		if strings.HasPrefix(line, "@") {
			firstErr = fmt.Errorf("etable has many header")
			break
		}

		items, err := ParseEText(line, delim)
		if err != nil {
			firstErr = fmt.Errorf("etable parse error, line: %s", line)
			break
		}

		if records.Len() == 0 {
			idx := 0
			for e := items.Front(); e != nil; e = e.Next() {
				if idx == 0 {
					idx++
					continue
				}
				row := klists.New[string]()
				row.PushBack(e.Value)
				records.PushBack(row)
				idx++
			}
		} else {
			idx := 0
			for e := items.Front(); e != nil; e = e.Next() {
				if idx == 0 {
					idx++
					continue
				}

				rowNum := 0
				for er := records.Front(); er != nil; er = er.Next() {
					if idx-1 == rowNum {
						er.Value.PushBack(e.Value)
						break
					}
					rowNum++
				}
				idx++
			}
		}
	}

	if nil != firstErr {
		return nil, firstErr
	}

	return records, nil
}

// @bref 横表式表格解析
//
// 例如:
//
// @顺序 单位名称 发生时间 次数
//
// #1 花花电网 无 1000
//
// #2 花花电网 无 1000
//
// #3 花花电网 无 1000
//
// #4 花花电网 无 1000
//
// #5 花花电网 无 1000
//
// #6 花花电网 无 1000
func parseTable(buf *bytes.Buffer, delim string) (*klists.KList[*klists.KList[string]], error) {
	var firstErr error = nil
	records := klists.New[*klists.KList[string]]()
	for {
		firstErr = nil
		line, err := buf.ReadString('\n')
		if nil != err {
			if io.EOF != err {
				firstErr = err
			}
			break
		}

		line = strings.TrimSpace(line)
		//忽略行注释 与 空行
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}

		// 多余的 header
		if strings.HasPrefix(line, "@") {
			firstErr = fmt.Errorf("etable has many header")
			break
		}

		items, err := ParseEText(line, delim)
		if err != nil {
			firstErr = fmt.Errorf("etable parse error, line: %s", line)
			break
		}

		row := klists.New[string]()
		for e := items.Front(); e != nil; e = e.Next() {
			row.PushBack(e.Value)
		}
		records.PushBack(row)

	}

	if nil != firstErr {
		return nil, firstErr
	}

	return records, nil
}

////////////////////////////////////////////////////////////////////
