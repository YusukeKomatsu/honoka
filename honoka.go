package honoka

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "time"

    homedir "github.com/mitchellh/go-homedir"
    "github.com/mitchellh/mapstructure"
)

type Client struct {
    Indexer IndexList
}

type IndexList map[string]Index

type Index struct {
    Key        string
    Bucket     string
    Expiration int64
}

type CleanResult struct {
    Bucket string
    Error  error
}

var (
    BucketFileNotFound = errors.New("Not found specified bucket file")
    IndexFileNotFound  = errors.New("Not found specified index file")
    CacheIsExpired     = errors.New("specified cache is expired")
)

func New() (*Client, error) {
    idx, err := getIndexList()
    if  err != nil && err != IndexFileNotFound {
        return nil, err
    }

    c := &Client{
        Indexer: idx,
    }
    return c, nil
}

func (c *Client) Get(key string, output interface{}) error {
    if c.Expire(key) {
        return CacheIsExpired
    }
    cache, err := c.GetJson(key)
    if err != nil {
        return err
    }
    var result interface{}
    err = json.Unmarshal(cache, &result)
    if err != nil {
        return err
    }
    err = mapstructure.WeakDecode(result, &output);
    return err
}

func (c *Client) GetJson(key string) ([]byte, error) {
    if c.Expire(key) {
        return nil, CacheIsExpired
    }

    idx := c.Indexer[key]
    cache, err := getCacheFromBucket(idx.Bucket)
    if err != nil {
        return nil, err
    }
    return cache, nil
}

func (c *Client) Set(key string, val interface{}, expire int64) error {
    if ! c.Expire(key) {
        return nil
    }

    exp := createExpiration(expire)
    name := getBucketName(key, exp)
    err := createNewBucket(name, val)
    if err != nil {
        return err
    }
    var idx IndexList
    idx, err = getIndexList()
    if err != nil {
        idx = c.Indexer
    }

    idx[key] = Index{
        Key:        key,
        Bucket:     name,
        Expiration: exp,
    }

    return nil
}

func (c *Client) Delete(key string) error {
    idx := c.Indexer[key]
    path, err := getBucketPath(idx.Bucket)
    if err != nil {
        return err
    }
    
    err = os.Remove(path)
    return err
}

func (c *Client) Expire(key string) bool {
    idx := c.Indexer[key]
    if idx.Expiration <= time.Now().Unix() {
        c.Delete(key)
        return true
    } else {
        return false
    }
}

func (c *Client) Outdated() ([]string, error) {
    idx, err := c.getIndexer(true)
    if err != nil {
        return nil, err
    }
    currents := make(map[string]string)
    for _, i := range idx {
        currents[i.Bucket] = ""
    }

    var list []string
    buckets, err := getBucketList()
    if err != nil {
        return nil, err
    }
    for _, bucket := range buckets {
        if _, exists := currents[bucket]; !exists {
            list = append(list, bucket)
        }
    }
    return list, nil
}

func (c *Client) Clean() ([]CleanResult, error) {
    bucketsDir, err := getBucketsDirPath()
    if err != nil {
        return nil, err
    }
    list, err := c.Outdated()
    if err != nil {
        return nil, err
    }

    var result []CleanResult
    for _, bucket := range list {
        err = os.Remove(filepath.Join(bucketsDir, bucket))
        r := CleanResult{
            Bucket: bucket,
            Error:  err,
        }
        result = append(result, r)
    }
    return result, nil
}

func (c *Client) List() ([]Index, error) {
    idx, err := c.getIndexer(true)
    if err != nil {
        return nil, err
    }
    var list []Index
    for _, i := range idx {
        list = append(list, i)
    }
    
    return list, nil  
}

func (c *Client) getIndexer(replace bool) (IndexList, error) {
    if replace || c.Indexer == nil {
        idx, err := getIndexList()
        if err != nil {
            return nil, err
        }
        c.Indexer = idx
    }
    return c.Indexer, nil
}

func (c *Client) setIndexer(indexes IndexList) error {
    idx, err := json.Marshal(indexes)
    if err != nil {
        return err
    }

    if err = updateIndexFile(idx); err != nil {
        return err
    }
    return nil
}

func getBucketsDirPath() (string, error) {
    home, err := homedir.Dir()
    if err != nil {
        return "", err
    }
    bucketsDir := filepath.Join(home, ".honoka", "buckets")
    os.MkdirAll(bucketsDir, 0700)
    return bucketsDir, err
}

func getBucketPath(bucketName string) (string, error) {
    bucketsDir, err := getBucketsDirPath()
    if err != nil {
        return "", err
    }
    return filepath.Join(bucketsDir, bucketName), nil
}

func getCacheFromBucket(bucketName string) ([]byte, error) {
    path, err := getBucketPath(bucketName)
    if err != nil {
        return nil, err
    }
    if !fileExists(path) {
        return nil, BucketFileNotFound
    }
    return ioutil.ReadFile(path);
}

func getBucketList() ([]string, error) {
    bucketsDir, err := getBucketsDirPath()
    if err != nil {
        return nil, err
    }
    files, err := ioutil.ReadDir(bucketsDir)
    var list []string
    for _, fi := range files {
        if !fi.IsDir() {
            filename := fi.Name()
            list = append(list, filename)
        }
    }
    return list, nil
}

func createNewBucket(name string, val interface{}) error {
    jval, err := json.Marshal(val)
    if err != nil {
        return err
    }
    path, err := getBucketPath(name)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(path, jval, 0644)
    return err
}

func getBucketName(key string, expiration int64) string {
    k := key + "." + strconv.FormatInt(expiration, 10)
    bytes := sha256.Sum256([]byte(k))
    return hex.EncodeToString(bytes[:])
}

func getIndexPath() (string, error) {
    home, err := homedir.Dir()
    if err != nil {
        return "", err
    }
    indexDir := filepath.Join(home, ".honoka")
    os.MkdirAll(indexDir, 0700)
    return filepath.Join(indexDir, "index"), err
}

func getIndexList() (IndexList, error) {
    b, err := getIndexFromFile()
    if err != nil {
        return nil, err
    }
    var list IndexList
    err = json.Unmarshal(b, &list)
    if  err != nil {
        return nil, err
    }
    return list, nil
}

func getIndexFromFile() ([]byte, error) {
    path, err := getIndexPath()
    if err != nil {
        return nil, err
    }
    if !fileExists(path) {
        return nil, IndexFileNotFound
    }
    return ioutil.ReadFile(path);
}

func updateIndexFile(indexes []byte) error {
    path, err := getIndexPath()
    if err != nil {
        return err
    }
    return ioutil.WriteFile(path, indexes, 0644);
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}

func createExpiration(expire int64) int64 {
    return time.Now().Unix() + expire
}
