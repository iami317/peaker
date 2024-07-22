# peaker

一款高质量的弱口令扫描工具
支持的协议
1. couchdb 
2. 达梦数据库
3. docker_api 
4. elasticsearch 
5. ftp
6. hadoop 
7. hive(暂不支持)
8. kibana (暂不支持)
9. ldap 
10. memcache 
11. mongodb 
12. mssql 
13. mysql 
14. oracle (暂不支持)
15. pgsql 
16. rdp 
17. redis 
18. smb 
19. snmp 
20. solr 
21. ssh 
22. telnet 
23. tomcat
24. sqlserver //todo
25. WinRM
26. VNC
27. SVN
28. WebLogic
29. jboss
30. zookeeper
31. wmi
32. smtp
33. pop3
34. iamp
35. webdav

使用方法参考 _example/main.go

参数说明：
1. --ip_list -i 目标文件地址 默认：iplist.txt，目标文件格式 192.168.103.156:22|SSH  不支持的协议会忽略
2. --user_dict -u 账号文件地址 默认：user.dic
3. --pass_dict -p 密码文件地址 默认：pass.dic
4. --check_alive -cA 运行是否检测目标是否存活
5. --verbose 是否详细展示
6. --thread -c 目标并发数量 默认：30
7. --timeout -t 单个目标最大执行时间 默认：20 * 60 秒
8. --timeout-single -tS 执行单个 ip port user pass 最大执行时间 默认：3 秒
9. --thread-single -tC 执行单个协议的并发数，例如：执行ssh 同时执行3组账号 密码
   
 ## 构建编译
(如果windows下执行错误，单个设置环境变量)
### windows：
    go build -o peaker.exe ./cmd/main.go
### linux: 
    set CGO_ENABLED=0 && set GOOS=linux && set GOARCH=amd64  
    go build -o peaker.exe ./cmd/main.go
### mac:
    set CGO_ENABLED=0 && set GOOS=darwin && set GOARCH=amd64  
    go build -o peaker.exe ./cmd/main.go
