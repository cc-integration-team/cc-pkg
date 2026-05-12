# logger

Package `logger` cung cấp interface logging thống nhất với zerolog và slog adapter, hỗ trợ structured fields và tính năng masking thông tin nhạy cảm trước khi ghi ra file/console.

## Sử dụng cơ bản

```go
import "github.com/cc-integration-team/cc-pkg/v3/pkg/logger"

// Khởi tạo và set default logger
logger.SetDefaultLogger(logger.NewZerologAdapter(cfg.Logger))

// Log đơn giản
logger.Info("server started")
logger.Errorf("connect failed: %v", err)

// Structured logging với WithFields (khuyến nghị)
logger.WithFields(logger.Fields{
    "customerPhone": phone,
    "callerID":      callerID,
    "agentID":       agentID,
}).Info("call event received")
```

## Cấu hình

```yaml
logger:
  service: "my-service"   # tên service, gắn vào mọi log line
  caller: false           # bật/tắt hiển thị file:line

  masking:
    enabled: true
    fields:               # top-level JSON fields cần mask
      - customerPhone
      - callerID
      - calleeID
    nestedFields:         # fields nằm trong nested object
      metadata:
        - callVariable8
        - callVariable3

  file:
    enabled: true
    path: "./log/app.log"
    level: "debug"
    maxSize: 100          # MB
    maxBackups: 7
    maxAge: 30            # ngày
    compress: true
    pretty: false

  console:
    enabled: true
    level: "error"
    pretty: true
```

## Masking thông tin nhạy cảm

### Mục đích

Che giấu số điện thoại và các thông tin nhạy cảm trong log output mà không ảnh hưởng đến dữ liệu thực tế được xử lý trong app.

### Format mask

Luôn dùng 6 dấu `*` cố định + 3 ký tự cuối, bất kể độ dài gốc:

```
"0901234567"  →  "******567"
"08123456789" →  "******789"
"1234"        →  "******234"
"123"         →  "123"         (≤3 ký tự: giữ nguyên)
""            →  ""            (empty: giữ nguyên)
```

Dùng prefix cố định để không lộ độ dài số điện thoại gốc qua log.

### Luồng hoạt động

```
App gọi logger.WithFields(...).Info("message")
          │
          ▼
zerolog.Event serialize → JSON bytes
{"level":"info","customerPhone":"0901234567","metadata":{"callVariable8":"0909090909"},...}
          │
          ▼
MaskingLevelWriter.WriteLevel()        ← intercept tại đây
  1. level < Info?  → bypass, ghi thẳng (debug không mask)
  2. json.Unmarshal → map[string]json.RawMessage
  3. Mask top-level fields (customerPhone, callerID, ...)
  4. Với mỗi nestedFields parent (metadata):
       - Unmarshal nested object
       - Mask các child keys được cấu hình
       - Marshal lại nested object
  5. json.Marshal → bytes mới
  6. Ghi vào writer gốc (file / console)
          │
          ▼
Output: {"level":"info","customerPhone":"******567","metadata":{"callVariable8":"******909"},...}
```

### Cấu hình fields

**`fields`** — top-level JSON keys, tương ứng với tên field trong `WithFields`:

```go
logger.WithFields(logger.Fields{
    "customerPhone": "0901234567",   // → "******567"
    "callerID":      "0912345678",   // → "******678"
}).Info("call")
```

**`nestedFields`** — keys nằm trong một nested map. Chỉ mask đúng các key được khai báo, các key khác trong cùng object giữ nguyên:

```go
logger.WithFields(logger.Fields{
    "metadata": map[string]any{
        "callVariable8": "0909090909",  // → "******909"  (trong config)
        "callVariable3": "other",       // giữ nguyên     (không trong config)
        "queueName":     "support",     // giữ nguyên     (không trong config)
    },
}).Info("call metadata")
```

### Debug level không bị mask

Masking chỉ áp dụng từ `Info` trở lên. Log `Debug` ghi nguyên giá trị để thuận tiện troubleshoot.

```go
logger.WithFields(logger.Fields{"customerPhone": "0901234567"}).Debug("raw value") // không mask
logger.WithFields(logger.Fields{"customerPhone": "0901234567"}).Info("masked")     // → ******567
```

### Bật/tắt

Khi `enabled: false` (hoặc không khai báo `masking:` trong config), `MaskingLevelWriter` không được tạo — writer chain giữ nguyên như trước, zero overhead.

### Giới hạn

Masking chỉ hoạt động khi dùng `WithFields` để log structured data. Các log dùng format string như `Infof("phone=%s", phone)` sẽ không được mask vì giá trị nằm trong chuỗi `message`, không phải JSON field riêng biệt.

```go
// Không mask được
logger.Infof("call from %s", phone)

// Mask được
logger.WithFields(logger.Fields{"customerPhone": phone}).Info("call")
```

## Interface

```go
type Logger interface {
    Debug(msg string)
    Debugf(msg string, args ...any)
    Info(msg string)
    Infof(msg string, args ...any)
    Warn(msg string)
    Warnf(msg string, args ...any)
    Error(msg string)
    Errorf(msg string, args ...any)
    Fatal(msg string)
    Fatalf(msg string, args ...any)
    WithFields(fields Fields) Logger
}

type Fields map[string]any
```

## Context

```go
// Inject logger vào context
ctx = logger.WithContext(ctx, myLogger)

// Lấy logger từ context (trả về defaultLogger nếu không có)
log := logger.FromContext(ctx)
log.Info("handler started")
```
