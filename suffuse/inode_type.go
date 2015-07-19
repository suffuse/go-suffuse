package suffuse

import (
  "os"
  "bazil.org/fuse"
)

/** The types of inodes, in a proper enum which we
 *  translate for the various consumers of such information
 *  such as C, go, and FUSE.
 */

type InodeType uint8

const (
  InodeNone InodeType = iota
  InodeSocket
  InodeLink
  InodeFile
  InodeBlock
  InodeDir
  InodeChar
  InodeFIFO
)

func (x *Inode) IsLink()bool   { return x.InodeType() == InodeLink }
func (x *Inode) IsDir()bool    { return x.InodeType() == InodeDir  }
func (x *Inode) IsAbsent()bool { return x.InodeType() == InodeNone }

func (x InodeType) String()string {
  switch x {
    case InodeSocket : return "socket"
    case InodeLink   : return "link"
    case InodeFile   : return "file"
    case InodeBlock  : return "blockdev"
    case InodeDir    : return "dir"
    case InodeChar   : return "chardev"
    case InodeFIFO   : return "fifo"
    default          : return "-"
  }
}
func (x InodeType) ToFuseType()fuse.DirentType {
  switch x {
    case InodeSocket : return fuse.DT_Socket
    case InodeLink   : return fuse.DT_Link
    case InodeFile   : return fuse.DT_File
    case InodeBlock  : return fuse.DT_Block
    case InodeDir    : return fuse.DT_Dir
    case InodeChar   : return fuse.DT_Char
    case InodeFIFO   : return fuse.DT_FIFO
    default   : return fuse.DT_Unknown
  }
}
func (x InodeType) ToFileMode() os.FileMode {
  switch x {
    case InodeSocket : return os.ModeSocket
    case InodeLink   : return os.ModeSymlink
    case InodeBlock  : return os.ModeDevice
    case InodeDir    : return os.ModeDir
    case InodeChar   : return os.ModeCharDevice | os.ModeDevice
    case InodeFIFO   : return os.ModeNamedPipe
    case InodeFile   : return os.FileMode(0)
    default          : return os.FileMode(0)
  }
}
func FuseTypeToInodeType(tp fuse.DirentType) InodeType {
  switch tp {
    case fuse.DT_Socket : return InodeSocket
    case fuse.DT_Link   : return InodeLink
    case fuse.DT_File   : return InodeFile
    case fuse.DT_Block  : return InodeBlock
    case fuse.DT_Dir    : return InodeDir
    case fuse.DT_Char   : return InodeChar
    case fuse.DT_FIFO   : return InodeFIFO
    default             : return InodeNone
  }
}

func FileModeToInodeType(tp os.FileMode) InodeType {
  return FuseTypeToInodeType(GoModeToDirentType(tp))
}
