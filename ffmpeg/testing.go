package ffmpeg

import (
  "io"
  "os"
  "io/ioutil"
  "encoding/json"
)

type MockFfmpeg struct {
  Embedded string
}

func (m *MockFfmpeg) OptimizeAlbumArt(s, d string) (string, error) {
  // temp file for optimizing
  tmp, err := ioutil.TempFile("", "")
  if err != nil {
    return "", err
  }
  defer os.Remove(tmp.Name())
  defer tmp.Close()

  b, err := ioutil.ReadFile(s)
  if err != nil {
    return "", err
  }

  // can make smaller
  contents := string(b)
  if len(contents) > 0 {
    _, err = io.WriteString(tmp, contents[:len(contents)-1])
    if err != nil {
      return "", err
    }
    // use instead of original source
    s = tmp.Name()
  }

  err = copyFile(s, d)
  if err != nil {
    return "", err
  }
  return "", nil
}

func (m *MockFfmpeg) Exec(args ...string) (string, error) {
  // hook on extract audio
  if len(args) == 4 {
    err := ioutil.WriteFile(args[3], []byte(m.Embedded), 0644)
    if err != nil {
      return "", err
    }
  }
  return "", nil
}

func (m *MockFfmpeg) ToMp3(c *Mp3Config) (string, error) {
  b, err := json.Marshal(c)
  if err != nil {
    return "", err
  }

  err = ioutil.WriteFile(c.Output, b, 0644)
  if err != nil {
    return "", err
  }

  return c.Output, nil
}

func copyFile(srcPath, destPath string) (err error) {
  srcFile, err := os.Open(srcPath)
  if err != nil {
    return
  }
  defer srcFile.Close()

  destFile, err := os.Create(destPath)
  if err != nil {
    return
  }
  defer destFile.Close()

  _, err = io.Copy(destFile, srcFile)
  if err != nil {
    return
  }

  err = destFile.Sync()
  return
}
