# Naglfar Server API

server_address: http://35.86.17.44:2048

作为 Demo，http status 200 为成功，其他情况一律认为失败，失败情况下 body 内容为 debug 信息。

## 获取日志片段列表

> 用户每次上传的一段日志叫日志片段

```
curl ${server_address}/api/v1/fragments
```

返回示例:

```json
["test-cluster"]
```

## 查询含有某个 event 的日志片段列表

```
curl ${server_address}/api/v1/${start}/${stop}/events/${eid}/fragments
```

参数说明:

- start: 搜索的开始时间
- end: 搜索的结束时间
- eid: event id，一个 event 代表一类日志

返回示例:

```json
["test-cluster"]
```

## 查询某个日志片段中的日志的 event 分布

```
curl ${server_address}/api/v1/{start}/{stop}/fragments/{fid}/logs/stats
```

参数说明:

- start: 搜索的开始时间
- end: 搜索的结束时间
- fid: fragment id，唯一代表一份日志

返回示例:

```json
[{
	"count": 666058,
	"event_id": "10107",
	"message": "[BIG_TXN]",
	"weight": 0.3990594877359161
}, {
	"count": 686393,
	"event_id": "10214",
	"message": "[TIME_COP_PROCESS]",
	"weight": 0.35494386275500994
}, {
	"count": 198540,
	"event_id": "10180",
	"message": "use snapshot schema",
	"weight": 0.19217699635458843
}, {
	"count": 57395,
	"event_id": "10099",
	"message": "kill",
	"weight": 0.05555554903682684
}, {
	"count": 96102,
	"event_id": "10087",
	"message": "wait response is cancelled",
	"weight": 0.04616782394648621
}, {
	"count": 18317,
	"event_id": "10023",
	"message": "unexpected `Flen` value(-1) in CONCAT's args",
	"weight": 0.024485539259457182
}
...
]
```

## 查询某个 event 的日志详情

```
curl ${server_address}/api/v1/{start}/{stop}/fragments/{fid}/events/{eid}/logs
```

参数说明:

- start: 搜索的开始时间
- end: 搜索的结束时间
- fid: fragment id，唯一代表一份日志
- eid: event id，一个 event 代表一类日志

返回示例:

```json
{
	"logs": [{
		"event_id": "10035",
		"f_job": "ID:112773, Type:create schema, State:none, SchemaState:queueing, SchemaID:112772, TableID:0, RowCount:0, ArgLen:0, start time: 2021-12-17 15:20:32.562 +0800 CST, Err:\u003cnil\u003e, ErrCount:0, SnapshotVersion:0",
		"f_worker": "worker 3, tp general",
		"level": "INFO",
		"message": "[ddl] run DDL job",
		"timestamp": 1639725632
	}, {
		"event_id": "10035",
		"f_job": "ID:112789, Type:create table, State:none, SchemaState:queueing, SchemaID:112772, TableID:112788, RowCount:0, ArgLen:0, start time: 2021-12-17 15:32:29.711 +0800 CST, Err:\u003cnil\u003e, ErrCount:0, SnapshotVersion:0",
		"f_worker": "worker 3, tp general",
		"level": "INFO",
		"message": "[ddl] run DDL job",
		"timestamp": 1639726349
	}, 
    ...
    ],
    "stats": {
		"f_job": 125,
		"f_worker": 2
	}
```

## 查询某个事件的字段值域分布

```
curl ${server_address}/api/v1/{start}/{stop}/fragments/{fid}/events/{eid}/fields/{field}/logs/stats
```

- start: 搜索的开始时间
- end: 搜索的结束时间
- fid: fragment id，唯一代表一份日志
- eid: event id，一个 event 代表一类日志
- field: 字段名

返回示例:

```json
{
    "worker 3, tp general":123,
    "worker 4, tp add index":2
}
```

## 查询某个日志片段的日志趋势

```
curl ${server_address}/api/v1/${start}/${stop}/fragments/${fid}/logs/trend?events={eid1,eid2...}
```

- start: 搜索的开始时间
- end: 搜索的结束时间
- fid: fragment id，唯一代表一份日志
- eidN: event id，不传则返回所有 event 的趋势

返回示例:

```json
[{
	"event_id": "10035",
	"name": "[ddl] run DDL job",
	"points": [{
		"timestamp": 1639725900,
		"value": 1
	}, {
		"timestamp": 1639726500,
		"value": 1
	}, {
		"timestamp": 1639731300,

	}, {
		"timestamp": 1640088300,
		"value": 8
	}, ...]
}]
```

## 查询某个事件按 field 聚合的趋势

```
curl ${server_address}/api/v1/{start}/{stop}/fragments/{fid}/events/{eid}/fields/{field}/logs/trend
```

- start: 搜索的开始时间
- end: 搜索的结束时间
- fid: fragment id，唯一代表一份日志
- eid: event id，一个 event 代表一类日志
- field: 字段名


返回示例:

```json
[{
	"event_id": "",
	"name": "github_events@172.16.4.36",
	"points": [{
		"timestamp": 1639969500,
		"value": 1
	}, {
		"timestamp": 1639971900,
		"value": 4
	}, {
		"timestamp": 1640068200,
		"value": 1
	}, {
		"timestamp": 1640069100,
		"value": 1
	}]
}, {
	"event_id": "",
	"name": "perfdata@172.16.4.36",
	"points": [{
		"timestamp": 1639821000,
		"value": 2
	}, {
		"timestamp": 1639997400,
		"value": 1
	}, {
		"timestamp": 1639997700,
		"value": 2
	}]
}, {
	"event_id": "",
	"name": "pingcap_user@172.16.4.36",
	"points": [{
		"timestamp": 1639651800,
		"value": 3
	}, {
		"timestamp": 1639652100,
		"value": 2
	}, {
		"timestamp": 1639996500,
		"value": 2
	}, {
		"timestamp": 1640052900,
		"value": 2
	}, {
		"timestamp": 1640059500,
		"value": 2
	}]
}, {
	"event_id": "",
	"name": "root@172.16.4.36",
	"points": [{
		"timestamp": 1639725900,
		"value": 4
	}, {
		"timestamp": 1639726800,
		"value": 3
	}, {
		"timestamp": 1640063100,
		"value": 1
	}]
}, {
	"event_id": "",
	"name": "sync_ci_data@172.16.4.36",
	"points": [{
		"timestamp": 1639671900,
		"value": 4
	}, {
		"timestamp": 1639674300,
		"value": 4
	}, {
		"timestamp": 1639722300,
		"value": 4
	}, {
		"timestamp": 1639729500,
		"value": 4
	}, {
		"timestamp": 1639736700,
		"value": 4
	}, {
		"timestamp": 1639794300,
		"value": 4
	}, {
		"timestamp": 1639799100,
		"value": 4
	}, {
		"timestamp": 1639806300,
		"value": 4
	}, {
		"timestamp": 1639807500,
		"value": 4
	}, {
		"timestamp": 1639820700,
		"value": 4
	}, {
		"timestamp": 1639825500,
		"value": 4
	}, {
		"timestamp": 1639832700,
		"value": 4
	}, {
		"timestamp": 1639839900,
		"value": 4
	}, {
		"timestamp": 1639841100,
		"value": 4
	}, {
		"timestamp": 1639866300,
		"value": 4
	}]
}, {
	"event_id": "",
	"name": "wh-write@172.16.4.36",
	"points": [{
		"timestamp": 1639995600,
		"value": 1
	}]
}]
```

## 获取一个日志片段和其他片段的相似度


```
curl ${server_address}/api/v1/fragments/{fid}/logs/similarity
```

- fid: fragment id，唯一代表一份日志


返回示例:

```json
{"tidb-3":0.3698447346687317,"tidb-7":0.8179337382316589}
```

## 获取特定日志的正常上下界

```
curl ${server_address}/api/v1/fragments/{fid}/events/{eid}/logs/threshold
```

- fid: fragment id，唯一代表一份日志
- eid: event id，一个 event 代表一类日志

返回示例:

```json
{"top":23897,"bottom":1}
```

## 获取特定日志的异常范围

```
curl ${server_address}/api/v1/fragments/{fid}/events/{eid}/logs/changepoints
```

- fid: fragment id，唯一代表一份日志
- eid: event id，一个 event 代表一类日志

返回示例:

```json
[{"start":1641177903,"stop":1641180011},{"start":1641176086,"stop":1641179409},{"start":1641182147,"stop":1641183961}]
```