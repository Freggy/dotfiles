package sh_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/freggy/dotfiles/sh"
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "simple pipe",
			run: func(tt *testing.T) {
				var (
					echo = sh.Cmd("echo hello")
					cat  = sh.Cmd("cat")
				)
				pipe := sh.P().P(echo).P(cat)
				assert.NoError(tt, pipe.Err)

				expected := "hello\n"
				assert.Equal(tt, expected, string(pipe.Out))
			},
		},
		{
			name: "into",
			run: func(t *testing.T) {
				var (
					cat  = sh.Cmd("cat")
					file = fmt.Sprintf("%s/file", t.TempDir())
				)

				// echo hello | cat > %s/file
				err := sh.P().P("echo hello").P(cat).Into(file)
				assert.NoError(t, err)

				// cat %s/file | cat
				pipe := sh.P().P(cat.Append(file)).P(cat)

				expected := "hello\n"
				assert.NoError(t, pipe.Err)
				assert.Equal(t, expected, string(pipe.Out))
			},
		},
		{
			name: "check pipefail behavior",
			run: func(t *testing.T) {
				var (
					cat  = sh.Cmd("cat")
					file = fmt.Sprintf("%s/file", t.TempDir())
				)
				// test | echo hello | cat > %s/file
				err := sh.P().P("test").P("echo hello").P(cat).Into(file)
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
			t.Cleanup(func() {
				if err := os.RemoveAll(t.TempDir()); err != nil {
					t.Logf("err while tmp dir cleanup: %v", err)
				}
			})
		})
	}
}
