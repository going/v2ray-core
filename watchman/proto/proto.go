package proto

type UserModel struct {
	ID       int64           `db:"id"`
	Email    string          `db:"email"`
	UUID     string          `db:"uuid"`
	AlterID  uint32          `db:"AlterId"`
	Password string          `db:"passwd"`
	Method   string          `db:"method"`
	Port     uint16          `db:"port"`
	Traffics *UserTrafficLog `db:"-"`
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
