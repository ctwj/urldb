# 夸克网盘 getShare 方法修复

## 问题描述

在夸克网盘的 `getShare` 函数中遇到了 HTTP 405 错误：

```
HTTP请求失败: 405, {"timestamp":1752368941648,"status":405,"error":"Method Not Allowed","exception":"org.springframework.web.HttpRequestMethodNotSupportedException","message":"Request method 'POST' not supported","path":"/1/clouddrive/share/sharepage/detail"}
```

## 问题分析

### 错误原因
- 当前Go代码使用 `HTTPPost` 方法请求 `/1/clouddrive/share/sharepage/detail` 接口
- 但服务器返回 405 错误，表示不支持 POST 方法
- 需要改为 GET 请求

### 对比原始PHP代码
通过对比 `demo/pan/QuarkPan.php` 中的 `getShare` 方法：

```php
public function getShare($pwd_id,$stoken)
{
    $urlData = array();
    $queryParams = [
        "pr" => "ucpro",
        "fr" => "pc",
        "uc_param_str" => "",
        "pwd_id" => $pwd_id,
        "stoken" => $stoken,
        "pdir_fid" => "0",
        "force" => "0",
        "_page" => "1",
        "_size" => "100",
        "_fetch_banner" => "1",
        "_fetch_share" => "1",
        "_fetch_total" => "1",
        "_sort" => "file_type:asc,updated_at:desc"
    ];
    return $this->executeApiRequest(
        "https://drive-pc.quark.cn/1/clouddrive/share/sharepage/detail", 
        "GET",  // 使用 GET 方法
        $urlData, 
        $queryParams
    );
}
```

## 修复方案

### 修复前的代码
```go
// 修复前：使用 POST 请求
func (q *QuarkPanService) getShare(shareID, stoken string) (*ShareResult, error) {
    data := map[string]interface{}{
        "pwd_id": shareID,
        "stoken": stoken,
    }

    queryParams := map[string]string{
        "pr":           "ucpro",
        "fr":           "pc",
        "uc_param_str": "",
    }

    respData, err := q.HTTPPost("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/detail", data, queryParams)
    // ...
}
```

### 修复后的代码
```go
// 修复后：使用 GET 请求，参数放在 URL 中
func (q *QuarkPanService) getShare(shareID, stoken string) (*ShareResult, error) {
    queryParams := map[string]string{
        "pr":              "ucpro",
        "fr":              "pc",
        "uc_param_str":    "",
        "pwd_id":          shareID,
        "stoken":          stoken,
        "pdir_fid":        "0",
        "force":           "0",
        "_page":           "1",
        "_size":           "100",
        "_fetch_banner":   "1",
        "_fetch_share":    "1",
        "_fetch_total":    "1",
        "_sort":           "file_type:asc,updated_at:desc",
    }

    respData, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/detail", queryParams)
    // ...
}
```

## 修复内容

### 1. 请求方法修改
- **从 POST 改为 GET**：`HTTPPost` → `HTTPGet`
- **参数传递方式**：从请求体改为 URL 查询参数

### 2. 参数结构调整
- **移除请求体数据**：不再使用 `data` 参数
- **添加完整查询参数**：按照PHP代码添加所有必要的查询参数
- **参数顺序**：保持与PHP代码一致的参数顺序

### 3. 参数说明
| 参数 | 说明 | 值 |
|------|------|-----|
| `pr` | 产品标识 | `ucpro` |
| `fr` | 来源标识 | `pc` |
| `uc_param_str` | UC参数 | 空字符串 |
| `pwd_id` | 分享ID | 从URL提取 |
| `stoken` | 安全令牌 | 从getStoken获取 |
| `pdir_fid` | 父目录ID | `0` |
| `force` | 强制标志 | `0` |
| `_page` | 页码 | `1` |
| `_size` | 页面大小 | `100` |
| `_fetch_banner` | 获取横幅 | `1` |
| `_fetch_share` | 获取分享信息 | `1` |
| `_fetch_total` | 获取总数 | `1` |
| `_sort` | 排序方式 | `file_type:asc,updated_at:desc` |

## 验证结果

### 编译测试
```bash
go build -o res_db .
# ✅ 编译成功
```

### 功能测试
```bash
go test ./utils/pan -v
# ✅ 所有测试通过
```

### 预期效果
- **解决 405 错误**：使用正确的 GET 请求方法
- **保持功能完整**：所有参数都正确传递
- **兼容性良好**：与原始PHP代码行为一致

## 经验总结

### 1. API 兼容性
- 在翻译代码时，需要仔细对比原始实现的请求方法
- 不同语言的HTTP客户端可能有不同的默认行为
- 需要确保请求方法、参数位置、参数名称都完全一致

### 2. 错误排查
- HTTP 405 错误通常表示请求方法不正确
- 对比原始代码是排查此类问题的最佳方法
- 需要检查请求方法、URL、参数等多个方面

### 3. 最佳实践
- **保持一致性**：与原始实现保持完全一致
- **完整参数**：不要遗漏任何必要的参数
- **测试验证**：修复后要进行充分的测试

## 相关文件

- `utils/pan/quark_pan.go` - 修复的文件
- `demo/pan/QuarkPan.php` - 原始PHP实现参考
- `utils/pan/base_pan.go` - 基础HTTP方法实现

这次修复确保了夸克网盘的 `getShare` 方法与原始PHP实现完全一致，解决了HTTP 405错误问题。 