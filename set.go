package redis

import (    
     
)

// Set命令

func (client *Client) Sadd(key string, value []byte) (bool, error) {
    res, err := client.SendCommand("SADD", key, string(value))

    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Srem(key string, value []byte) (bool, error) {
    res, err := client.SendCommand("SREM", key, string(value))

    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Spop(key string) ([]byte, error) {
    res, err := client.SendCommand("SPOP", key)
    if err != nil {
        return nil, err
    }

    if res == nil {
        return nil, RedisError("Spop failed")
    }

    data := res.([]byte)
    return data, nil
}

func (client *Client) Smove(src string, dst string, val []byte) (bool, error) {
    res, err := client.SendCommand("SMOVE", src, dst, string(val))
    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Scard(key string) (int, error) {
    res, err := client.SendCommand("SCARD", key)
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Sismember(key string, value []byte) (bool, error) {
    res, err := client.SendCommand("SISMEMBER", key, string(value))

    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Sinter(keys ...string) ([][]byte, error) {
    res, err := client.SendCommand("SINTER", keys...)
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Sinterstore(dst string, keys ...string) (int, error) {
    args := make([]string, len(keys)+1)
    args[0] = dst
    copy(args[1:], keys)
    res, err := client.SendCommand("SINTERSTORE", args...)
    if err != nil {
        return 0, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Sunion(keys ...string) ([][]byte, error) {
    res, err := client.SendCommand("SUNION", keys...)
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Sunionstore(dst string, keys ...string) (int, error) {
    args := make([]string, len(keys)+1)
    args[0] = dst
    copy(args[1:], keys)
    res, err := client.SendCommand("SUNIONSTORE", args...)
    if err != nil {
        return 0, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Sdiff(key1 string, keys []string) ([][]byte, error) {
    args := make([]string, len(keys)+1)
    args[0] = key1
    copy(args[1:], keys)
    res, err := client.SendCommand("SDIFF", args...)
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Sdiffstore(dst string, key1 string, keys []string) (int, error) {
    args := make([]string, len(keys)+2)
    args[0] = dst
    args[1] = key1
    copy(args[2:], keys)
    res, err := client.SendCommand("SDIFFSTORE", args...)
    if err != nil {
        return 0, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Smembers(key string) ([][]byte, error) {
    res, err := client.SendCommand("SMEMBERS", key)

    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Srandmember(key string) ([]byte, error) {
    res, err := client.SendCommand("SRANDMEMBER", key)
    if err != nil {
        return nil, err
    }

    return res.([]byte), nil
}
