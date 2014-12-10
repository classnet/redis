package redis

import (     
    "strconv"
)
// sorted set命令

func (client *Client) Zadd(key string, value []byte, score float64) (bool, error) {
    res, err := client.SendCommand("ZADD", key, strconv.FormatFloat(score, 'f', -1, 64), string(value))
    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Zrem(key string, value []byte) (bool, error) {
    res, err := client.SendCommand("ZREM", key, string(value))
    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func (client *Client) Zincrby(key string, value []byte, score float64) (float64, error) {
    res, err := client.SendCommand("ZINCRBY", key, strconv.FormatFloat(score, 'f', -1, 64), string(value))
    if err != nil {
        return 0, err
    }

    data := string(res.([]byte))
    f, _ := strconv.ParseFloat(data, 64)
    return f, nil
}

func (client *Client) Zrank(key string, value []byte) (int, error) {
    res, err := client.SendCommand("ZRANK", key, string(value))
    if err != nil {
        return 0, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Zrevrank(key string, value []byte) (int, error) {
    res, err := client.SendCommand("ZREVRANK", key, string(value))
    if err != nil {
        return 0, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Zrange(key string, start int, end int) ([][]byte, error) {
    res, err := client.SendCommand("ZRANGE", key, strconv.Itoa(start), strconv.Itoa(end))
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Zrevrange(key string, start int, end int) ([][]byte, error) {
    res, err := client.SendCommand("ZREVRANGE", key, strconv.Itoa(start), strconv.Itoa(end))
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Zrangebyscore(key string, start float64, end float64) ([][]byte, error) {
    res, err := client.SendCommand("ZRANGEBYSCORE", key, strconv.FormatFloat(start, 'f', -1, 64), strconv.FormatFloat(end, 'f', -1, 64))
    if err != nil {
        return nil, err
    }

    return res.([][]byte), nil
}

func (client *Client) Zcount(key string, min float64, max float64) (int, error) {
	res, err := client.SendCommand("ZCOUNT", key, strconv.FormatFloat(min, 'f', -1, 64), strconv.FormatFloat(max, 'f', -1, 64))
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

func (client *Client) ZcountAll(key string) (int, error) {
	res, err := client.SendCommand("ZCOUNT", key, "-INF", "+INF")
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

func (client *Client) Zcard(key string) (int, error) {
    res, err := client.SendCommand("ZCARD", key)
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Zscore(key string, member []byte) (float64, error) {
    res, err := client.SendCommand("ZSCORE", key, string(member))
    if err != nil {
        return 0, err
    }

    data := string(res.([]byte))
    f, _ := strconv.ParseFloat(data, 64)
    return f, nil
}

func (client *Client) Zremrangebyrank(key string, start int, end int) (int, error) {
    res, err := client.SendCommand("ZREMRANGEBYRANK", key, strconv.Itoa(start), strconv.Itoa(end))
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}

func (client *Client) Zremrangebyscore(key string, start float64, end float64) (int, error) {
    res, err := client.SendCommand("ZREMRANGEBYSCORE", key, strconv.FormatFloat(start, 'f', -1, 64), strconv.FormatFloat(end, 'f', -1, 64))
    if err != nil {
        return -1, err
    }

    return int(res.(int64)), nil
}

