package redis

import (
    "bytes"
    "strconv"
)

//key命令
func (client *Client) Exists(key string) (bool, error) {
    res, err := client.SendCommand("EXISTS", key)
    if err != nil {
        return false, err
    }
    return res.(int64) == 1, nil
}

func (client *Client) Del(key string) (bool, error) {
    res, err := client.SendCommand("DEL", key)

    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Type(key string) (string, error) {
    res, err := client.SendCommand("TYPE", key)

    if err != nil {
        return "", err
    }

    return res.(string), nil
}

func (client *Client) Keys(pattern string) ([]string, error) {
    res, err := client.SendCommand("KEYS", pattern)

    if err != nil {
        return nil, err
    }

    var ok bool
    var keydata [][]byte

    if keydata, ok = res.([][]byte); ok {
        // key data is already a double byte array
    } else {
        keydata = bytes.Fields(res.([]byte))
    }
    ret := make([]string, len(keydata))
    for i, k := range keydata {
        ret[i] = string(k)
    }
    return ret, nil
}

func (client *Client) Randomkey() (string, error) {
    res, err := client.SendCommand("RANDOMKEY")
    if err != nil {
        return "", err
    }
    return res.(string), nil
}

func (client *Client) Rename(src string, dst string) error {
    _, err := client.SendCommand("RENAME", src, dst)
    if err != nil {
        return err
    }
    return nil
}

func (client *Client) Renamenx(src string, dst string) (bool, error) {
    res, err := client.SendCommand("RENAMENX", src, dst)
    if err != nil {
        return false, err
    }
    return res.(int64) == 1, nil
}

func (client *Client) Dbsize() (int, error) {
    res, err := client.SendCommand("DBSIZE")
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Expire(key string, time int64) (bool, error) {
    res, err := client.SendCommand("EXPIRE", key, strconv.FormatInt(time, 10))

    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Ttl(key string) (int64, error) {
    res, err := client.SendCommand("TTL", key)
    if err != nil {
        return -1, err
    }

    return res.(int64), nil
}

func (client *Client) Move(key string, dbnum int) (bool, error) {
    res, err := client.SendCommand("MOVE", key, strconv.Itoa(dbnum))

    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}
