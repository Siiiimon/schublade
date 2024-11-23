package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestScan(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []Event
	}{
		"hello world": {
			input: "hello world",
			want: []Event{
				{
					name: "Text",
					data: "hello world",
				},
			},
		},
		"empty": {
			input: "",
			want:  []Event{},
		},
		"null": {
			input: "\x00",
			want:  []Event{},
		},
		"Carriage Return": {
			input: "\r",
			want: []Event{
				{
					name: "CarriageReturn",
					data: nil,
				},
			},
		},
		"string with carriage return": {
			input: "foobar\r",
			want: []Event{
				{
					name: "Text",
					data: "foobar",
				},
				{
					name: "CarriageReturn",
					data: nil,
				},
			},
		},
		"bell": {
			input: "\x07",
			want: []Event{
				{
					name: "Bell",
					data: nil,
				},
			},
		},
		"backspace": {
			input: "\x08",
			want: []Event{
				{
					name: "Backspace",
					data: nil,
				},
			},
		},
		"tab": {
			input: "\t",
			want: []Event{
				{
					name: "Tab",
					data: nil,
				},
			},
		},
		"line feed": {
			input: "\n",
			want: []Event{
				{
					name: "LineFeed",
					data: nil,
				},
			},
		},
		"vertical tab": {
			input: "\v",
			want: []Event{
				{
					name: "VerticalTab",
					data: nil,
				},
			},
		},
		"form feed": {
			input: "\x0C",
			want: []Event{
				{
					name: "FormFeed",
					data: nil,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rc := io.NopCloser(bytes.NewReader([]byte(test.input)))
			parseChan := make(chan Event)
			var got []Event

			go Parse(rc, parseChan)
			for ev := range parseChan {
				got = append(got, ev)
			}

			if len(test.want) != len(got) {
				t.Errorf("expected %d events, got %d", len(test.want), len(got))
				println("Events:")
				for _, ev := range got {
					fmt.Printf("%#v\n", ev)
				}
				t.FailNow()
			}

			for i, w := range test.want {
				g := got[i]

				if w.name != g.name {
					t.Errorf("event[%d]: expected name %q, got %q", i, w.name, g.name)
				}

				if w.data != nil && g.data != nil {
					if w.data != g.data {
						t.Errorf("event[%d]: expected data %q, got %q", i, w.data, g.data)
					}
				} else if w.data != g.data {
					t.Errorf("event[%d]: expected data %v, got %v", i, w.data, g.data)
				}
			}
		})
	}
}
