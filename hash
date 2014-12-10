package redis
import (
    "reflect"
    "strconv"
    "errors"
)
// hash命令
func (client *Client) Hset(key string, field string, val []byte) (bool, error) {
    res, err := client.SendCommand("HSET", key, field, string(val))
    if err != nil {
        return false, err
    }
    return res.(int64) == 1, nil
}
func (client *Client) Hget(key string, field string) ([]byte, error) {
    res, _ := client.SendCommand("HGET", key, field)
    if res == nil {
        return nil, RedisError("Hget failed")
    }
    data := res.([]byte)
    return data, nil
}
//从这里拷贝相当多的JSON代码
func valueToString(v reflect.Value) (string, error) {
    if !v.IsValid() {
        return "null", nil
    }
    switch v.Kind() {
    case reflect.Ptr:
        return valueToString(reflect.Indirect(v))
    case reflect.Interface:
        return valueToString(v.Elem())
    case reflect.Bool:
        x := v.Bool()
        if x {
            return "true", nil
        } else {
            return "false", nil
        }
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return strconv.FormatInt(v.Int(), 10), nil
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        return strconv.FormatUint(v.Uint(), 10), nil
    case reflect.UnsafePointer:
        return strconv.FormatUint(uint64(v.Pointer()), 10), nil
    case reflect.Float32, reflect.Float64:
        return strconv.FormatFloat(v.Float(), 'g', -1, 64), nil
    case reflect.String:
        return v.String(), nil
    //This is kind of a rough hack to replace the old []byte
    //detection with reflect.Uint8Type, it doesn't catch
    //zero-length byte slices
    case reflect.Slice:
        typ := v.Type()
        if typ.Elem().Kind() == reflect.Uint || typ.Elem().Kind() == reflect.Uint8 || typ.Elem().Kind() == reflect.Uint16 || typ.Elem().Kind() == reflect.Uint32 || typ.Elem().Kind() == reflect.Uint64 || typ.Elem().Kind() == reflect.Uintptr {
            if v.Len() > 0 {
                if v.Index(0).OverflowUint(257) {
                    return string(v.Interface().([]byte)), nil
                }
            }
        }
    }
    return "", errors.New("Unsupported type")
}
func containerToString(val reflect.Value, args *[]string) error {
    switch v := val; v.Kind() {
    case reflect.Ptr:
        return containerToString(reflect.Indirect(v), args)
    case reflect.Interface:
        return containerToString(v.Elem(), args)
    case reflect.Map:
        if v.Type().Key().Kind() != reflect.String {
            return errors.New("Unsupported type - map key must be a string")
        }
        for _, k := range v.MapKeys() {
            *args = append(*args, k.String())
            s, err := valueToString(v.MapIndex(k))
            if err != nil {
                return err
            }
            *args = append(*args, s)
        }
    case reflect.Struct:
        st := v.Type()
        for i := 0; i < st.NumField(); i++ {
            ft := st.FieldByIndex([]int{i})
            *args = append(*args, ft.Name)
            s, err := valueToString(v.FieldByIndex([]int{i}))
            if err != nil {
                return err
            }
            *args = append(*args, s)
        }
    }
    return nil
}
func (client *Client) Hmset(key string, mapping interface{}) error {
    var args []string
    args = append(args, key)
    err := containerToString(reflect.ValueOf(mapping), &args)
    if err != nil {
        return err
    }
    _, err = client.SendCommand("HMSET", args...)
    if err != nil {
        return err
    }
    return nil
}
func (client *Client) Hmget(key string, fields ...string) ([][]byte, error) {
    var args []string
    args = append(args, key)
    for _, field := range fields {
        args = append(args, field)
    }
    res, err := client.SendCommand("HMGET", args...)
    if err != nil {
        return nil, err
    }
    return res.([][]byte), nil
}
func (client *Client) Hincrby(key string, field string, val int64) (int64, error) {
    res, err := client.SendCommand("HINCRBY", key, field, strconv.FormatInt(val, 10))
    if err != nil {
        return -1, err
    }
    return res.(int64), nil
}
func (client *Client) Hexists(key string, field string) (bool, error) {
    res, err := client.SendCommand("HEXISTS", key, field)
    if err != nil {
        return false, err
    }
    return res.(int64) == 1, nil
}
func (client *Client) Hdel(key string, field string) (bool, error) {
    res, err := client.SendCommand("HDEL", key, field)
    if err != nil {
        return false, err
    }
    return res.(int64) == 1, nil
}
func (client *Client) Hlen(key string) (int, error) {
    res, err := client.SendCommand("HLEN", key)
    if err != nil {
        return -1, err
    }
    return int(res.(int64)), nil
}

func (client *Client) Hkeys(key string) ([]string, error) {
    res, err := client.SendCommand("HKEYS", key)
    if err != nil {
        return nil, err
    }
    data := res.([][]byte)
    ret := make([]string, len(data))
    for i, k := range data {
        ret[i] = string(k)
    }
    return ret, nil
}
func (client *Client) Hvals(key string) ([][]byte, error) {
    res, err := client.SendCommand("HVALS", key)
    if err != nil {
        return nil, err
    }
    return res.([][]byte), nil
}
func writeTo(data []byte, val reflect.Value) error {
    s := string(data)
    switch v := val; v.Kind() {
    // if we're writing to an interace value, just set the byte data
    // TODO: should we support writing to a pointer?
    case reflect.Interface:
        v.Set(reflect.ValueOf(data))
    case reflect.Bool:
        b, err := strconv.ParseBool(s)
        if err != nil {
            return err
        }
        v.SetBool(b)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        i, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return err
        }
        v.SetInt(i)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        ui, err := strconv.ParseUint(s, 10, 64)
        if err != nil {
            return err
        }
        v.SetUint(ui)
    case reflect.Float32, reflect.Float64:
        f, err := strconv.ParseFloat(s, 64)
        if err != nil {
            return err
        }
        v.SetFloat(f)
    case reflect.String:
        v.SetString(s)
    case reflect.Slice:
        typ := v.Type()
        if typ.Elem().Kind() == reflect.Uint || typ.Elem().Kind() == reflect.Uint8 || typ.Elem().Kind() == reflect.Uint16 || typ.Elem().Kind() == reflect.Uint32 || typ.Elem().Kind() == reflect.Uint64 || typ.Elem().Kind() == reflect.Uintptr {
            v.Set(reflect.ValueOf(data))
        }
    }
    return nil
}
func writeToContainer(data [][]byte, val reflect.Value) error {
    switch v := val; v.Kind() {
    case reflect.Ptr:
        return writeToContainer(data, reflect.Indirect(v))
    case reflect.Interface:
        return writeToContainer(data, v.Elem())
    case reflect.Map:
        if v.Type().Key().Kind() != reflect.String {
            return errors.New("Invalid map type")
        }
        elemtype := v.Type().Elem()
        for i := 0; i < len(data)/2; i++ {
            mk := reflect.ValueOf(string(data[i*2]))
            mv := reflect.New(elemtype).Elem()
            writeTo(data[i*2+1], mv)
            v.SetMapIndex(mk, mv)
        }
    case reflect.Struct:
        for i := 0; i < len(data)/2; i++ {
            name := string(data[i*2])
            field := v.FieldByName(name)
            if !field.IsValid() {
                continue
            }
            writeTo(data[i*2+1], field)
        }
    default:
        return errors.New("Invalid container type")
    }
    return nil
}
func (client *Client) Hgetall(key string, val interface{}) error {
    res, err := client.SendCommand("HGETALL", key)
    if err != nil {
        return err
    }
    data := res.([][]byte)
    if data == nil || len(data) == 0 {
        return RedisError("Key `" + key + "` does not exist")
    }
    err = writeToContainer(data, reflect.ValueOf(val))
    if err != nil {
        return err
    }
    return nil
}

