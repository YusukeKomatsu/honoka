package honoka

import (
  "testing"
  "fmt"
)

func TestGetIndexPath(t *testing.T) {
    path, err := getIndexPath()
    if err != nil {
        t.Errorf("occurred error when get index path: %v", err)
    }
    fmt.Println(path)
}

func TestUpdateIndexFile(t *testing.T) {
    idx := "{'test': {'foo': 'a', 'hoge': 'fuga'}}"

    err := updateIndexFile(idx)
    if err != nil {
        t.Errorf("occurred error when update index file: %v", err)
    }
}