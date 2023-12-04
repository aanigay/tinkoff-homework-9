package file

type File struct {
	Name     string
	Contents []byte
}

type Info struct {
	Name string
	Size int64
}
