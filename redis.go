package redis

import (
    "bufio"
    "bytes"   
    "fmt"
    "io"
    "io/ioutil"
    "net"      
    "strings"  
    "strconv"
)


type Client struct {
    Addr        string                 //地址
    Db          int                    //数据库号
    Password    string                 //密码
    MaxPoolSize int                    //最大连接池大小
    pool chan net.Conn                 //连接池
}

//Redis错误信息设置
type RedisError string

func (err RedisError) Error() string { return "Redis Error: " + string(err) }

var doesNotExist = RedisError("Key does not exist ")


//读取块回复
func readBulk(reader *bufio.Reader, head string) ([]byte, error) {
    var err error
    var data []byte
    //空返回空错误
    if head == "" {
        head, err = reader.ReadString('\n')
        if err != nil {
            return nil, err
        }
    }
    switch head[0] {
    //第一个字符是":" 截取后面内容
    case ':':
        data = []byte(strings.TrimSpace(head[1:]))
    //第一个字符是"$" 转化为int
    case '$':
        size, err := strconv.Atoi(strings.TrimSpace(head[1:]))
        if err != nil {
            return nil, err
        }
        if size == -1 {
            return nil, doesNotExist
        }
        lr := io.LimitReader(reader, int64(size))
        data, err = ioutil.ReadAll(lr)
        if err == nil {
            // read end of line
            _, err = reader.ReadString('\n')
        }
    default:
        return nil, RedisError("Expecting Prefix'$'or':'")
    }

    return data, err
}

//写请求
func writeRequest(writer io.Writer, cmd string, args ...string) error {
    b := commandBytes(cmd, args...)
    _, err := writer.Write(b)
    return err
}

//输出命令行
func commandBytes(cmd string, args ...string) []byte {
    var cmdbuf bytes.Buffer
    fmt.Fprintf(&cmdbuf, "*%d\r\n$%d\r\n%s\r\n", len(args)+1, len(cmd), cmd)
    for _, s := range args {
        fmt.Fprintf(&cmdbuf, "$%d\r\n%s\r\n", len(s), s)
    }
    return cmdbuf.Bytes()
}

//读响应
func readResponse(reader *bufio.Reader) (interface{}, error) {

    var line string
    var err error

    //读取第一个非空白行，空返回空错误
    for {
        line, err = reader.ReadString('\n')
        if len(line) == 0 || err != nil {
            return nil, err
        }
	
        line = strings.TrimSpace(line)   //删除前后空格
        if len(line) > 0 {
            break
        }
    }  
    //第一个字符是"+" 截取后面内容
    if line[0] == '+' {
        return strings.TrimSpace(line[1:]), nil
    }

    if strings.HasPrefix(line, "-ERR ") {
        errmesg := strings.TrimSpace(line[5:])
        return nil, RedisError(errmesg)
    }
    //第一个字符是":" 转化为int64
    if line[0] == ':' {
        n, err := strconv.ParseInt(strings.TrimSpace(line[1:]), 10, 64)
        if err != nil {
            return nil, RedisError("Int reply is not a number")
        }
        return n, nil
    }
    //第一个字符是"*" 转化为int,如果是负数改为0
    if line[0] == '*' {
        size, err := strconv.Atoi(strings.TrimSpace(line[1:]))
        if err != nil {
            return nil, RedisError("MultiBulk reply expected a number")
        }
        if size <= 0 {
            return make([][]byte, 0), nil
        }
        res := make([][]byte, size)
        for i := 0; i < size; i++ {
            res[i], err = readBulk(reader, "")
            if err == doesNotExist {
                continue
            }
            if err != nil {
                return nil, err
            }
            // dont read end of line as might not have been bulk
        }
        return res, nil
    }
    return readBulk(reader, line)
}
//原始发送
func (client *Client) rawSend(c net.Conn, cmd []byte) (interface{}, error) {
    _, err := c.Write(cmd)
    if err != nil {
        return nil, err
    }

    reader := bufio.NewReader(c)

    data, err := readResponse(reader)
    if err != nil {
        return nil, err
    }

    return data, nil
}
//打开连接
func (client *Client) openConnection() (c net.Conn, err error) {

    addr := "127.0.0.1:6379"   //默认地址

    if client.Addr != "" {
        addr = client.Addr
    }
    c, err = net.Dial("tcp", addr)//在网络上连接地址addr，并返回一个Conn接口
    if err != nil {
        return
    }
    
    //处理密码验证
    if client.Password != "" {
        cmd := fmt.Sprintf("AUTH %s\r\n", client.Password)
        _, err = client.rawSend(c, []byte(cmd))
        if err != nil {
            return
        }
    }

    if client.Db != 0 {
        cmd := fmt.Sprintf("SELECT %d\r\n", client.Db)
        _, err = client.rawSend(c, []byte(cmd))
        if err != nil {
            return
        }
    }

    return
}

func (client *Client) SendCommand(cmd string, args ...string) (data interface{}, err error) {
    
    //抓取池中的连接
    var b []byte
    c, err := client.popCon()
    if err != nil {
        println(err.Error())
        goto End
    }

    b = commandBytes(cmd, args...)
    data, err = client.rawSend(c, b)
    if err == io.EOF {
        c, err = client.openConnection()
        if err != nil {
            println(err.Error())
            goto End
        }

        data, err = client.rawSend(c, b)
    }

End:

    //添加客户端返回到队列
    client.pushCon(c)

    return data, err
}

func (client *Client) SendCommands(cmdArgs <-chan []string, data chan<- interface{}) (err error) {
    // grab a connection from the pool
    c, err := client.popCon()
    var reader *bufio.Reader
    var pong interface{}
    var errs chan error
    var errsClosed = false

    if err != nil {
        goto End
    }

    reader = bufio.NewReader(c)

    // Ping first to verify connection is open
    err = writeRequest(c, "PING")

    // On first attempt permit a reconnection attempt
    if err == io.EOF {
        // Looks like we have to open a new connection
        c, err = client.openConnection()
        if err != nil {
            goto End
        }
        reader = bufio.NewReader(c)
    } else {
        // Read Ping response
        pong, err = readResponse(reader)
        if pong != "PONG" {
            return RedisError("Unexpected response to PING.")
        }
        if err != nil {
            goto End
        }
    }

    errs = make(chan error)

    go func() {
        for cmdArg := range cmdArgs {
            err = writeRequest(c, cmdArg[0], cmdArg[1:]...)
            if err != nil {
                if !errsClosed {
                    errs <- err
                }
                break
            }
        }
        if !errsClosed {
            errsClosed = true
            close(errs)
        }
    }()

    go func() {
        for {
            response, err := readResponse(reader)
            if err != nil {
                if !errsClosed {
                    errs <- err
                }
                break
            }
            data <- response
        }
        if !errsClosed {
            errsClosed = true
            close(errs)
        }
    }()

    // Block until errs channel closes
    for e := range errs {
        err = e
    }

End:

    // Close client and synchronization issues are a nightmare to solve.
    c.Close()

    // Push nil back onto queue
    client.pushCon(nil)

    return err
}

func (client *Client) popCon() (net.Conn, error) {
    if client.pool == nil {
        poolSize := client.MaxPoolSize
        if poolSize == 0 {
            poolSize =  5    //默认连接池大小
        }
        client.pool = make(chan net.Conn, poolSize)
        for i := 0; i < poolSize; i++ {
            //add dummy values to the pool
            client.pool <- nil
        }
    }
    // grab a connection from the pool
    c := <-client.pool

    if c == nil {
        return client.openConnection()
    }
    return c, nil
}

func (client *Client) pushCon(c net.Conn) {
    client.pool <- c
}
