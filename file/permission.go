package file

type FilePermission interface {
	HasPermission(File, string) (bool, error)
}
