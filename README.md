<code>

package main<br />

import "github.com/classnet/redis"<br />

func main() {<br />
    var client redis.Client<br />
    var key = "hello"<br />
    client.Set(key, []byte("world"))<br />
    val, _ := client.Get("hello")<br />
    println(key, string(val))<br />
}
</code>
