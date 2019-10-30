package config

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      []Host
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "Todos os hosts devem ser carregados",
			args: func(*testing.T) args {
				content := "host server\nHostName 192.168.0.1\nUser root\nPort 22\nhost server2\nHostName 192.168.0.2"
				return args{
					reader: strings.NewReader(content),
				}
			},
			want1: []Host{
				Host{
					Alias: "server",
					Options: map[string]string{
						"hostname": "192.168.0.1",
						"user":     "root",
						"port":     "22",
					},
				},
				Host{
					Alias: "server2",
					Options: map[string]string{
						"hostname": "192.168.0.2",
					},
				},
			},
			wantErr:    false,
			inspectErr: func(error, *testing.T) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := parseConfig(tArgs.reader)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseConfig got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("parseConfig error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
