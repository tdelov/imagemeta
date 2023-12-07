package xmp

import (
    "testing"
_    "fmt"
  _  "regexp"
)




func TestIpmcDateString(t *testing.T) {
    msg, err := parseDate("2021-01-10T17:30:57.04")
    if err != nil {
    	t.Fatalf(`parseDate("2021-01-10T17:30:57.04") = %q, %v, want "", error`, msg, err)
    }
    t.Log("Sent 2021-01-10T17:30:57.04 Got:" , msg )
    return
}


func TestExifDateString1(t *testing.T) {
    msg, err := parseDate("2007:09:11 13:53:33.006")
    if err != nil {
        t.Fatalf(`parseDate("2007:09:11 13:53:33.006") = %q, %v, want "", error`, msg, err)
    }  
     	t.Log("Sent 2007:09:11 13:53:33.006. Got:" , msg )
   	return
}

func TestExifDateString2(t *testing.T) {
    msg, err := parseDate("2007:09:11T13:53:33.006")
    if err != nil {
        t.Fatalf(`parseDate("2007:09:11T13:53:33.006") = %q, %v, want "", error`, msg, err)
    }  
     	t.Log("Sent 2007:09:11T13:53:33.006. Got:" , msg )
   	return
}

func TestExifDateStringEmpty(t *testing.T) {
    msg, err := parseDate("")
    if err != nil {
	t.Fatalf(`parseDate("") = %q, %v, want "", error`, msg, err)
    }  
	t.Log("Sent  Got:" , msg )
	return
}

