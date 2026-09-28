package slicex

// List 可链式处理的切片。
type List[T comparable] []T

// quark 空结构体占位，提高可读性
type quark struct{}

// New 包一层供 Filter 和 Unique 链式调用。不复制底层数组。
func New[T comparable](src []T) *List[T] {
	list := List[T](src)
	return &list
}

// Filter 留下谓词为真的元素，并保留原顺序。
func (l *List[T]) Filter(keep func(T) bool) *List[T] {
	out := make(List[T], 0, len(*l))
	for _, item := range *l {
		if !keep(item) {
			continue
		}
		out = append(out, item)
	}
	// 写回新切片，调用方原来的切片头不变。
	*l = out
	return l
}

// Unique 按首次出现去重，并保留原顺序。
func (l *List[T]) Unique() *List[T] {
	seen := make(map[T]quark, len(*l))
	out := make(List[T], 0, len(*l))
	for _, item := range *l {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = quark{}
		out = append(out, item)
	}
	*l = out
	return l
}

// Values 取出底层切片。没有元素时返回空切片。
func (l *List[T]) Values() []T {
	if l == nil || len(*l) == 0 {
		return []T{}
	}
	return *l
}
