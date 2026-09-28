package slicex

// Func 定义单个对象的转换函数。
// 适用于 repository 层显式声明 model <-> domain 转换逻辑。
type Func[Src any, Dst any] func(Src) Dst

// Mapper 批量转换对象列表。
func Mapper[Src any, Dst any](src []Src, fn Func[Src, Dst]) []Dst {
	dst := make([]Dst, len(src))
	for i, item := range src {
		dst[i] = fn(item)
	}
	return dst
}
