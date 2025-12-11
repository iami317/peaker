# peaker

一款高质量的弱口令扫描工具,目前支持市面大部分协议31种
支持的协议
1. amqp
2. couchdb 
3. 达梦数据库
4. docker_api
5. elastic
6. ftp
7. hadoop 
8. hive
9. kibana
10. ldap 
11. memcache 
12. mongodb
13. mqtt
14. mssql 
15. mysql
16. neutron
17. oracle 
18. pop3
19. pgsql 
20. rdp 
21. redis
22. rsync
23. smb 
24. snmp 
25. socks5
26. solr 
27. ssh 
28. telnet 
29. tomcat
30. vnc
31. zookeeper
32. sqlserver // todo
33. WinRM // todo
34. SVN // todo
35. WebLogic // todo
36. jboss // todo
37. wmi // todo
38. smtp // todo
39. iamp // todo
40. webdav // todo

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
    go build -o peaker.exe ./_example/main.go
### linux: 
    set CGO_ENABLED=0 && set GOOS=linux && set GOARCH=amd64  
    go build -o peaker.exe ./_example/main.go
### mac:
    set CGO_ENABLED=0 && set GOOS=darwin && set GOARCH=amd64  
    go build -o peaker.exe ./_example/main.go
