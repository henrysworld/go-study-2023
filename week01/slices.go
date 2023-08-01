package week01

import "reflect"

func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j]
	return append(s[:i], s[j:]...)
}

func Shrink[S ~[]E, E any](s S) S {
	v := reflect.ValueOf(s)
	n := v.Len()

	// 创建一个新的与原来类型相同的切片，容量和长度为原来的一半
	ns := reflect.MakeSlice(v.Type(), n/2, n/2)
	reflect.Copy(ns, v.Slice(0, n/2))

	// 将新切片转回原来的类型
	return ns.Interface().(S)
}

func DeleteAndShrink[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check

	v := reflect.ValueOf(s)
	// 删除元素后的长度
	newLen := len(s) - (j - i)
	// 判断是否需要缩容
	var ns reflect.Value
	if cap(s) > newLen*4 {
		if newLen == 0 {
			ns = reflect.MakeSlice(v.Type(), newLen, 4)
		} else {
			// 创建一个新的与原来类型相同的切片，容量和长度为删除元素后长度的2倍
			ns = reflect.MakeSlice(v.Type(), newLen, newLen*2)
		}

	} else {
		ns = reflect.MakeSlice(v.Type(), newLen, cap(s))
	}

	// 将 s[:i] 和 s[j:] 的元素复制到新切片中
	reflect.Copy(ns.Slice(0, i), v.Slice(0, i))
	reflect.Copy(ns.Slice(i, ns.Len()), v.Slice(j, v.Len()))

	// 将新切片转回原来的类型
	return ns.Interface().(S)
}

func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	return append(S([]E{}), s...)
}

func Equal[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
