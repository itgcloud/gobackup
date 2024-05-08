package archive

import (
	"strings"
	"testing"

	"github.com/longbridgeapp/assert"
	"github.com/spf13/viper"

	"github.com/itgcloud/gobackup/config"
	"github.com/itgcloud/gobackup/helper"
)

func TestRun(t *testing.T) {
	// with nil Archive
	model := config.ModelConfig{
		Archive: nil,
	}
	err := Run(model)
	assert.NoError(t, err)
}

func TestOptions(t *testing.T) {
	type args struct {
		excludes       []string
		includes       []string
		additionalArgs []string
	}
	tests := []struct {
		name        string
		args        args
		wantOptsGnu string
		wantOpts    string
		wantErr     bool
	}{
		{
			name: "default test",
			args: args{
				includes: []string{
					"/foo/bar/dar",
					"/bar/foo",
					"/ddd",
				},
				excludes: []string{
					"/hello/world",
					"/cc/111",
				},
				additionalArgs: []string{},
			},
			wantOptsGnu: "--ignore-failed-read -cP -f ~/work/dir/archive.tar --exclude=/hello/world --exclude=/cc/111 /foo/bar/dar /bar/foo /ddd",
			wantOpts:    "-cP -f ~/work/dir/archive.tar --exclude=/hello/world --exclude=/cc/111 /foo/bar/dar /bar/foo /ddd",
			wantErr:     false,
		},
		{
			name: "test with additional arguments",
			args: args{
				includes: []string{
					"/foo/bar/dar",
					"/bar/foo",
					"/ddd",
				},
				excludes: []string{
					"/hello/world",
					"/cc/111",
				},
				additionalArgs: []string{"-h"},
			},
			wantOptsGnu: "--ignore-failed-read -cP -h -f ~/work/dir/archive.tar --exclude=/hello/world --exclude=/cc/111 /foo/bar/dar /bar/foo /ddd",
			wantOpts:    "-cP -h -f ~/work/dir/archive.tar --exclude=/hello/world --exclude=/cc/111 /foo/bar/dar /bar/foo /ddd",
			wantErr:     false,
		},
		{
			name: "test with multiple arguments",
			args: args{
				includes: []string{
					"/foo/bar/dar",
					"/bar/foo",
					"/ddd",
				},
				excludes:       []string{},
				additionalArgs: []string{"-h", "-s"},
			},
			wantOptsGnu: "--ignore-failed-read -cP -h -s -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
			wantOpts:    "-cP -h -s -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
			wantErr:     false,
		},
		{
			name: "test without includes",
			args: args{
				includes:       []string{},
				excludes:       []string{},
				additionalArgs: []string{},
			},
			wantOptsGnu: "",
			wantOpts:    "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			model := config.ModelConfig{
				DumpPath: "~/work/dir",
				Archive:  viper.New(),
			}
			// init model
			model.Archive.Set("additional_arguments", tt.args.additionalArgs)
			model.Archive.Set("includes", tt.args.includes)
			model.Archive.Set("excludes", tt.args.excludes)

			opts, err := options(model)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			args := strings.Join(opts, " ")
			if helper.IsGnuTar {
				assert.Equal(t, args, tt.wantOptsGnu)
			} else {
				assert.Equal(t, args, tt.wantOpts)
			}
		})
	}
}
