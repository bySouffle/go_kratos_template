# Auth

## 1. JSON Web Token (JWT)

### 格式  AAAAA.BBBBB.CCCCC

* 三个子串分别称作 头部（Header） 、 负载（Payload） 和 签名（Signature）

#### 头部 Header

虽说只要相关几方之间有共识，则在头部中放什么是没有限制的，但通常由两部分组成：

* typ: 表示令牌类型（type），值为 ~JWT~
* alg: 表示签名此令牌的算法（algorithm），如 ~HMAC~、~RSA~、~SHA~

#### 负载 Payload

JWT 的第二部分表示负载，这部分由声明（claims）组成。

* iss (issuer): 声明了发行人，也就是发行 JWT 的主体。处理此声明通常是因应用而异的。“iss” 值是一个大小写敏感的字符串，包含一个普通字符串或者一个
  URL。该声明是可选的
* sub (subject): 表示 JWT 的主体 (用户)。值必须要么是全局唯一的，要么在发行人上下文范围内局部唯一。处理该声明通常也是因应用而异的。“sub”
  值是一个大小写敏感的字符串，包含一个普通字符串或者一个 URL。该声明是可选的
* aud (audience): 表示 JWT 的目标接收方。如果当该声明存在且处理该声明的一方不能通过 “aud” 的值进行自我身份验证时，则 JWT
  必须被拒绝。大多数情况下，这个值是由大小写敏感的字符串（包含一个普通字符串或者一个 URL）组成的数组。该声明是可选的
* exp (expiration): 表示过期时间，即等于或晚于那个时刻再处理 JWT 则绝不可被接受。其值通常是以秒记的时间戳（译注：按 POSIX
  中定义的 “seconds since epoch” 标准，也就是 PHP 等语言中常用的那种）。该声明是可选的
* nbf (not before) : 表示一个时间，即早于那个时刻再处理 JWT 则绝不可被接受。 其值通常是以秒记的时间戳。该声明是可选的
* iat (issued at): 表示发出 JWT 的时刻。可用于判断 JWT 的寿命。必须是一个时间戳。该声明是可选的
* jti (JWT ID): 为 JWT 提供一个唯一的身份识别符，其值必须难以重复，以防 JWT 被重复执行。该声明是可选的
* 自定义声明: 自定义key value

#### 签名

签名先是通过对头部和负载 Base64 编码而生成，其后会与一个密钥联合，最好被头部中指定的算法签名。

### HS256加密

* 对称加密，由于头部和负载仅仅对数据进行base64编码，在线jwt验证可以解码出信息所，因此在负载处不要存放敏感信息。
* 校验部分在签名处