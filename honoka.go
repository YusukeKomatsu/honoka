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

    // for Debug
    // "fmt"
    // "github.com/davecgh/go-spew/spew"
)

// Cache client
type Client struct {
    // Cache Index list
    Indexer IndexList
}

// Cache index list
type IndexList map[string]Index

// Cache index
type Index struct {
    // The index key.
    Key        string

    // The bucket name that saved cache data.
    Bucket     string

    // The maximum elapsed time since the last file update.
    Expiration int64
}

// The structure is used when use clean method.
type CleanResult struct {
    // The bucket name that saved cache data.
    Bucket string

    // Error when delete the specified bucket.
    Error  error
}

type UpdateFunc func() (interface{}, error)

var (
    Version = "0.0.1"
    BucketFileNotFound = errors.New("Not found specified bucket file")
    IndexFileNotFound  = errors.New("Not found specified index file")
    CacheIsExpired     = errors.New("specified cache is expired")
)

// New is a function for making a new cache
func New() (*Client, error) {
    idx, err := getIndexList()
    if err != nil {
        if err == IndexFileNotFound {
            idx = nil
        } else {
            return nil, err
        }
    }

    c := &Client{
        Indexer: idx,
    }
    return c, nil
}

// Get is used to retrieve a cache by specified key.
// 
// Example:
//   cli, err := honoka.New()
//   var output interface{}
//   cli.Get("foobar", &output)
//   // OR
//   result, err := cli.Get("foobar", &output)
func (c *Client) Get(key string, output interface{}) (interface{}, error) {
    if c.Expire(key) {
        return nil, CacheIsExpired
    }
    cache, err := c.GetJson(key)
    if err != nil {
        return nil, err
    }
    var result interface{}
    err = json.Unmarshal(cache, &result)
    if err != nil {
        return nil, err
    }
    err = mapstructure.WeakDecode(result, &output);
    return &output, err
}

// Get is used to retrieve a cache by specified key.
// Return value is JSON string
// Example:
//   cli, err := honoka.New()
//   result, err := cli.GetJson("foobar")
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

// Get is used to create a cache if specified key has not used yet.
// 
// Example:
//   cli, err := honoka.New()
//   err := cli.Set("foobar", "fizzbizz", 100)
func (c *Client) Set(key string, val interface{}, expire int64) error {
    if ! c.Expire(key) {
        return nil
    }

    exp := createExpiration(expire)
    name := getBucketName(key, exp)
    _, err := createNewBucket(name, val)
    if err != nil {
        return err
    }
    var idx IndexList
    idx, err = getIndexList()
    if err != nil {
        if (err == IndexFileNotFound) {
            idx = map[string]Index{}
        } else {
            return err
        }
    }

    idx[key] = Index{
        Key:        key,
        Bucket:     name,
        Expiration: exp,
    }
    c.setIndexer(idx)

    return nil
}

// Update calls the cache update function on the cached data.
// Get is used to retrieve a cache by specified key.
// 
// Example:
//   cli, err := honoka.New()
//   var output interface{}
//   f := func() { return "fizzbizz" }
//   cli.Update("foobar", f, 100, &output)
//   // OR
//   result, err := cli.Get("foobar", f, 100, &output)
func (c *Client) Update(key string, updater UpdateFunc, expire int64, output interface{}) (interface{}, error) {
    b, err := c.UpdateJson(key, updater, expire)
    if b != nil {
        var result interface{}
        e := json.Unmarshal(b, &result)
        if e != nil {
            return nil, e
        }

        e = mapstructure.WeakDecode(result, &output);
        if e != nil {
            return nil, e
        }
    }

    return output, err
}

// Update calls the cache update function on the cached data.
// Return value is JSON string.
// 
// Example:
//   cli, err := honoka.New()
//   f := func() { return "fizzbizz" }
//   result, err := cli.UpdateJson("foobar", f, 100)
func (c *Client) UpdateJson(key string, updater UpdateFunc, expire int64) ([]byte, error) {
    if ! c.Expire(key) {
        return c.GetJson(key)
    }

    val, err := updater()
    if err != nil {
        return nil, err
    }

    exp := createExpiration(expire)
    name := getBucketName(key, exp)
    jval, err := createNewBucket(name, val)
    if err != nil {
        return jval, err
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
    c.setIndexer(idx)

    return jval, nil
}

// Delete is used to delete a cache by specified key.
// 
// Example:
//   cli, err := honoka.New()
//   err = cli.Delete("foobar")
func (c *Client) Delete(key string) error {
    idx := c.Indexer[key]
    path, err := getBucketPath(idx.Bucket)
    if err != nil {
        return err
    }
    if fileExists(path) {
        err = os.Remove(path)
        if err != nil {
            return err
        }
    }

    delete(c.Indexer, key)
    c.setIndexer(c.Indexer)
    return nil
}

// Expire is a predicate which determines if the cache should be updated.
// 
// Example:
//   cli, err := honoka.New()
//   expired := cli.Expire("foobar")
func (c *Client) Expire(key string) bool {
    if nil == c.Indexer {
        return true
    }

    idx, exists := c.Indexer[key]
    if exists {
        if idx.Expiration <= time.Now().Unix() {
            c.Delete(key)
            return true
        } else {
            return false
        }
    }
    return true
}

// Outdated is used to retrive no-indexed bucket.
// 
// Example:
//   cli, err := honoka.New()
//   list, err := cli.Outdated()
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

// Clean is used to delete no-indexed bucket.
// Example:
//   cli, err := honoka.New()
//   result, err := cli.Clean()
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
        e := os.Remove(filepath.Join(bucketsDir, bucket))
        r := CleanResult{
            Bucket: bucket,
            Error:  e,
        }
        result = append(result, r)
    }
    return result, nil
}

// List is used to retrive cache indexes.
// 
// Example:
//   cli, err := honoka.New()
//   list, err := cli.List()
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
    c.Indexer = indexes
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

func createNewBucket(name string, val interface{}) ([]byte, error) {
    jval, err := json.Marshal(val)
    if err != nil {
        return nil, err
    }
    path, err := getBucketPath(name)
    if err != nil {
        return jval, err
    }
    err = ioutil.WriteFile(path, jval, 0644)
    return jval, err
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
