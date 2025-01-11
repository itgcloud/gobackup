package compressor

import (
	"strings"
	"testing"

	"github.com/longbridgeapp/assert"
	"github.com/spf13/viper"

	"github.com/itgcloud/gobackup/config"
	"github.com/itgcloud/gobackup/helper"
)

func TestTar_options(t *testing.T) {
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
			wantOptsGnu: "--ignore-failed-read -a -cP --exclude=/hello/world --exclude=/cc/111 -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
			wantOpts:    "-a -cP --exclude=/hello/world --exclude=/cc/111 -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
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
			wantOptsGnu: "--ignore-failed-read -a -cP -h --exclude=/hello/world --exclude=/cc/111 -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
			wantOpts:    "-a -cP -h --exclude=/hello/world --exclude=/cc/111 -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
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
			wantOptsGnu: "--ignore-failed-read -a -cP -h -s -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
			wantOpts:    "-a -cP -h -s -f ~/work/dir/archive.tar /foo/bar/dar /bar/foo /ddd",
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

			tar := &Tar{
				Base: Base{model: model},
			}
			opts, err := tar.options("~/work/dir/archive.tar")
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			args := strings.Join(opts, " ")
			if helper.IsGnuTar {
				assert.Equal(t, tt.wantOptsGnu, args)
			} else {
				assert.Equal(t, tt.wantOpts, args)
			}
		})
	}

}
