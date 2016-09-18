package lookuptree

import (
	"errors"
)

type GeoIndexInfo struct {
	IPstart  string
	IPend    string
	Country  string
	Province string
	City     string
}

/*
 * 节点索引
 */
type index struct {
	next *node         // 下层节点
	prev *index        // 上一个next不为空的索引位置
	info *GeoIndexInfo // 底层节点所包含的IP位置信息
}

/*
 * 查找树节点
 */
type node struct {
	level int      // 节点所在层数
	prev  *node    // 本层前一个节点
	next  *node    // 本层后一个节点
	i     []*index // 节点索引
}

/*
 * IP查找树
 */
type LookUpTree struct {
	root *node
}

func newNode() *node {
	n := &node{level: 0, prev: nil, i: make([]*index, 256, 256)}
	for seq := 0; seq < 256; seq++ {
		n.i[seq] = &index{}
	}
	return n
}

func (n *node) insert(geoindex *GeoIndexInfo) error {
	// 查找ip在本层的索引
	index, err := GetIpSection(geoindex.IPstart, n.level)
	if err != nil {
		return err
	} else if index < 0 || index > 255 {
		err = errors.New("Node index error.")
		return err
	}

	if n.level == 3 {
		// 处理底层node
		n.i[index].info = geoindex
	} else {
		// 处理中间层node
		if n.i[index].next == nil {
			n.i[index].next = newNode()
			n.i[index].next.level = n.level + 1
			// 对新建node建立prev连接
			var prevnode *node
			curnode := n.i[index].next
			if n.i[index].prev != nil {
				prevnode = n.i[index].prev.next
			} else if n.prev != nil {
				if n.prev.i[255].next != nil {
					prevnode = n.prev.i[255].next
				} else {
					prevnode = n.prev.i[255].prev.next
				}
			}
			// 插入新node
			if prevnode != nil {
				curnode.prev = prevnode
				curnode.next = prevnode.next
				if curnode.next != nil {
					curnode.next.prev = curnode
				}
				prevnode.next = curnode
			}
		}
		n.i[index].next.insert(geoindex)
	}
	// 建立节点内索引prev连接
	for seq := index + 1; seq < 256; seq++ {
		n.i[seq].prev = n.i[index]
		if n.i[seq].info != nil || n.i[seq].next != nil {
			break
		}
	}
	return nil
}

func (n *node) search(ip string) (geoindex *GeoIndexInfo, err error) {
	// 查找ip在本层的索引
	index, err := GetIpSection(ip, n.level)
	if err != nil {
		return nil, err
	} else if index < 0 || index > 255 {
		err = errors.New("Node index error.")
		return nil, err
	}

	// 处理底层node
	if n.level == 3 {
		// 查找上一个有效索引
		if n.i[index].info != nil {
			return n.i[index].info, nil
		} else if n.i[index].prev != nil {
			if n.i[index].prev.info != nil {
				return n.i[index].prev.info, nil
			} else {
				return nil, errors.New("Cannot find matched info!")
			}
		} else if n.prev != nil {
			return n.prev.search(ip)
		} else {
			return nil, errors.New("Cannot find matched info!")
		}
	}
	// 处理中间层node
	// 查找上一个有效索引
	if n.i[index].next != nil {
		return n.i[index].next.search(ip)
	} else if n.i[index].prev != nil {
		if n.i[index].prev.next != nil {
			return n.i[index].prev.next.search(EnlargeIP(ip, n.level))
		} else {
			return nil, errors.New("Cannot find matched info!")
		}
	} else if n.prev != nil {
		return n.prev.search(EnlargeIP(ip, n.level))
	} else {
		return nil, errors.New("Cannot find matched info!")
	}
}

func NewLookUpTree() *LookUpTree {
	return &LookUpTree{root: newNode()}
}

func (t *LookUpTree) Insert(geoindex *GeoIndexInfo) {
	t.root.insert(geoindex)
}

func (t *LookUpTree) Search(ip string) (geoindex *GeoIndexInfo, err error) {
	// 查找匹配的起始ip
	geoindex, err = t.root.search(ip)
	if err != nil {
		return nil, err
	} else {
		ipvalue, err := IpToLong(ip)
		if err != nil {
			return nil, err
		}
		ipstartvalue, err := IpToLong(geoindex.IPstart)
		if err != nil {
			return nil, err
		}
		ipendvalue, err := IpToLong(geoindex.IPend)
		if err != nil {
			return nil, err
		}
		// 确认ip在匹配的ip段内
		if ipvalue < ipstartvalue || ipvalue > ipendvalue {
			return nil, errors.New("Cannot find matched info!")
		}
	}
	return
}
