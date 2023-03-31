package interfaces

// Implemented on structs that are serializable
// I think the terms 'Marshal/Unmarshal' are pretentious
// And I refuse to read them in my code.
type ISerializable interface {
	Serialize() (string, error)
	Deserialize(jsonData []byte) error
}
