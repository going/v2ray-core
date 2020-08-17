package main

import (
	"flag"

	"v2ray.com/core/watchman"
)

var (
	address    = flag.String("address", "127.0.0.1:4321", "v2ray gRPC server address")
	inboundTag = flag.String("tag", "vmess-proxy", "v2ray inbound tag name")
	dbUrl      = flag.String("db", "222kingshard333:sddssds3322we@tcp(45.76.99.66:1795)/dotunnel005", "database address")
	nodeId     = flag.Int64("node", 1, "node id")
)

func main() {
	flag.Parse()

	watchman.Start(*address, *inboundTag, *dbUrl, *nodeId)
}
