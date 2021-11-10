package storage

type Interface interface {
	GetFile(path string) ([]byte, error)
	UploadFile(userID int, fileName string, data []byte) (string, error)
	DeleteFile(path string) error
}

type Storage struct {
	Interface
}

func NewStorage() (Interface, error) {
	stor, err := NewMinioStorage("localhost:9000",
		"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", false)
	if err != nil {
		return nil, err
	}

	return stor, nil
}
