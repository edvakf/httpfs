package httpfs

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var port string = ":9999"

func TestMain(m *testing.M) {
	os.RemoveAll("/tmp/httpfs_test")

	go startServer()
	time.Sleep(100 * time.Millisecond)

	code := m.Run()
	defer os.Exit(code)

	os.RemoveAll("/tmp/httpfs_test")
}

func startServer() {
	http.HandleFunc("/fs/", Handle)
	log.Fatal(http.ListenAndServe(port, nil))
}

func TestPost(t *testing.T) {
	defer os.RemoveAll("/tmp/httpfs_test/foo")

	url := "http://127.0.0.1" + port + "/fs/httpfs_test/foo/bar"
	t.Logf("post %s", url)
	resp, err := http.Post(url, "", strings.NewReader("abc"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("response: %s", resp.Status)
	if resp.StatusCode != 200 {
		t.Fatal("bad response")
	}

	b, err := ioutil.ReadFile("/tmp/httpfs_test/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "abc" {
		t.Error("file content does not match the posted content")
	}
}

func TestPut(t *testing.T) {
	defer os.RemoveAll("/tmp/httpfs_test/foo")

	url := "http://127.0.0.1" + port + "/fs/httpfs_test/foo/bar"
	t.Logf("put %s", url)
	req, err := http.NewRequest("PUT", url, strings.NewReader("abc"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("response: %s", resp.Status)
	if resp.StatusCode != 200 {
		t.Fatal("bad response")
	}

	b, err := ioutil.ReadFile("/tmp/httpfs_test/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "abc" {
		t.Error("file content does not match the put content")
	}
}

func TestGet(t *testing.T) {
	os.RemoveAll("/tmp/httpfs_test/foo")
	err := os.MkdirAll("/tmp/httpfs_test/foo", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("/tmp/httpfs_test/foo/bar", []byte("abc"), 0644)

	url := "http://127.0.0.1" + port + "/fs/httpfs_test/foo/bar"
	t.Logf("get %s", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("response: %s", resp.Status)
	if resp.StatusCode != 200 {
		t.Fatal("bad response")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "abc" {
		t.Error("get content does not match the file content")
	}
}

func TestDeleteFile(t *testing.T) {
	os.RemoveAll("/tmp/httpfs_test/foo")
	err := os.MkdirAll("/tmp/httpfs_test/foo", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("/tmp/httpfs_test/foo/bar", []byte("abc"), 0644)

	url := "http://127.0.0.1" + port + "/fs/httpfs_test/foo/bar"
	t.Logf("delete %s", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("response: %s", resp.Status)
	if resp.StatusCode != 200 {
		t.Fatal("bad response")
	}

	_, err = os.Stat("/tmp/httpfs_test/foo/bar")
	if err != nil {
		t.Log("file deleted")
	} else {
		t.Fatal("file still exists!")
	}
}

func TestDeleteDir(t *testing.T) {
	os.RemoveAll("/tmp/httpfs_test/foo")
	err := os.MkdirAll("/tmp/httpfs_test/foo", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("/tmp/httpfs_test/foo/bar", []byte("abc"), 0644)

	url := "http://127.0.0.1" + port + "/fs/httpfs_test/foo"
	t.Logf("delete %s", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("response: %s", resp.Status)
	if resp.StatusCode != 200 {
		t.Fatal("bad response")
	}

	_, err = os.Stat("/tmp/httpfs_test/foo")
	if err != nil {
		t.Log("directory deleted")
	} else {
		t.Fatal("directory still exists!")
	}
}
