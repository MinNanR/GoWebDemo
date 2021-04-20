package main

type Server struct {
	Port string
}

var server = Server{Port: "8989"}

type DataSource struct {
	Host       string
	Port       string
	DataBase   string
	DriverName string
	Username   string
	Password   string
}

var dataSource = DataSource{
	Host:       "minnan.site",
	Port:       "3306",
	DataBase:   "link_manager",
	Username:   "Minnan",
	Password:   "minnan",
	DriverName: "mysql",
}

var filterList = make([]func(*HttpContext, *FilterChain), 0)

func createFilterChain() FilterChain {
	return FilterChain{
		chain:      filterList,
		chainIndex: -1,
	}
}

var authorityPath = []string{"/login"}

type Redis struct {
	Host     string
	Port     string
	Database int
	Password string
}

var redisConfig = Redis{
	Host:     "minnan.site",
	Port:     "6379",
	Database: 5,
	Password: "minnan",
}
