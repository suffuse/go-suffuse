#!/bin/bash

# We need the following inside User Mode Linux
CURDIR="`pwd`"
GOPATH="$GOPATH"
GOROOT="$GOROOT"
PATH="$PATH"

cat > umltest.inner.sh <<EOF
#!/bin/sh
(
   set -e
   set -x

   # Enable fuse
   insmod /usr/lib/uml/modules/\`uname -r\`/kernel/fs/fuse/fuse.ko

   # Navigate to the current directory
   cd "$CURDIR"

   # Mount the processes of the external environment
   mount -t proc proc /proc

   # Enable the network
   ifconfig lo up; ifconfig eth0 10.0.2.15; ip route add default via 10.0.2.1

   # Make sure we live in the same environment
   export GOPATH="$GOPATH"
   export GOROOT="$GOROOT"
   export PATH="$PATH"

   # List Go version and environment
   go version
   go env

   # Run tests
   go test -v ./...

   echo Success
)
echo "\$?" > "$CURDIR"/umltest.status
halt -f
EOF

chmod +x umltest.inner.sh

# Execute the script in user mode linux
/usr/bin/linux.uml init=`pwd`/umltest.inner.sh eth0=slirp rootfstype=hostfs mem=384M rw 2>&1 | \
   egrep -v 'modprobe: FATAL: Could not load /lib/modules|\[no test files\]'

exit $(<umltest.status)
