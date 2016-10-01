goredis

Installation

Just run  go get github.com/classnet/redis 

Examples

Most of the examples connect to a redis database running in the default port -- 6379. 

Hello World example

package main

import "github.com/classnet/redis"

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

More examples coming soon. See  redis_test.go  for more usage examples.

Commands not supported yet
•MULTI/EXEC/DISCARD/WATCH/UNWATCH
•SORT
•ZUNIONSTORE / ZINTERSTORE
