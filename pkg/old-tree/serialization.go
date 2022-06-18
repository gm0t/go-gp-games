package old_tree

type Serialized struct {
	Key        NodeTypeKey            `json:"key"`
	Parameters map[string]interface{} `json:"parameters"`
	Children   map[string]Serialized  `json:"children"`
}
