package docs

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/templates/docs.ghtml": {
		local:   "templates/docs.ghtml",
		size:    639,
		modtime: 1491862874,
		compressed: `
H4sIAAAAAAAA/1xQz46eIBA/16eY0jOSr2eXy+656Rs0UxmVLA4G0U1DfPcGPzTkOzEZ5vc3JUODZQJh
fL9Kxl0cR9M5qxsAgA5hCjS8iR9CdxZ6h+v6JgaEAaXlwed3+BK6U1bDh++3mThitJ6hWxfkCoEh+PMy
73WnsChs7jpi3IFxlyv1no10tJMTuvmWz1IKyCNB++H73zjSehwn/KS43N4LrCils/wpSgyVU6qU2l84
03EIfY+3oZNAZcoiTGyKWPddSlDtq0eQsmRRm9PNE3zhmpeGFxwpV5xSpHlxGAlEv/yZCI2ANndv7H7Z
Pxt7Ulfb3jvpRvn4KSrH0+P6zgoy81E4871vIRDHUtwdd3pU6Ip9xvBp/BfLv978qxRSyt2/e47EsW5f
GbuXszKX5zXj4H2k8ExZ6vkfAAD//68Y2Bt/AgAA
`,
	},

	"/templates/helping-out.md": {
		local:   "templates/helping-out.md",
		size:    3245,
		modtime: 1491862874,
		compressed: `
H4sIAAAAAAAA/4xWW4vcxhJ+16+otR9iw1xw4sSQt9hLLpDEhjWEYAIuqUtSe9VdfbpLO6sT8t8P1d0z
o7F3D3mUuq5ffXX584ef3l2/BpugJesHcGgI0BuI6KFd4AUM8wIobgOOwC2Jpj6/2x3cMKBfYKQpgIUO
PQwkagpDiNRZFDK7pnk/UiT9neiOIk5wwCXBwnNWydoOb9V5DaYlEYqbLGLYfyUw4h2BMNx6PmSfHRtV
EM76V00DW/gjWtGfhrvZkRcUy34H70eClgU8kUlwvX68ghsi6OcoI0VoaeID9BzBsQbsewb2cHjQ7Nph
x4byj59pCvqDZwHONgNxmNQWGJs6jibL/Wjv4TCiKB4gS+CUIR0iOocRKEaOOfveegPPJAOIkcBp7tyD
jOSeq6XX86IfgAZSwI6a5unTp6fAZLSfZdw054Ib6q23QtMC7KmaLbl37MJkOy2gYpdyRjmMTY5UBSNh
Yg+YIEkRpF51FUAZMRPBQojUU1SD7ZzEtnaysmQbibo56gcrChPzbdoUH2oBYWLRmC5gV5taRzLZTUVX
GPp5mhaYxU72v4W/szcUkxyj7QlljpTyh9IhRL6zhtKuad6M1N3mh2G2hioRbK8l+CozQShS0hwz8/ND
4aTW4z8zJQ0u21ZAbrVVeFVyLcpTuLFOoy0+2B9Jj9CN6IecBn6WbsCBmgZe7OBHjre1PxoAgC28iYRC
+195sF51Byvj3NbH3/HODijZqCYWIn+iTrJB+KAg//VsFAnp+/2+KO46dvtP7DG9evlqv+AQTPu8WnsX
KRXkeg2jnUXYN/D1Dv7kGdLI86TlB2xLNRIRJHYkoyY42VsqVMz4XH3IubyrEf2OjuAXhwP99ezvv5VK
toMnhrttBmr7Yqs+d8EPT/7553kD3+xq5oDg6QBtRN+N0Ed2OUJDd/XfF8FX0RI+5Ka6+vD64qd9PJCv
t0V/W0TPEZVqTLZy6IsA3vqOjlRi/5nMBrqT5mV8OKD1pduWkEcIglewlHwH9B5hsHf5M146vPpQEaoG
/09S3xyT6rLGOamXO/ihF6qtvFyUmX1xeoZ/80D8VV4ndsIFPhagv1/HW/LRQnxs4Nsd/IaZKVQ7ItV8
fulXOTuVWfdMmObB+m0K1NnednkhbE7rJU/QMo48fKyy6naPKZGkvYpvtS12znz8Vw4H8nmNXToyfHZj
uEt7IRcmFEr7+/v7+5X1t4EKC3o7ZTBLA5bJGk5sJWPliOWpsOfKvdyqwIqMDxf4C7EVZzPcuRwV7xKC
LrycMTtnBRylhANtwFDqom21pQ9HUhhrHgru221R3lblx+N7SHIV4rl5azTH2fPdii1hniaIlOdw1fuJ
oUVtqzL+IrNAz5OhCM8GLgu6B5ym0o42UiccLaWj45vT+Ko8xhAI4wYOVkbAY2E83Ys6saIkV/Enb9gF
3dW5mqvInjyA03fbivzjAF2IrJC5IYEW06nrhOGa7h7w8WqrYrXPH/fzhdjK13WpPF3W/UzXOm4u6/Bq
B3+glbyoI6XAPlEZ045K2zg7jJJvRl0WucLkJZXKZZxD5MCJzJmgrhSdUwF45EBl91uBg50mnU+O4kDm
qinH0Bv2Em07n061fJHmG8CmTHYhr2078Pm4OVBb1n6PXb5d56TqYdKZ/PP7334t8bXMkiRiKBeNZpGc
suoT3qFiFqTsmWo30m41XDJzjuPlFCX7sjF7e595NmwAjQE8njDVVZ3CR3qGyC220wKexfaLXiC9jUn0
hj8dGjalmYAjJPL5fnbqunbe6mJR86e5Xy7wAyapw0Ksy+uoU/ZZfQGcEsOBY/aTdxxKOXv+zU3cNNv3
b6/fbovC63kpt+v6qi0C/wsAAP//WuEmt60MAAA=
`,
	},

	"/templates/quickstart.md": {
		local:   "templates/quickstart.md",
		size:    187,
		modtime: 1491862874,
		compressed: `
H4sIAAAAAAAA/xzLsU4DMQwG4N1P8XMTZKjEK8DCwHgbQqfQsxKrJgmxA62q5tkR7N9HaxZDi4nxI6rI
rA2XOvA15HjSC4wdo8Ez46P6geiFO0MMsUA+/x7dvUV1OJ/9/f56NY8uRyzO5v/g0EpabjcsrzVVrOLK
WPnseFweiEJ4qrqHQBSeRzf55kC0baPs3FUK79tGNKd5lxN77nWkPCei7b8BAAD//yQZdPa7AAAA
`,
	},

	"/templates/templates.md": {
		local:   "templates/templates.md",
		size:    284,
		modtime: 1491862874,
		compressed: `
H4sIAAAAAAAA/1TPsU4DMQzG8f2e4ttopZbszEhMqEs3xGAlvsT0Ekdnh7Zvj66IgfEv/ST7e1dzpHuj
KhHON4eNWECGL5WGymaU2Q6Iw1wrotZKLRmopV8qhsyNV3JOGCYtw7n2hZwNu6orwzpHmSXSstzx8XZ6
ssep8OfQKV4o8+euuHd7CSHrQi0/65pDv+TwT4f9fprORQydMoOkGlzBt76QNBS9bpnZ4YVRt4E6HDpv
XQ+Pz6+FHPNo0UWbgVZGX/VbEqdpOuJ8ej3h+BMAAP///0PdJBwBAAA=
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/templates": {
		isDir: true,
		local: "templates",
	},
}
