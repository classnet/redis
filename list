package redis

import (     
    "strconv"
)

// List命令

func (client *Client) Rpush(key string, val []byte) error {
    _, err := client.SendCommand("RPUSH", key, string(val))

    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Lpush(key string, val []byte) error {
    _, err := client.SendCommand("LPUSH", key, string(val))

    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Llen(key string) (int, error) {
    res, err := client.SendCommand("LLEN", key)
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Lrange(key string, start int, end int) ([][]byte, error) {
    res, err := client.SendCommand("LRANGE", key, strconv.Itoa(start), strconv.Itoa(end))
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Ltrim(key string, start int, end int) error {
    _, err := client.SendCommand("LTRIM", key, strconv.Itoa(start), strconv.Itoa(end))
    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Lindex(key string, index int) ([]byte, error) {
    res, err := client.SendCommand("LINDEX", key, strconv.Itoa(index))
    if err != nil {
        return nil, err
    }

    return res.([]byte), nil
}

func (client *Client) Lset(key string, index int, value []byte) error {
    _, err := client.SendCommand("LSET", key, strconv.Itoa(index), string(value))
    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Lrem(key string, count int, value []byte) (int, error) {
    res, err := client.SendCommand("LREM", key, strconv.Itoa(count), string(value))
    if err != nil {
        return -1, err
    }
    return int(res.(int64)), nil
}

func (client *Client) Lpop(key string) ([]byte, error) {
    res, err := client.SendCommand("LPOP", key)
    if err != nil {
        return nil, err
    }

    return res.([]byte), nil
}

func (client *Client) Rpop(key string) ([]byte, error) {
    res, err := client.SendCommand("RPOP", key)
    if err != nil {
        return nil, err
    }

    return res.([]byte), nil
}

func (client *Client) Blpop(keys []string, timeoutSecs uint) (*string, []byte, error) {
    return client.bpop("BLPOP", keys, timeoutSecs)
}
func (client *Client) Brpop(keys []string, timeoutSecs uint) (*string, []byte, error) {
    return client.bpop("BRPOP", keys, timeoutSecs)
}

func (client *Client) bpop(cmd string, keys []string, timeoutSecs uint) (*string, []byte, error) {
    args := append(keys, strconv.FormatUint(uint64(timeoutSecs), 10))
    res, err := client.SendCommand(cmd, args...)
    if err != nil {
        return nil, nil, err
    }
    kv := res.([][]byte)
    // Check for timeout
    if len(kv) != 2 {
        return nil, nil, nil
    }
    k := string(kv[0])
    v := kv[1]
    return &k, v, nil
}

func (client *Client) Rpoplpush(src string, dst string) ([]byte, error) {
    res, err := client.SendCommand("RPOPLPUSH", src, dst)
    if err != nil {
        return nil, err
    }

    return res.([]byte), nil
}
