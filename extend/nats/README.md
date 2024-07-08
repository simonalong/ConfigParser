# nats

对业内的emqx客户端进行配置化封装，用于简化获取
### 全部配置
```yaml
gole:
  nats:
    # 是否开启emqx，默认关闭
    enable: true
    servers:
      # 域名格式1
      - "tcp://{user}:{password}@{host}:{port}"
      # 域名格式2
      - "tcp://{host}:{port}"
```
提供封装的 `nats客户端api`
```go
func NewEmqxClient() (mqtt.Client, error) {}
```
#### 示例：
```yaml

```yaml
gole:
  nats:
    enable: true
    servers:
      - "tcp://{host}:{port}"
```
```go
import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/simonalong/gole/extend/emqx"
 )

```
