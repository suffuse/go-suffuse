package suffuse

/*


The code below is not implemented by Inode and might be valueable research results

func (x *IdNode) Read(ctx context.Context, req *f.ReadRequest, resp *f.ReadResponse) error {
  path := x.Path
  trace("[%v] Read(%+v)", path, *req)

  for _, rule := range rules {
    bytes := rule.FileData(path)
    if bytes != nil {
      handleRead(req, resp, bytes)
      return nil
    }
  }
  return nil
}


// HandleRead handles a read request assuming that data is the entire file content.
// It adjusts the amount returned in resp according to req.Offset and req.Size.
func handleRead(req *f.ReadRequest, resp *f.ReadResponse, data []byte) {
  if req.Offset >= int64(len(data)) {
    data = nil
  } else {
    data = data[req.Offset:]
  }
  if len(data) > req.Size {
    data = data[:req.Size]
  }
  n := copy(resp.Data[:req.Size], data)
  resp.Data = resp.Data[:n]
}
*/
