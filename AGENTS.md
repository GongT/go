# Project Instruction

This repository is a Go library project, with a few CLI tools.

## Scope

- Applies to all Go source files in this repository.
- Prefer consistency with existing repository architecture, naming, and package boundaries.

## Language And Version Rules

- Treat Go as the default implementation language unless the task explicitly requires another language.
- Prefer modern (newest) Go syntax and standard-library capabilities.
- When multiple approaches are valid, choose the clearer and more idiomatic modern Go approach.
- Use `require` and `assert` in `github.com/stretchr/testify` for testing, avoid directly calling `t.Xxx`.
	- Add `myenv.RedirectDebugTesting(t)` statement at the beginning of each Test function.
	- Use same package name for test files as the package being tested, not using `<package>_test`.

## Comment And Error Message Localization

- Write code comments in Simplified Chinese.
- Write developer-facing and user-facing and testing error messages in Simplified Chinese.
- Keep logging and error text concise, actionable, and context-rich.
- Do not translate or modify third-party error strings when wrapping; preserve original errors with `%w` and add Chinese context around them. Never use `%s` or `%v` for errors.

## Misc

- Always return errors when they occur, either wrap them with context or return them as-is, or completely ignore it as intended; never silently swallow errors
- Keep protocol fields, external API keys, environment variable names, and machine-readable constants unchanged.
- Use `pkg/errors` in this repository for error stack trace and context extension
	- (rarely) Except when using it would cause cyclic dependencies.
- Use `pkg/fsys/fpath` for path manipulation instead of `path/filepath`. but still use `os` for real file operations. Strictly prohibited any string manipulation of paths.

## Example Patterns

```go
// 正确: 注释使用中文，错误上下文使用中文，并保留原始错误
if err := doWork(); err != nil {
    return errors.Extend(err, "执行任务失败")
}

// 正确: 错误信息清晰且可定位
return errors.Extend(err, "读取配置文件失败，路径=%s", path)

// 避免: 英文注释和英文错误信息
// load config from file
return errors.Extend(err, "failed to load config")
```

```go
// 避免: err变量被忽略
if err := doWork(); err != nil { return errors.New("doWork失败") }

// 正确: 返回带有中文上下文的错误
if err := doWork(); err != nil {
    return errors.Extend(err, "执行任务失败")
}
```

```go
// 正确: 现代化语法
for i := range 100 {}
func example(value any) {}

// 避免: 传统语法
for i := 0; i < 100; i++ {}
func example(value interface{}) {}
```
