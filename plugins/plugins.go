package plugins

import (
	"time"
)

type Protocol string
type Class uint

var ClassMap = map[Class]string{
	Unauthorized: "未授权",
	WeakPass:     "弱口令",
	UnKnow:       "未知",
}

/**单个执行任务结构**/
type Single struct {
	Ip       string
	Port     uint
	Protocol string
	Username string
	Password string
	TimeOut  time.Duration
}

type ScanResult struct {
	Single Single
	Class  Class //1.未授权  2.弱口令 3.未知
	Result bool
}

type ScanFunc func(i interface{}) interface{}
type UnauthorizedFunc func(i interface{}) interface{}

// Scan 扫描协议定义包含执行的协程池及其不同协议用到的参数
type Scan struct {
	Thread           int
	Ts               time.Duration
	ScanFunc         ScanFunc
	UnauthorizedFunc UnauthorizedFunc
}

var (
	ScanMap map[Protocol]*Scan
)

const (
	SshThread           = 3 //ssh并发数
	DefaultThread       = 200
	DefaultTs           = time.Second * 3
	Unauthorized  Class = 1
	WeakPass      Class = 2
	UnKnow        Class = 3
)

const (
	COUCHDB    Protocol = "couchdb"
	DOCKER     Protocol = "docker"
	ELASTIC    Protocol = "elastic"
	FTP        Protocol = "ftp"
	HADOOP     Protocol = "hadoop"
	HIVE       Protocol = "hive"
	KIBANA     Protocol = "kibana"
	LDAP       Protocol = "ldap"
	MEMCACHE   Protocol = "memcache"
	MONGODB    Protocol = "mongodb"
	MSSQL      Protocol = "mssql"
	MYSQL      Protocol = "mysql"
	ORACLE     Protocol = "oracle"
	POSTGRESQL Protocol = "postgresql"
	RDP        Protocol = "rdp"
	REDIS      Protocol = "redis"
	SMB        Protocol = "smb"
	SNMP       Protocol = "snmp"
	SOLR       Protocol = "solr"
	SSH        Protocol = "ssh"
	TELNET     Protocol = "telnet"
	TOMCAT     Protocol = "tomcat"
	DM         Protocol = "dm"
	VNC        Protocol = "vnc"
	AMQP       Protocol = "amqp"
	MQTT       Protocol = "mqtt"
	HTTP       Protocol = "http"
	POP3       Protocol = "pop3"
	SOCKS5     Protocol = "socks5"
	RSYNC      Protocol = "rsync"
	ZOOKEEPER  Protocol = "zookeeper"
	NEUTRON    Protocol = "neutron"
)

/**初始化构建各个执行协议的调用池供引擎调用*/
func init() {
	ScanMap = map[Protocol]*Scan{
		COUCHDB:    {Thread: DefaultThread, ScanFunc: ScanCouchdb, UnauthorizedFunc: UnauthorizedCouchdb},
		DOCKER:     {Thread: DefaultThread, ScanFunc: ScanDocker, UnauthorizedFunc: UnauthorizedDocker},
		ELASTIC:    {Thread: DefaultThread, ScanFunc: ScanElastic, UnauthorizedFunc: UnauthorizedElastic},
		FTP:        {Thread: DefaultThread, ScanFunc: ScanFtp, UnauthorizedFunc: UnauthorizedFtp},
		HADOOP:     {Thread: DefaultThread, ScanFunc: ScanHadoop, UnauthorizedFunc: UnauthorizedHadoop},
		HIVE:       {Thread: DefaultThread, ScanFunc: ScanHive, UnauthorizedFunc: UnauthorizedHive},
		KIBANA:     {Thread: DefaultThread, ScanFunc: ScanKibana, UnauthorizedFunc: UnauthorizedKibana},
		LDAP:       {Thread: DefaultThread, ScanFunc: ScanLdap, UnauthorizedFunc: UnauthorizedLdap},
		MEMCACHE:   {Thread: DefaultThread, ScanFunc: ScanMemcache, UnauthorizedFunc: UnauthorizedMemcache},
		MONGODB:    {Thread: DefaultThread, ScanFunc: ScanMongodb, UnauthorizedFunc: UnauthorizedMongodb},
		MYSQL:      {Thread: DefaultThread, ScanFunc: ScanMysql, UnauthorizedFunc: UnauthorizedMysql},
		POSTGRESQL: {Thread: DefaultThread, ScanFunc: ScanPostgres, UnauthorizedFunc: UnauthorizedPostgres},
		RDP:        {Thread: DefaultThread, ScanFunc: ScanRdp, UnauthorizedFunc: UnauthorizedRdp},
		REDIS:      {Thread: DefaultThread, ScanFunc: ScanRedis, UnauthorizedFunc: UnauthorizedRedis},
		SMB:        {Thread: DefaultThread, ScanFunc: ScanSmb, UnauthorizedFunc: UnauthorizedSmb},
		SNMP:       {Thread: DefaultThread, ScanFunc: ScanSNMP, UnauthorizedFunc: UnauthorizedSnmp},
		SOLR:       {Thread: DefaultThread, ScanFunc: ScanSolr, UnauthorizedFunc: UnauthorizedSolr},
		SSH:        {Thread: SshThread, ScanFunc: ScanSsh, UnauthorizedFunc: UnauthorizedSsh},
		TELNET:     {Thread: DefaultThread, ScanFunc: ScanTelnet, UnauthorizedFunc: UnauthorizedTelnet},
		TOMCAT:     {Thread: DefaultThread, ScanFunc: ScanTomcat, UnauthorizedFunc: UnauthorizedTomcat},
		DM:         {Thread: DefaultThread, ScanFunc: ScanDm, UnauthorizedFunc: UnauthorizedDm},
		MSSQL:      {Thread: DefaultThread, ScanFunc: ScanMssql, UnauthorizedFunc: UnauthorizedMssql},
		ORACLE:     {Thread: DefaultThread},
		VNC:        {Thread: DefaultThread, ScanFunc: ScanVnc, UnauthorizedFunc: UnauthorizedVnc},
		AMQP:       {Thread: DefaultThread, ScanFunc: ScanAmqp, UnauthorizedFunc: UnauthorizedAmqp},
		MQTT:       {Thread: DefaultThread, ScanFunc: ScanMqtt, UnauthorizedFunc: UnauthorizedMqtt},
		POP3:       {Thread: DefaultThread, ScanFunc: ScanPop3, UnauthorizedFunc: UnauthorizedPop3},
		SOCKS5:     {Thread: DefaultThread, ScanFunc: ScanSocks5, UnauthorizedFunc: UnauthorizedSocks5},
		RSYNC:      {Thread: DefaultThread, ScanFunc: ScanVnc, UnauthorizedFunc: UnauthorizedVnc},
		ZOOKEEPER:  {Thread: DefaultThread, ScanFunc: ScanZookeeper, UnauthorizedFunc: UnauthorizedZookeeper},
		NEUTRON:    {Thread: DefaultThread, ScanFunc: ScanVnc, UnauthorizedFunc: UnauthorizedVnc},
		HTTP:       {Thread: DefaultThread, ScanFunc: ScanVnc, UnauthorizedFunc: UnauthorizedVnc},
	}
}
