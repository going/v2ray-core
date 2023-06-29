package proto

type UserModel struct {
	ID          int64           `db:"id"`
	Email       string          `db:"email"`
	UUID        string          `db:"uuid"`
	AlterID     uint32          `db:"AlterId"`
	TrafficRate float64         `db:"traffic_rate"`
	Traffics    *UserTrafficLog `db:"-"`
}

type UserTrafficLog struct {
	Email     string
	Uploads   int64
	Downloads int64
	Clients   int64
	IPs       []string
}

type DBConfig struct {
	Master  string `yaml:"master"`
	MaxOpen int    `yaml:"max_open"`
	MaxIdle int    `yaml:"max_idle"`
}

type NodeModel struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Server      string  `db:"server"`
	NodeClass   int64   `db:"node_class"`
	Port        int64   `db:"port"`
	VlessPort   int64   `db:"vlessport"`
	TrafficRate float64 `db:"traffic_rate"`
	Reality     bool    `db:"reality"`
}
