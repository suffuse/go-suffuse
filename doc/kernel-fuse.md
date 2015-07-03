# Kernel protocol

These are [bazil.org/fuse](https://github.com/bazil/fuse) internal types used for kernel/userspace bookkeeping. You don't need to know about them unless you are hacking on the fuse library.

- RequestID: to match response to request lifetime ends with response
- NodeID: directory entry kernels knows about kernel tells when to forget
- HandleID: open file kernel tells when to destroy
