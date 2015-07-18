package xattr

/** Packages confined to this file: bazil.org/fuse/syscallx
 */

import (
  "bytes"
  "howett.net/plist"
  "bazil.org/fuse"
  "bazil.org/fuse/syscallx"
  // "syscall"
)

// TODO: Other ideas
// "com.apple.metadata:kMDItemWhereFroms"
// "com.apple.diskimages.fsck"
// "com.apple.quarantine"
// "com.apple.diskimages.recentcksum"
// "com.dropbox.attributes"
// "com.apple.TextEncoding"
// "com.apple.metadata:com_apple_backup_excludeItem"
// "com.apple.metadata:kMDItemDownloadedDate"
const (
  ListBufsize  = 2048
  GetBufsize   = 4096
  OsxTags      = "com.apple.metadata:_kMDItemUserTags"
  FinderInfo   = "com.apple.FinderInfo"
  ResourceFork = "com.apple.ResourceFork"
)

func FuseXattrByteBuffer(x interface{}, altSize int) []byte {
  size := altSize
  switch x := x.(type) {
    case *fuse.GetxattrRequest  : if x.Size > 0 { size = int(x.Size) }
    case *fuse.ListxattrRequest : if x.Size > 0 { size = int(x.Size) }
  }
  return make([]byte, size)
}

func Get(path string, name string) (bytes []byte) {
  buf := make([]byte, GetBufsize)
  read, err := syscallx.Getxattr(path, name, buf)
  if err == nil && read > 0 {
    bytes = buf[:read]
  }
  return
}
func List(path string) (names []string) {
  buf := make([]byte, ListBufsize)
  read, err := syscallx.Listxattr(path, buf)
  if err == nil && read > 0 {
    names = ListToStrings(buf[:read])
  }
  return
}
func Set(path string, key string, value []byte) error {
  return syscallx.Setxattr(path, key, value, 0)
}
func Remove(path string, key string) error {
  return syscallx.Removexattr(path, key)
}

// The extended attribute names are simple NULL-terminated
// UTF-8 strings and are returned in arbitrary order.
func ListToStrings(buf []byte) []string {
  buflen := len(buf)
  var res []string
  idx := 0

  for idx < buflen {
    start := idx
    for buf[idx] != 0 && idx < buflen { idx++ }
    if start < idx && buf[idx] == 0 {
      res = append(res, string(buf[start:idx]))
    }
    idx++
  }

  return res
}

// TODO: plist encoding/decoding.
type PlistWriter struct {
  Bytes []byte
}
func (x PlistWriter) Write(p []byte) (n int, err error) {
  x.Bytes = append(x.Bytes, p...)
  return len(p), nil
}

func AnyToBinaryPlist(x interface{}) (result []byte) {
  w := PlistWriter { }
  encoder := plist.NewBinaryEncoder(w)
  encoder.Encode(x)
  return w.Bytes
}
func BinaryPlistToMap(bs []byte) (result map[string]string) {
  buf := bytes.NewReader(bs)
  decoder := plist.NewDecoder(buf)
  decoder.Decode(&result)
  return
}
