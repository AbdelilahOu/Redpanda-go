package consumer

type Event[T any] struct {
	Schema  Schema     `json:"schema"`
	Payload Payload[T] `json:"payload"`
}

type Payload[T any] struct {
	Before      *T          `json:"before"`
	After       *T          `json:"after"`
	Source      Source      `json:"source"`
	Op          string      `json:"op"`
	TsMs        int64       `json:"ts_ms"`
	Transaction interface{} `json:"transaction"`
}

type Schema struct {
	Type   string `json:"type"`
	Fields []struct {
		Type     string `json:"type"`
		Field    string `json:"field"`
		Optional bool   `json:"optional"`
	} `json:"fields"`
	Optional bool   `json:"optional"`
	Name     string `json:"name"`
}

type Source struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	Snapshot  string `json:"snapshot"`
	Db        string `json:"db"`
	Table     string `json:"table"`
	ServerID  int    `json:"server_id"`
	GTID      string `json:"gtid"`
	File      string `json:"file"`
	Pos       int    `json:"pos"`
	Row       int    `json:"row"`
	Thread    int    `json:"thread"`
	Query     string `json:"query"`
}
