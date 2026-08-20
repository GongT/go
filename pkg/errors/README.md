# 超级无敌牛逼错误包

此程序库提供一系列错误处理的工具，扩展标准库的错误处理能力

特性:
* 所有错误对象都带有栈信息
* 与标准库非常相似的Wrap链，使用`fmt.Errorf("%w", err)`也不影响输出栈信息
* 支持与[gitlab.com/tozd/go/errors](https://gitlab.com/tozd/go/errors)相同的附加数据协议
* 提供一个紧凑的CLI错误输出方法
* 通过`Details`实现与JavaScript中的`Error.cause`类似的功能    
    > 注意: 此cause和go社区标准中的Cause()完全无关！
	>
	> 为避免歧义，本库中不使用、不兼容Cause()方法，仅使用新标准库`Unwrap()`一种方法实现错误链
	>
	> 使用“reason”作为Details中的键

## 创建错误对象

#### 常用方法

| 方法 | 对应操作 | 说明 |
|---|----|-----|
| `NewAnonymous(...)` | `errors.New(fmt.Sprintf(...))` | 创建一个匿名错误对象 |
| `var tmpl = NewTemplate("failed to do %s")` <br> `tmpl.Instance("something")` | `var FsErr = errors.New("failed");` <br> `fmt.Errorf(FsErr, "something")` | 创建一个模板错误对象（并实例化） |
| `Extend(err, msg, ...)` <br> `ExtendTrace(...)` | `fmt.Errorf("msg: %w", err, ...)` | 给错误添加上下文（Wrap核心机制） |
| `Join(errs...)` <br> `Concat(errs, msg, ...)` | `errors.Join(errs...)` | 将多个错误对象合并为一个 |

#### 其他

| 方法 |  说明 |
|---|----|
| `Ensure(err)` | 确保err是本库中的类型，如果不是则创建一个 |
| `EnsureTrace(err)` | 和`Ensure(err)`类似，但只在err没有`StackTrace`的时候才创建 |
| `GetStackTrace(err)` | 获取错误对象的栈信息，如果没有则返回nil |
| `Unjoin(err)` | 如果err是Join类型，则返回所有子错误，否则返回nil |
| `IterateErrorStack(err)` | 迭代错误对象的栈信息，返回`runtime.Frame`类型的迭代器 |

#### 附加数据

| 方法 |  说明 |
|---|----|
| `SetDetail(err, key, value)` | 设置一个key-value |
| `SetDetails(err, map[string]any)` | maps.Copy |
| `WithDetails(err, key1, value1, key2, value2...)` | 设置多个key-value，和`tozd/errors`相同API |
| `GetDetails(err)` | 获取所有附加数据，返回map |
| `GetReason(err)` | 获取附加数据中的`reason`字段，返回error类型 |
| `GetCode(err)` | 获取附加数据中的`code`字段，返回int类型 |

其中 `SetDetail(s)`、`WithDetails` 也作为`Err`类型的方法使用。

#### 递归

- 大部分**用于读取的**函数式API都是递归的，例如 GetStackTrace、GetDetails 等
- 对象上的方法则全都是仅读写自身
- Unjoin是例外，它仅仅简单对err对象调用Unwrap
