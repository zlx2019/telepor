## Socks5
- 参考 文献[RFC 1928](https://datatracker.ietf.org/doc/rfc1928/)
> SOCKS is an Internet protocol that exchanges network packets between a client and server through a proxy server

Socks是一种网络传输协议，该协议可以通过一个**中间服务器(代理)**来实现**客户端**和**目标服务端**之间的通信传输。简单来说它就是一个代理协议，扮演一个无情传话人。

### Socks5代理流程

流程分为四个步骤，如下:

- 握手
  
- 协商 && 认证
- 连接
- 转发

简述流程：

1. 首先**客户端**与**代理服务器**建立连接，完成TCP三次握手。

2. 连接完成后，双方开始协商，**客户端**告诉**代理服务器**自己支持的认证方式等信息，报文格式如下：

   ```tex
   +----+----------+----------+
   |VER | NMETHODS | METHODS  |
   +----+----------+----------+
   | 1  |    1     | 1 to 255 |
   +----+----------+----------+
   ```

   - `VER`:  协议版本，`0x05`表示Socks5。

   - `METHODS`：客户端支持的认证方式列表。即一个`uint8`(1-255)类型的数组。可选项有如下:

     ```tex
     `0x00`: 无需认证。
     `0x01`: GSSAPI认证。
     `0x02`: 用户名密码认证。
     `0x03`: IANA认证。
     `0x80`: 保留的认证方式。
     `0xFF`: 不支持任何认证方式，收到该值后，应该立刻关闭连接.
     ```

   - `NMETHODS`: 客户端支持多少种认证方式。即`METHODS`字段的长度。

3. **代理服务器**获取报文后，从报文中的`METHODS`选取一项认证方式，响应给**客户端**，表示使用该方式认证，报文格式如下:

   ```tex
   +----+--------+
   |VER | METHOD |
   +----+--------+
   | 1  |   1    |
   +----+--------+
   ```

   - `VER`: 协议版本。
   - `METHOD`: 协商后，要使用的认证方式。

4. 然后双方通过选定的方式进行认证，认证通过后，**客户端**就可以向**代理服务器**发送代理请求，代理请求的报文格式如下:UD

   ```tex
   +----+-----+-------+------+----------+----------+
   |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
   +----+-----+-------+------+----------+----------+
   | 1  |  1  | X'00' |  1   | Variable |    2     |
   +----+-----+-------+------+----------+----------+
   ```

   - `VER`: 协议版本号 (1Byte)。
   - `CMD`:  代理请求的指令码 (1Byte)。
     - `0x01`:  表示`CONNECT`请求，用于TCP代理。
     - `0x02`: 表示`BIND`请求。
     - `0x03`:  用于UDP代理。
   - `RSV`：保留字段，设置为`0x00`即可 (1Byte)。
   - `ATYP`：目标服务器地址类型 (1Byte)。
     - `0x01`：IPv4。
     - `0x03`：表示不是IP，而是域名(并且`DST.ADDR`的首个字节为域名的长度)。
     - `0x04`：IPv6。
   - `DST.ADDR`：目标服务器地址( [4|16] Byte)。
   - `DST.PORT`：目的地端口(2Byte)。

5. **代理服务器**解析该报文，得到目标服务器信息，然后与**目标服务器**尝试建立连接，然后响应客户端，报文格式如下：

   ```tex
   +----+-----+-------+------+----------+----------+
   |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
   +----+-----+-------+------+----------+----------+
   | 1  |  1  | X'00' |  1   | Variable |    2     |
   +----+-----+-------+------+----------+----------+
   ```

   - `VER`:  协议版本号。
   - `REP`： 代理请求结果，结果有如下：
     - `0x00`：连接成功。
     - `0x01`： 目标服务器故障。
     - `0x02`： 目标服务器不允许的连接。
     - `0x03`： 网络不可达。
     - `0x04`：目标主机无法访问。
     - `0x05`：连接被拒绝。
     - `0x06`： TTL 过期。
     - `0x07`：不支持的命令。
     - `0x08`：不支持的地址类型。
     - `0xFF`：未知的错误。

   - `RSV`：预留字节。
   - `ATYP`：`BND.ADDR`地址类型。
   - `BND.ADDR`：代理服务器连接目标服务器成功后的代理服务器 IP。
   - `BND.PORT`: 代理服务器连接目标服务器成功后的代理服务器端口。

6. **代理服务器**与双方桥梁搭建完毕，开始通信。

![Socks5流程图.png](Socks5%E6%B5%81%E7%A8%8B%E5%9B%BE.png)