boot    /dev/hda
lba32
prompt
timeout 500
delay   500
vga     normal
root    current

images
  linux-devfs
    kernel  /pkg/kernel/bzImage

  linux-reiser
    kernel  /pkg/kernel/bz-2.2.19-devfs-reiser

  install
    kernel  /pkg/kernel/bz-2.2.19-devfs-ram
    initrd  /pkg-src/sorcerer-linux/initrd
    append  root=/dev/ram0
