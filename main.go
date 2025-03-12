package main

import "fmt"

type Getter[T1 any, T2 any] interface {
	get2() *T2
}

type Container1[T1 any, T2 any] struct {
	v1 *T1
	v2 *T2
}

type Container2[T1 any, T2 any] struct {
	v2 *T2
	v1 *T1
}

func (c Container1[T1, T2]) get2() *T2 { return c.v2 }
func (c Container2[T1, T2]) get2() *T2 { return c.v2 }

func transmute[T1 any, T2 any](v1 *T1) *T2 {
	if v1 == nil {
		return nil
	}

	c1 := Container1[T1, T2]{
		v1: v1,
		v2: nil,
	}

	c2 := Container2[T1, T2]{
		v2: nil,
		v1: v1,
	}

	var v Getter[T1, T2] = c1
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				v = c1
				v = c2
			}
		}
	}()
	for {
		if v2 := v.get2(); v2 != nil {
			done <- struct{}{}
			return v2
		}
	}
}

func ptr2addr[T any](ptr *T) int {
	return *transmute[*T, int](&ptr)
}

func addr2ptr[T any](addr int) *T {
	return *transmute[int, *T](&addr)
}

type RawString struct {
	ptr *byte
	len int
}

func string2raw(s *string) *RawString {
	return transmute[string, RawString](s)
}

func getDynamicString() string {
	return fmt.Sprintf("foobar")
}

func main() {
	s := getDynamicString()

	fmt.Printf("s: %s\n", s)

	raw := string2raw(&s)
	fmt.Printf("ptr=%p, len=%d\n", raw.ptr, raw.len)
	fmt.Println()

	fmt.Println("Shortening len:")
	raw.len = 4
	fmt.Printf("s: %s\n", s)
	fmt.Println()

	fmt.Println("Overwriting first byte:")
	*raw.ptr = byte(65)
	fmt.Printf("s: %s\n", s)
	fmt.Println()

	fmt.Println("Overwriting second by using addr2ptr and ptr2addr:")
	addr := ptr2addr(raw.ptr)
	fmt.Printf("address of ptr as int: 0x%x", addr)
	byte2 := addr2ptr[byte](addr + 1)
	*byte2 = 65
	fmt.Printf("s: %s\n", s)
	fmt.Println()

	fmt.Println("Taking substring:")
	s2 := s[2:]
	fmt.Printf("s[2:]: %s\n", s2)
	raw2 := string2raw(&s2)
	fmt.Printf("ptr=%p, len=%d\n", raw2.ptr, raw2.len)
	fmt.Println()

	fmt.Println("Overwriting third byte (first byte of substring):")
	*raw2.ptr = byte(65)

	fmt.Printf("s: %s\n", s)
	fmt.Println()

	fmt.Println("Increasing len:")
	raw.len = 10000
	fmt.Printf("s: %s\n", s)
}
