# 序列化、反序列化生成器

[示例与测试](./assets/serialize)
[包体](pkg/serialize)
[生成器main](cmd/serialize/main.go)
[二进制数据操作工具](pkg/strings/packer)

## 功能
- （通过packages和ast）静态分析指定结构体，生成2个函数 Marshal、Unmarshal，对每个属性依次进行序列化、反序列化操作
	- 此函数将写入定义的文件旁边，所以可以直接访问私有属性
	- 无需记录类型等信息，因为序列化和反序列化是基于结构体的字段顺序进行的
- 遇到嵌套结构体、interface、error时，调用serialize.Marshal和serialize.Unmarshal实现序列化和反序列化
	- 如过它们返回nil（未定义支持的接口），则返回运行时错误
	- 不尝试递归解析类型
- 支持map、slice
- 遇到chan、func时，调用的
