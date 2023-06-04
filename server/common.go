package server

type PerfEventConfig struct {
	Cgroup string `json:"cgroup" form:"cgroup"`
	Cpuset string `json:"cpuset" form:"cpuset"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
