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

var filterMethodList = make([]func(*HttpContext, *Filter), 0)

func createFilterChain() *Filter {
	chain := new(Filter)
	current := chain
	for _, filterMethod := range filterMethodList {
		newFilter := &Filter{filterMethod: filterMethod, next: nil}
		current.addFilterAt(newFilter)
		current = current.next
	}
	return chain
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
