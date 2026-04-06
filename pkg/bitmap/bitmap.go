package bitmap

type Bitmap struct {
	bits []byte // 注意此处为byte数组
	size int
}

func NewBitmap(size int) *Bitmap {
	if size == 0 {
		size = 250 // 默认大小为250
	}
	//  [0,0,0,0][0,0,0,0] 每个byte中有8个bit
	return &Bitmap{
		// 指定有size个byte
		bits: make([]byte, size),
		// 总共有多少个bit
		size: size * 8,
	}
}

func (b *Bitmap) Set(id string) {
	// 先计算id在哪个bit
	idx := hash(id) % b.size
	// 再根据bit在哪个位置去计算在哪个字节
	byteIdx := idx / 8
	// 在这个字节中的哪个bit位置
	bitIdx := idx % 8

	// 将00000001向左移动bitIdx位：结果是一个掩码，只有第 bitIdx 位是1，其余位都是0
	// 再位或运算：通过位或的方式来设置为1
	b.bits[byteIdx] |= 1 << bitIdx
}

// IsSet 检查特定位是否为1
func (b *Bitmap) IsSet(id string) bool {
	idx := hash(id) % b.size
	byteIdx := idx / 8
	bitIdx := idx % 8
	// 将00000001左移bitIdx位后
	// 再与原二进制数值进行位与运算判断特定位是否已为1
	return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}

// 导出
func (b *Bitmap) Export() []byte {
	return b.bits // 输出当前的字节数组
}

// 导入（加载）
func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}

	return &Bitmap{
		bits: bits,
		size: len(bits) * 8, // 注意将byte转换为bit长度
	}
}

func hash(id string) int {
	// BKDR hash
	seed := 131313
	hash := 0
	for _, c := range id {
		hash = hash*seed + int(c)
	}
	return hash & 0x7FFFFFFF
}
