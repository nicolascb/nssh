package actions

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/nicolascb/nssh/internal/config"
)

func Test_parseURI(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 string
		want2 string
		want3 string
	}{
		{
			name: "Test parse complete uri",
			args: func(*testing.T) args {
				return args{
					uri: "root@localhost:22",
				}
			},
			want1: "root",
			want2: "22",
			want3: "localhost",
		},
		{
			name: "Test without port",
			args: func(*testing.T) args {
				return args{
					uri: "root@localhost",
				}
			},
			want1: "root",
			want2: "",
			want3: "localhost",
		},
		{
			name: "Test without user and port",
			args: func(*testing.T) args {
				return args{
					uri: "localhost",
				}
			},
			want1: "",
			want2: "",
			want3: "localhost",
		},
		{
			name: "Test without user with port",
			args: func(*testing.T) args {
				return args{
					uri: "localhost:22",
				}
			},
			want1: "",
			want2: "22",
			want3: "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, got2, got3 := parseURI(tArgs.uri)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseURI got1 = %v, want1: %v", got1, tt.want1)
			}

			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("parseURI got2 = %v, want2: %v", got2, tt.want2)
			}

			if !reflect.DeepEqual(got3, tt.want3) {
				t.Errorf("parseURI got3 = %v, want3: %v", got3, tt.want3)
			}
		})
	}
}

func Test_getHostOptions(t *testing.T) {
	type args struct {
		alias   string
		uri     string
		sshkey  string
		options []string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]string
	}{
		{
			name: "Get full options",
			args: func(*testing.T) args {
				return args{
					alias:  "server",
					uri:    "root@nicolascb.com:22",
					sshkey: "~/.ssh/mykey",
					options: []string{
						"opt1=val1",
						"opt2=val2",
					},
				}
			},
			want1: map[string]string{
				"user":         "root",
				"port":         "22",
				"hostname":     "nicolascb.com",
				"identityfile": "~/.ssh/mykey",
				"opt1":         "val1",
				"opt2":         "val2",
			},
		},
		{
			name: "Test with general",
			args: func(*testing.T) args {
				return args{
					alias:  config.GeneralKey,
					uri:    "root@nicolascb.com:22",
					sshkey: "~/.ssh/mykey",
					options: []string{
						"opt1=val1",
						"opt2=val2",
					},
				}
			},
			want1: map[string]string{
				"identityfile": "~/.ssh/mykey",
				"opt1":         "val1",
				"opt2":         "val2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := getHostOptions(tArgs.alias, tArgs.uri, tArgs.sshkey, tArgs.options)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getHostOptions got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_confirmProceedUpdate(t *testing.T) {
	type args struct {
		out io.Reader
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 bool
	}{
		{
			name: "Press Y return true",
			args: func(*testing.T) args {
				return args{
					out: strings.NewReader("Y\n"),
				}
			},
			want1: true,
		},
		{
			name: "Press y lowcase test",
			args: func(*testing.T) args {
				return args{
					out: strings.NewReader("y\n"),
				}
			},
			want1: true,
		},
		{
			name: "dont press enter return false",
			args: func(*testing.T) args {
				return args{
					out: strings.NewReader("y"),
				}
			},
			want1: false,
		},
		{
			name: "dont press y return false",
			args: func(*testing.T) args {
				return args{
					out: strings.NewReader("x\n"),
				}
			},
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := confirmProceedUpdate(tArgs.out)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("confirmProceedUpdate got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
