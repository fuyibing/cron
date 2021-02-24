# Crontab

> Package `github.com/fuyibing/cron/v2`

### Ticker表达式.

1. `10s` 每隔`10秒`执行一次
1. `5m` 每隔`5分钟`执行一次
1. `3h` 每隔`3小时`执行一次
1. `1d` 每隔`1天`执行一次
1. `00:00` 每天的`00:00:00`时执行一次
1. `00:15:30` 每天的`00:15:30`时执行一次
1. `00:15:30, 01:15` 每天的`00:15:30`和`01:15:00`时各执行一次
1. **Crontab格式**
    1. 第`1`列 : Second, 秒.
    1. 第`2`列 : Minute, 分钟.
    1. 第`3`列 : Hour, 时.
    1. 第`4`列 : Day, 天.
    1. 第`5`列 : Month, 月.
    1. 第`6`列 : Week, 周.
    1. 第`7`列 : Year, 年.

### 单节点模式

```text
tick := cron.NewTicker("test1", "10m", callback)
tick.SingleNode(true)
```
