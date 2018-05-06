package main

import (
	"testing"

	"google.golang.org/appengine/aetest"
	"fmt"
)

func Test_fetchCalendar(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	result:=fetchCalendar(ctx)
	fmt.Printf("%#v", result.Items)
}
