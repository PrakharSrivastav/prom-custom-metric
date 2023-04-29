package main

const (
	template    = "sftp_outbound_status{name=\"%s\",dimension=\"connectivity\"} %d"
	logfilePath = "/u01/oracle/batch/Run_scripts/sftp_conn.log"
	//logfilePath = "sftp_conn.log"
	promFilePath = "/var/cache/prometheus/sftp_conn.prom"
	//promFilePath = "sftp_conn.prom"
	configPath = "/u01/oracle/batch/Run_scripts/metrics/config.yaml"
	//configPath = "config.yaml"
)

type Config struct {
	Items []struct {
		Name string `yaml:"name"`
		Exe  string `yaml:"exe"`
		Args string `yaml:"args"`
	} `yaml:"items"`
}
