package abstractions

type IUseCase[TInput IParamValidator, TOutput interface{}] interface {
	Execute(input TInput) (TOutput, error)
}
