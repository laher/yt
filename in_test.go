package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGetSources(t *testing.T) {
	origfs := fs
	defer func() { fs = origfs }()
	fs = afero.NewMemMapFs()

	_ = fs.MkdirAll("a", 0777)
	afero.WriteFile(fs, "a.txt", []byte("file a.txt"), 0644)
	afero.WriteFile(fs, "a/b.txt", []byte("file b.txt"), 0644)
	tis := []struct {
		name    string
		sources []string
		out     map[string]source
		err     bool
	}{
		{
			name:    "single file",
			sources: []string{"a.txt"},
			out: map[string]source{
				mainSource: source{
					path: "a.txt",
					typ:  file,
				},
				"a.txt": source{
					path: "a.txt",
					typ:  file,
				},
			},
		}, {
			name:    "named file",
			sources: []string{"t=a.txt"},
			out: map[string]source{
				"t": source{
					path: "a.txt",
					typ:  file,
				},
				"t:a.txt": source{
					path: "a.txt",
					typ:  file,
				},
			},
		}, {
			name:    "named dir",
			sources: []string{"t=a/*"},
			out: map[string]source{
				"t:a/b.txt": source{
					path: "a/b.txt",
					typ:  file,
				},
			},
		}, {
			name:    "nonexistent",
			sources: []string{"t=b"},
			err:     true,
		},
	}
	for _, ti := range tis {
		t.Run(ti.name, func(t *testing.T) {
			s, err := getSources(ti.sources)
			switch ti.err {
			case false:
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(s) != len(ti.out) {
					t.Errorf("output wrong length [%d] vs expected [%d]", len(s), len(ti.out))
					return
				}
				for k, v := range ti.out {
					si, ok := s[k]
					if !ok {
						t.Errorf("Output key '%s' does not exist in map %v", k, s)
						return
					}
					if v.path != si.path || v.typ != si.typ {
						t.Errorf("Output key '%s' with value '%+v' did not match expected '%+v'", k, s[k], v)
					}
				}
			default:
				if err == nil {
					t.Errorf("Error expected")
					return
				}
			}
		})
	}
}
