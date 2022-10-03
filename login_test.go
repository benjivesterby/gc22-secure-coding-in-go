package main

import "testing"

func Test_ShaTest(t *testing.T) {
	ShaTest("mostcommon.txt")
}

func Test_ArgonTest(t *testing.T) {
	ArgonTest("mostcommon.txt")
}
