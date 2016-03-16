//index.go
package main

import (
	"./model"
	"net/http"
)

func main() {

	http.HandleFunc("/reg", model.Reg) //设置访问的路由
	http.HandleFunc("/login", model.Login)
	http.ListenAndServe(":85", nil) //设置监听的端口
}

Examples

Most of the examples connect to a redis database running in the default port -- 6379. 

Hello World example

package main

import "github.com/hoisie/redis"

func main() {
    var client redis.Client
    var key = "hello"
    client.Set(key, []byte("world"))
    val, _ := client.Get("hello")
    println(key, string(val))
}

Strings

var client redis.Client
client.Set("a", []byte("hello"))
val, _ := client.Get("a")
println(string(val))
client.Del("a")

Lists

var client redis.Client
vals := []string{"a", "b", "c", "d", "e"}
for _, v := range vals {
    client.Rpush("l", []byte(v))
}
dbvals,_ := client.Lrange("l", 0, 4)
for i, v := range dbvals {
    println(i,":",string(v))
}
client.Del("l")

Publish/Subscribe

sub := make(chan string, 1)
sub <- "foo"
messages := make(chan Message, 0)
go client.Subscribe(sub, nil, nil, nil, messages)

time.Sleep(10 * 1000 * 1000)
client.Publish("foo", []byte("bar"))

msg := <-messages
println("received from:", msg.Channel, " message:", string(msg.Message))

close(sub)
close(messages)

