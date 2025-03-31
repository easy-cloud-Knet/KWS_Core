package userconfig



type User_specific struct {
	Name                string   `yaml:"name,omitempty"`
	Passwd              string   `yaml:"passwd,omitempty"`
	Lock_passwd         bool     `yaml:"lock_passwd"`
	Ssh_authorized_keys []string `yaml:"ssh_authorized_keys,omitempty"`
	Groups              string   `yaml:"groups,omitempty"`
	SuGroup             string   `yaml:"sudo,omitempty"`
	Shell               string   `yaml:"shell,omitempty"`
}


type User_write_file struct {
	Path        string `yaml:"path"`
	Permissions string `yaml:"permissions"`
	Content     string `yaml:"content"`
}
type User_data_yaml struct {
	PackageUpdatable bool              `yaml:"package_update"`
	PredownProjects  []string          `yaml:"packages"`
	Users            []interface{}     `yaml:"users"`
	Write_files      []User_write_file `yaml:"write_files"`
	Runcmd           []string          `yaml:"runcmd"`
}



type Meta_data_yaml struct {
	Instance_ID   string `yaml:"instance-id"`
	Local_Host_Id string `yaml:"local-hostname"`
}
