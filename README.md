goredis

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
