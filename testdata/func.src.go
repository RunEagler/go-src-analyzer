package testdata

import (
	"time"
	"fmt"
)

type sampleStruct struct {
	i int
	s string
}

type constType int

const (
	A constType = iota
	B
	C
)

func test1() {

}

func test2(arg1 int) (err error) {

	return fmt.Errorf("")
}

func test3(
	arg1 int,
	arg2 string,
	arg3 []int,
	arg4 sampleStruct,
	arg5 *int,
	arg6 []*string,
	arg7 chan string,
	arg8 constType,
	arg9 []constType,
	arg10 map[int]chan string,
	arg11 map[constType]*int,
	arg12 *map[time.Time]string,
	arg13 func(int) error,
	arg14 []*func(string),
	arg15 map[int]func(sampleStruct2 *sampleStruct),
	arg16 chan *string,
	arg17 map[chan *int]*sampleStruct,
	arg18 map[string]map[*func()][]int,
	arg19 []interface{},
	arg20 map[interface{}]chan interface{},
) {

}

func (s *sampleStruct) test4(args ...*string) (error, string, time.Time) {

	return fmt.Errorf(""), "", time.Now()
}
