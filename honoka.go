package honoka

import (
    // "crypto/sha256"
    // "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
    // "strconv"
    // "time"

    homedir "github.com/mitchellh/go-homedir"
    // "github.com/mitchellh/mapstructure"
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

var (
    FileNotFound = errors.New("Not found the specified file")
)

func New() (*Client, error) {
    idx, err := getIndexList()
    if err != nil {
        return nil, err
    }

    c := &Client{
        Indexer: idx,
    }
    return c, nil
}

func (c *Client) Expire(key string) (bool, error) {
    idx := c.Indexer[key]
    if idx.Expiration <= time.Now().Unix() {
        c.Delete(key)
        return true
    } else {
        return false
    }
}


// func (c *Client) setIndexer(indexes IndexList) error {
//     idxj, err := json.Marshal(indexes)
//     if err != nil {
//         return nil, err
//     }

//     path, err := getIndexPath()
//     if err != nil {
//         return nil, err
//     }

//     if err = ioutil.WriteFile(path, []byte(idxj), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
//         return nil, err
//     }

//     return nil
// }

// func (c *Client) getIndexer(replace bool) (*Client.Indexer, error) {
//     if replace || c.Indexer == nil {
//         path, err := getIndexPath()
//         if err != nil {
//             return
//         }

//         idx, err := getIndexList()
//         if err != nil {
//             return nil, err
//         }
//         err = c.setIndexer(idx)
//         if err != nil {
//             return nil, err
//         }
//     }
//     return c.Indexer, nil
// }

func getIndexPath() (string, error) {
    home, err := homedir.Dir()
    if err != nil {
        return "", err
    }
    indexDir := filepath.Join(home, ".honoka")
    os.MkdirAll(indexDir, 0700)
    return filepath.Join(indexDir, "index"), err
}

func getIndexFromFile() ([]byte, error) {
    _, err := getIndexPath()
    if err != nil {
        return nil, err
    }
    return ioutil.ReadFile(path);
}

func updateIndexFile(indexes []byte) error {
    path, err := getIndexPath()
    if err != nil {
        return err
    }
    return ioutil.ReadFile(path, indexes, 0644);
}
