package fastlocalcache

type Ring struct {
	head uint64
	tail uint64
	data []byte
}

// 可能会因空间不够报错
func (r *Ring) Push(ety entry) error {
	// 在 tail 的位置追加元素
	panic("")
}

func (r *Ring) SetState(keyIndex uint32, state entryState) error {
	// 设置index处的entry的state
	panic("")
}

func (r *Ring) Pop() (entry, error) {
	//
	panic("")
}

func (r *Ring) Peek() (entry, error) {
	// 从head处，拿出来元素
	panic("")
}

// 取出来keyIndex处的entry
func (r *Ring) Get(keyIndex uint32) (entry, error) {
	panic("")
}
