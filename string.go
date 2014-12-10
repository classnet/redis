package redis

import (
 
    "strconv"
    
)


// String命令

func (client *Client) Set(key string, val []byte) error {
    _, err := client.SendCommand("SET", key, string(val))

    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Get(key string) ([]byte, error) {
    res, _ := client.SendCommand("GET", key)
    if res == nil {
        return nil, RedisError("Key `" + key + "` does not exist")
    }

    data := res.([]byte)
    return data, nil
}

func (client *Client) Getset(key string, val []byte) ([]byte, error) {
    res, err := client.SendCommand("GETSET", key, string(val))

    if err != nil {
        return nil, err
    }

    data := res.([]byte)
    return data, nil
}

func (client *Client) Mget(keys ...string) ([][]byte, error) {
    res, err := client.SendCommand("MGET", keys...)
    if err != nil {
        return nil, err
    }

    data := res.([][]byte)
    return data, nil
}

func (client *Client) Setnx(key string, val []byte) (bool, error) {
    res, err := client.SendCommand("SETNX", key, string(val))

    if err != nil {
        return false, err
    }
    if data, ok := res.(int64); ok {
        return data == 1, nil
    }
    return false, RedisError("Unexpected reply to SETNX")
}

func (client *Client) Setex(key string, time int64, val []byte) error {
    _, err := client.SendCommand("SETEX", key, strconv.FormatInt(time, 10), string(val))

    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Mset(mapping map[string][]byte) error {
    args := make([]string, len(mapping)*2)
    i := 0
    for k, v := range mapping {
        args[i] = k
        args[i+1] = string(v)
        i += 2
    }
    _, err := client.SendCommand("MSET", args...)
    if err != nil {
        return err
    }
    return nil
}

func (client *Client) Msetnx(mapping map[string][]byte) (bool, error) {
    args := make([]string, len(mapping)*2)
    i := 0
    for k, v := range mapping {
        args[i] = k
        args[i+1] = string(v)
        i += 2
    }
    res, err := client.SendCommand("MSETNX", args...)
    if err != nil {
        return false, err
    }
    if data, ok := res.(int64); ok {
        return data == 0, nil
    }
    return false, RedisError("Unexpected reply to MSETNX")
}

func (client *Client) Incr(key string) (int64, error) {
    res, err := client.SendCommand("INCR", key)
    if err != nil {
        return -1, err
    }

    return res.(int64), nil
}

func (client *Client) Incrby(key string, val int64) (int64, error) {
    res, err := client.SendCommand("INCRBY", key, strconv.FormatInt(val, 10))
    if err != nil {
        return -1, err
    }

    return res.(int64), nil
}

func (client *Client) Decr(key string) (int64, error) {
    res, err := client.SendCommand("DECR", key)
    if err != nil {
        return -1, err
    }

    return res.(int64), nil
}

func (client *Client) Decrby(key string, val int64) (int64, error) {
    res, err := client.SendCommand("DECRBY", key, strconv.FormatInt(val, 10))
    if err != nil {
        return -1, err
    }

    return res.(int64), nil
}

func (client *Client) Append(key string, val []byte) error {
    _, err := client.SendCommand("APPEND", key, string(val))

    if err != nil {
        return err
    }

    return nil
}

func (client *Client) Substr(key string, start int, end int) ([]byte, error) {
    res, _ := client.SendCommand("SUBSTR", key, strconv.Itoa(start), strconv.Itoa(end))

    if res == nil {
        return nil, RedisError("Key `" + key + "` does not exist")
    }

    data := res.([]byte)
    return data, nil
}

func (client *Client) Strlen(key string) (int, error) {
    res, err := client.SendCommand("STRLEN", key)
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}
