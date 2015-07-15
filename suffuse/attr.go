package suffuse

import (
  "os"
  "time"
)

type AttrKey string
type AttrValue interface{}
type AttrMap map[AttrKey]AttrValue

type DirList    map[Name]*Inode
type Gid        uint32
type InodeNum   uint64
type LinkTarget Path
type Name       string
type PermBits   os.FileMode
type Times      [2]time.Time
type Uid        uint32
type Umask      os.FileMode

const BytesKey      = AttrKey("Bytes")
const DirListKey    = AttrKey("DirList")
const GidKey        = AttrKey("Gid")
const InodeTypeKey  = AttrKey("InodeType")
const InodeNumKey   = AttrKey("InodeNum")
const LinkTargetKey = AttrKey("LinkTarget")
const NameKey       = AttrKey("Name")
const PermBitsKey   = AttrKey("PermBits")
const TimesKey      = AttrKey("Times")
const UidKey        = AttrKey("Uid")

// man 2 chown:
//   One of the owner or group id's may be left unchanged by specifying it as -1.
// Thanks chown, we'll try to forget that uid_t and gid_t are unsigned and pass a
// bit pattern which should look like -1 after you cast it back to a signed int.
const MaxValueUint32 = uint32(1 << 32 - 1)

var AllPermBits  = PermBits(os.ModePerm)
var NoGid        = Gid(MaxValueUint32)
var NoInodeNum   = InodeNum(0)
var NoLinkTarget = LinkTarget("")
var NoName       = Name("")
var NoUid        = Uid(MaxValueUint32)
var OsGid        = Gid(os.Getgid())
var OsUid        = Uid(os.Getuid())
var OsUmask      = GetUmask()

func (x *Inode) Bytes()Bytes           { return x.AttrOr(BytesKey, NoBytes).(Bytes)                }
func (x *Inode) DirList()DirList       { return x.AttrOr(DirListKey, nil).(DirList)                }
func (x *Inode) Gid()Gid               { return x.AttrOr(GidKey, OsGid).(Gid)                      }
func (x *Inode) InodeType()InodeType   { return x.AttrOr(InodeTypeKey, InodeNone).(InodeType)      }
func (x *Inode) InodeNum()InodeNum     { return x.AttrOr(InodeNumKey, NoInodeNum).(InodeNum)       }
func (x *Inode) LinkTarget()LinkTarget { return x.AttrOr(LinkTargetKey, NoLinkTarget).(LinkTarget) }
func (x *Inode) Name()Name             { return x.AttrOr(NameKey, NoName).(Name)                   }
func (x *Inode) PermBits()PermBits     { return x.AttrOr(PermBitsKey, AllPermBits).(PermBits)      }
func (x *Inode) Times()Times           { return x.AttrOr(TimesKey, TimesNow()).(Times)             }
func (x *Inode) Uid()Uid               { return x.AttrOr(UidKey, OsUid).(Uid)                      }

func BasePermBits()PermBits {
  return PermBits(os.ModePerm & ^OsUmask)
}
func TimesNow()Times {
  now := time.Now()
  return Times { now, now }
}
func (x *Inode) FuseMode()os.FileMode {
  return x.InodeType().ToFileMode() | os.FileMode(x.PermBits())
}
func (x PermBits) String()string {
  return Sprintf("0%o", os.FileMode(x))
}
