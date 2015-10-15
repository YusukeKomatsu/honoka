package honoka

import (
  "testing"
  "fmt"
)

func TestNew(t *testing.T) {
    cli, err := New()
    if err != nil {
        t.Errorf("occurred error when get index path: %#v", err)
    }
    // fmt.Println(cli)
}

func TestGetIndexPath(t *testing.T) {
    path, err := getIndexPath()
    if err != nil {
        t.Errorf("occurred error when get index path: %v", err)
    }
    // fmt.Println(path)
}

func TestUpdateIndexFile(t *testing.T) {
    idx := []byte("{\"test\": {\"key\":\"test\",\"bucket\":\"qwertyuiopasdfghjkl567ghjk\",\"expiration\":32528049000},\"foobar\": {\"key\":\"foobar\",\"bucket\":\"qwertyuioyujhgdhjkl567ghjk\",\"expiration\": 32518049000}}")

    err := updateIndexFile(idx)
    if err != nil {
        t.Errorf("occurred error when update index file: %v", err)
    }
}