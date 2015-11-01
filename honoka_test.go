package honoka

import (
  "testing"
  "path/filepath"
  "time"
  "strconv"
  homedir "github.com/mitchellh/go-homedir"
  // for Debug
  // "github.com/davecgh/go-spew/spew"
)

func TestGetIndexPath(t *testing.T) {
    actual, err := getIndexPath()
    if err != nil {
        t.Errorf("occurred error when get index path: %v", err)
    }

    home, _ := homedir.Dir()
    expected := filepath.Join(home, ".honoka", "index")
    if actual != expected {
        t.Errorf("actual does not match expected. actual: %s , expected: %s", actual, expected)
    }
}

func TestGetBucketsDirPath(t *testing.T) {
    actual, err := getBucketsDirPath()
    if err != nil {
        t.Errorf("occurred error when get bucket directory path: %v", err)
    }

    home, _ := homedir.Dir()
    expected := filepath.Join(home, ".honoka", "buckets")
    if actual != expected {
        t.Errorf("actual does not match expected. actual: %s , expected: %s", actual, expected)
    }
}

func TestGetBucketPath(t *testing.T) {
    dummyBucket := "foobar"
    actual, err := getBucketPath(dummyBucket)
    if err != nil {
        t.Errorf("occurred error when get bucket directory path: %v", err)
    }

    home, _ := homedir.Dir()
    expected := filepath.Join(home, ".honoka", "buckets", dummyBucket)
    if actual != expected {
        t.Errorf("actual does not match expected. actual: %s , expected: %s", actual, expected)
    }
}

func TestNew(t *testing.T) {
    _, err := New()
    if err != nil {
        t.Errorf("occurred error when get cache client: %#v", err)
    }
}

func TestSaveCache(t *testing.T) {
    cli, err := New()
    if err != nil {
        t.Errorf("occurred error when get cache client: %#v", err)
    }

    err = cli.Set("TestCacheString", "foobar", 100)
    if err != nil {
        t.Errorf("occurred error when set cache (string): %#v", err)
    }

    err = cli.Set("TestCacheInt", "12345656789", 100)
    if err != nil {
        t.Errorf("occurred error when set cache (slice): %#v", err)
    }

    err = cli.Set("TestCacheSlice", []string{"foo", "bar", "fizz"}, 100)
    if err != nil {
        t.Errorf("occurred error when set cache (slice): %#v", err)
    }
}

func TestExpire(t *testing.T) {
    cli, err := New()
    if err != nil {
        t.Errorf("occurred error when get cache client: %#v", err)
    }

    err = cli.Set("testCache", "foobar", 3)
    if err != nil {
        t.Errorf("occurred error when set cache (string): %#v", err)
    }

    b, err := cli.GetJson("testCache")
    actual := string(b)
    expected := "\"foobar\""
    if actual != expected {
        t.Errorf("actual does not match expected. actual: %s , expected: %s", b, expected)
    }

    time.Sleep(3 * time.Second)

    b, err = cli.GetJson("testCache")
    if err != CacheIsExpired {
        t.Errorf("cache is not expired: %#v", err)
    }
}

func TestUpdate(t *testing.T) {
    callback := func() (interface{}, error) {
        i := 1000 + 1
        val := "test.update." + strconv.Itoa(i)
        return val, nil
    }

    cli, err := New()
    if err != nil {
        t.Errorf("occurred error when get cache client: %#v", err)
    }

    var output interface{}
    actual, err := cli.Update("testUpdate", callback, 100, &output)
    if err != nil {
        t.Errorf("occurred error when update cache: %#v", err)
    }

    expected := "test.update.1001"
    if actual != expected {
        t.Errorf("actual does not match expected. actual: %v , expected: %v", actual, expected)
    }
}

// func TestUpdateIndexFile(t *testing.T) {
//     idx := []byte("{\"test\": {\"key\":\"test\",\"bucket\":\"qwertyuiopasdfghjkl567ghjk\",\"expiration\":32528049000},\"foobar\": {\"key\":\"foobar\",\"bucket\":\"qwertyuioyujhgdhjkl567ghjk\",\"expiration\": 32518049000}}")

//     err := updateIndexFile(idx)
//     if err != nil {
//         t.Errorf("occurred error when update index file: %v", err)
//     }
// }