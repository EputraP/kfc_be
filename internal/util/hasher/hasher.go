package hasher

type Hasher interface {
	Hash(string) (string, error)
	IsEqual(hashed string, rawValue string) (bool, error)
}
