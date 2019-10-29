# Mywater

A tool to fetch your water usage from a [Smart Energy Water] powered
portal. The current output is formatted to be uploaded to [InfluxDB]

## Usage

```sh
> DATE=10/20/2019 USERNAME=<your username> PASSWORD=<your password> go run main.go 2>/dev/null >data
> curl -XPOST http://influxdb:8086/write?db=whatever --data-binary @data
```

[Smart Energy Water]: https://www.smartenergywater.com/scm_water-conservation.html
[InfluxDB]: https://www.influxdata.com/
