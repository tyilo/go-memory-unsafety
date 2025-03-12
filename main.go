package main

import "fmt"

type RawString struct {
	ptr *byte
	len int
}

type Getter interface {
	getString() *string
	getRawString() *RawString
}

type Container1 struct {
	s   *string
	raw *RawString
}

type Container2 struct {
	raw *RawString
	s   *string
}

func (c Container1) getString() *string       { return c.s }
func (c Container2) getString() *string       { return c.s }
func (c Container1) getRawString() *RawString { return c.raw }
func (c Container2) getRawString() *RawString { return c.raw }

func string2raw(s *string) *RawString {
	raw := RawString{
		ptr: nil,
		len: -1,
	}
	c1 := Container1{
		s:   s,
		raw: &raw,
	}

	c2 := Container2{
		raw: &raw,
		s:   s,
	}

	var v Getter = c1
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
		if raw := v.getRawString(); raw.ptr != nil {
			done <- struct{}{}
			return raw
		}
	}
}

func getDynamicString() string {
	return fmt.Sprintf("foobar")
}

func main() {
	s := getDynamicString()

	fmt.Printf("s: %s\n", s)

	raw := string2raw(&s)
	fmt.Printf("ptr=%p, len=%d\n", raw.ptr, raw.len)

	fmt.Println("Shortening len:")
	raw.len = 4
	fmt.Printf("s: %s\n", s)

	fmt.Println("Overwriting first byte:")
	*raw.ptr = byte(65)
	fmt.Printf("s: %s\n", s)

	fmt.Println("Taking substring:")
	s2 := s[2:]
	fmt.Printf("s[2:]: %s\n", s2)
	raw2 := string2raw(&s2)
	fmt.Printf("ptr=%p, len=%d\n", raw2.ptr, raw2.len)

	fmt.Println("Overwriting third byte (first byte of substring):")
	*raw2.ptr = byte(65)

	fmt.Printf("s: %s\n", s)

	fmt.Println("Increasing len:")
	raw.len = 10000
	fmt.Printf("s: %s\n", s)
}
