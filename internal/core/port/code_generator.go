package port

type CodeGenerator interface {
	Generate() (string, error)
}
