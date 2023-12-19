# 生成TLS证书

## 1. docker方式

https://distribution.github.io/distribution/about/insecure/

```shell
# CNu与subjectAltName一致
openssl req \
-config /etc/ssl/openssl.cnf \
-newkey rsa:4096 -nodes -sha256 -keyout ./domain.key \
-addext "subjectAltName = DNS:myregistry.domain.com" \
-x509 -days 365 -out domain.crt
```

```shell
Generating a RSA private key
...............................................++++
.......................++++
writing new private key to './domain.key'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:TJ
Locality Name (eg, city) []:TJ
Organization Name (eg, company) [Internet Widgits Pty Ltd]:SS
Organizational Unit Name (eg, section) []:SS
Common Name (e.g. server FQDN or YOUR name) []:myregistry.domain.com                            
Email Address []:ss@ss.com
```

## 2. 常用方法

```shell
# 生成服务器私钥
openssl genrsa -out server.key 1024
# 根据私钥和输入的信息生成证书请求文件 查找系统openssl.cnf 文件 不同系统可能不一致
openssl req -config /etc/ssl/openssl.cnf -new -key server.key -out server.csr
# 用第一步的私钥和第二步的请求文件生成证书
openssl x509 -req -in server.csr -out server.crt -signkey server.key -days 3650
```