# two ways of writing a block of text

description
   "Order is significant, and duplicate nodes
    are allowed. Per definition, each node is
    accesible thru a path. To distinguish between
    duplicated nodes, indexes can be used."

description \
    Order is significant, and duplicate nodes
    are allowed. Per definition, each node is
    accesible thru a path. To distinguish between
    duplicated nodes, indexes can be used.

"^(.*)#([^#]+)$"
  "$1 $2"
