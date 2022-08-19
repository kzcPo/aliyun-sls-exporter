# Aliyun SlsMonitor Exporter

exporter for Aliyun SlS Monitor. Written in Golang.
inspired by [aliyun-exporter](https://github.com/IAOTW/aliyun-exporter)

## Develop

```bash
cd aliyun-sls-exporter
make tidy
```

## Usage

```bash
# build
make build
# run
./aliyun-sls-exporter serve
```

## 1. Based on configuration files
Provide a configuration file containing authentication information.
```yaml
credentials:
  # You can obtain monitoring data of multiple tenants by configuring multiple Tenant information.
  tenantId1:
    accessKey: xxxxxxxxxxxx
    accessKeySecret: xxxxxxxxxxxx
    region: cn-hangzhou
  tenantId2:
    accessKey: xxxxxxxxxxxx
    accessKeySecret: xxxxxxxxxxxx
    region: cn-hangzhou
```

## 2. configuration sls params
Variable description

| Name                      | Description                                       |
|---------------------------|---------------------------------------------------|
| `project`                 | sls request params project                        |
| `logstore`                | sls request params logstore                       |
| `name`                    | metrics name                                      |
| `desc`                    | metrics Help                                      |
| `query`                   | sls request params query, Query statement         |
| `dimensions`              | Query statement result copy to metrics label name |
| `measure`                 | metrics value label                               |

```yaml
logsMetric:
  taqu:
    - project: taqu
      logstore: nginx_logdb
      name: system_nginx_requests_count
      desc: all level, group by system uri
      query: env:online  | select "system", pv from( select count(1) as pv , "system" from log group by "system" order by pv desc) order by pv desc
      dimensions:
        - system
      measure: pv
```
You can visit metrics in http://aliyun-sls-exporter:9528/metrics


## Ref

- https://next.api.aliyun.com/api/Sls/2020-12-30/GetLogs?sdkStyle=dara&params={}&tab=DEBUG
- https://help.aliyun.com/document_detail/29029.html
- https://github.com/fengxsong/aliyun-exporter
- https://github.com/aliyun/alibaba-cloud-sdk-go
