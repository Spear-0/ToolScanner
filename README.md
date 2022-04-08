### http yaml编写
1、name(必须) 当前yaml名称

2、protect(必须) 使用的协议（http）

3、tool

    tool_name(必须) 检测工具的名称
    
    tool_version(选填) 检测工具的版本

4、request:

    method(必须):请求方式
    
    path(选填):请求路径
    
    data(选填):数据

5、response: (一下两个选项选择必须要有一个以上不为空)

    pcre_body:匹配http body部分
    
    pcre_header:匹配header格式为key: value 
    
    pcre_status:匹配响应首行

### tcp yaml编写

1、name

2、ptotect

3、tool

    tool_name

    tool_version
4、request

    data: 请求数据

5、response

    pcre_body: 相应匹配规则，默认接受2048字节的数据