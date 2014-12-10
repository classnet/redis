package redis

import ()

//一般命令

func (client *Client) Auth(password string) error {
	_, err := client.SendCommand("AUTH", password)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) Flush(all bool) error {
	var cmd string
	if all {
		cmd = "FLUSHALL"
	} else {
		cmd = "FLUSHDB"
	}
	_, err := client.SendCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

//发布/订阅

// 从出版商渠道，收到我们订阅的消息容器。
type Message struct {
	ChannelMatched string
	Channel        string
	Message        []byte
}

// 订阅redis服务channels,该方法将阻塞直到一个封闭的sub/unsub通道。
// 有两条渠道订阅/退订和psubscribe / punsubscribe。
// The former does an exact match on the channel, the later uses glob patterns on the redis channels.
// Closing either of these channels will unblock this method call.
// Messages that are received are sent down the messages channel.
func (client *Client) Subscribe(subscribe <-chan string, unsubscribe <-chan string, psubscribe <-chan string, punsubscribe <-chan string, messages chan<- Message) error {
	cmds := make(chan []string, 0)
	data := make(chan interface{}, 0)

	go func() {
		for {
			var channel string
			var cmd string

			select {
			case channel = <-subscribe:
				cmd = "SUBSCRIBE"
			case channel = <-unsubscribe:
				cmd = "UNSUBSCRIBE"
			case channel = <-psubscribe:
				cmd = "PSUBSCRIBE"
			case channel = <-punsubscribe:
				cmd = "UNPSUBSCRIBE"

			}
			if channel == "" {
				break
			} else {
				cmds <- []string{cmd, channel}
			}
		}
		close(cmds)
		close(data)
	}()

	go func() {
		for response := range data {
			db := response.([][]byte)
			messageType := string(db[0])
			switch messageType {
			case "message":
				channel, message := string(db[1]), db[2]
				messages <- Message{channel, channel, message}
			case "subscribe":
				// Ignore
			case "unsubscribe":
				// Ignore
			case "pmessage":
				channelMatched, channel, message := string(db[1]), string(db[2]), db[3]
				messages <- Message{channelMatched, channel, message}
			case "psubscribe":
				// Ignore
			case "punsubscribe":
				// Ignore

			default:
				// log.Printf("Unknown message '%s'", messageType)
			}
		}
	}()

	err := client.SendCommands(cmds, data)
	return err
}

// Publish a message to a redis server.
func (client *Client) Publish(channel string, val []byte) error {
	_, err := client.SendCommand("PUBLISH", channel, string(val))
	if err != nil {
		return err
	}
	return nil
}

//服务器的命令

func (client *Client) Save() error {
	_, err := client.SendCommand("SAVE")
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Bgsave() error {
	_, err := client.SendCommand("BGSAVE")
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Lastsave() (int64, error) {
	res, err := client.SendCommand("LASTSAVE")
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

func (client *Client) Bgrewriteaof() error {
	_, err := client.SendCommand("BGREWRITEAOF")
	if err != nil {
		return err
	}
	return nil
}
