package migrations

type DatabaseMigrator interface {
	Up() error
	Down() error
}
