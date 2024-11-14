package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	p := "D:/test/2024/11/14/Snipaste-2.5.6-Beta-x64_1731548525393472700.rar"
	fmt.Println(strings.Split(strings.TrimSuffix(p, filepath.Ext(p)), "_")[0], filepath.Ext(p))
}
