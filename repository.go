package main

type Repository interface {
	InsertEvent(e Event) error
}

func MakeRepository(output Output) Repository {
	var repo Repository
	switch output.Kind {
	case "sql":
		repo = ConnectToDatabase(output.ConnectionString)
	case "file":
		repo = CreateFileRepository(output.ConnectionString)
	default:
		panic("Unknown output kind. Valid kinds: [sql,file]")
	}

	return repo
}
