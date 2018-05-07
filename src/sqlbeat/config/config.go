// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period   time.Duration `config:"period"`
	DBType   string        `yaml:"dbtype"`
	Hostname string        `yaml:"hostname"`
	Port     string        `yaml:"port"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Database string        `yaml:"database"`
	Queries  string        `yaml:"queries"`
}

var DefaultConfig = Config{
	Period:   5 * time.Second,
	DBType:   "mssql",
	Hostname: "VimalKumarPC",
	Database: "TestDB",
	Username: "sa",
	Port:     "1433",
	Password: "tesla",
	Queries:  "SELECT  [Id] ,[Name] ,[ModifiedOn]  ,[calcol] FROM [TestDB].[dbo].[Category]",
}
