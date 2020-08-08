module logagent

go 1.13

require (
	github.com/Shopify/sarama v1.26.4
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/hpcloud/tail v1.0.0
	go.etcd.io/etcd v3.3.22+incompatible
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.14.1

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
