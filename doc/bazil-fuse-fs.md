# [bazil.org/fuse/fs](https://github.com/bazil/fuse/blob/master/fs/serve.go) data types

fs/serve.go is the provided "high level" FUSE interface, as opposed to the "low level"
interface where you work directly with request and response structure pointers. In the
high level interface you get to exploit some common infrastructure code like keeping
track of open filehandles.

These are the interface names. bazil/fs places one method in each, then performs
type assertions to see if you've implemented that one method.

Node is an inode, Handle is a file descriptor.

```go
# An FS is the interface required of a file system.
FS
FSDestroyer
FSIniter
FSInodeGenerator
FSStatfser

# A Handle is the interface required of an opened file or directory.
Handle
HandleFlusher
HandleReadAller
HandleReadDirAller
HandleReader
HandleReleaser
HandleWriter

# A Node is the interface required of a file or directory.
Node
NodeAccesser
NodeCreater
NodeForgetter
NodeFsyncer
NodeGetattrer
NodeGetxattrer
NodeLinker
NodeListxattrer
NodeMkdirer
NodeMknoder
NodeOpener
NodeReadlinker
NodeRemover
NodeRemovexattrer
NodeRenamer
NodeRequestLookuper
NodeSetattrer
NodeSetxattrer
NodeStringLookuper
NodeSymlinker
```

These are the methods found in the interfaces above. You implement these as "instance" methods on some struct type which conforms to the FS, Handle, or Node interface.
```go
# Required - only one FS method and one Node method are mandatory.
func (FS) Root() (Node, error)
func (Node) Attr(ctx context.Context, attr *fuse.Attr) error

# Returning (Node, err)
Link       (ctx, req *LinkRequest, old Node) Node
Lookup     (ctx, name string) Node
Lookup     (ctx, req *LookupRequest, resp *LookupResponse) Node
Mkdir      (ctx, req *MkdirRequest) Node
Mknod      (ctx, req *MknodRequest) Node
Symlink    (ctx, req *SymlinkRequest) Node

# Returning (..., err)
Create     (ctx, req *CreateRequest, resp *CreateResponse) (Node, Handle)
Open       (ctx, req *OpenRequest, resp *OpenResponse) Handle
ReadAll    (ctx) []byte
ReadDirAll (ctx) []Dirent
Readlink   (ctx, req *ReadlinkRequest) string

# Returning err
Access     (ctx, req *AccessRequest)
Attr       (ctx, attr *Attr)
Flush      (ctx, req *FlushRequest)
Fsync      (ctx, req *FsyncRequest)
Getattr    (ctx, req *GetattrRequest, resp *GetattrResponse)
Getxattr   (ctx, req *GetxattrRequest, resp *GetxattrResponse)
Init       (ctx, req *InitRequest, resp *InitResponse)
Listxattr  (ctx, req *ListxattrRequest, resp *ListxattrResponse)
Read       (ctx, req *ReadRequest, resp *ReadResponse)
Release    (ctx, req *ReleaseRequest)
Remove     (ctx, req *RemoveRequest)
Removexattr(ctx, req *RemovexattrRequest)
Rename     (ctx, req *RenameRequest, newDir Node)
Setattr    (ctx, req *SetattrRequest, resp *SetattrResponse)
Setxattr   (ctx, req *SetxattrRequest)
Statfs     (ctx, req *StatfsRequest, resp *StatfsResponse)
Write      (ctx, req *WriteRequest, resp *WriteResponse)
```
