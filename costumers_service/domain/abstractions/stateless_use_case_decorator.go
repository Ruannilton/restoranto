package abstractions

type StatelessUseCaseDecorator[TInput IParamValidator, TOutput interface{}] struct {
	useCase IUseCase[TInput, TOutput]
}

func NewStatelessUseCase[TInput IParamValidator, TOutput interface{}](useCase IUseCase[TInput, TOutput]) StatelessUseCaseDecorator[TInput, TOutput] {
	return StatelessUseCaseDecorator[TInput, TOutput]{
		useCase: useCase,
	}
}

func (decorator StatelessUseCaseDecorator[TInput, TOutput]) Execute(input TInput) (TOutput, error) {
	var zeroValue TOutput

	if validationErr := input.Validate(); validationErr != nil {
		return zeroValue, validationErr
	}

	return decorator.useCase.Execute(input)
}
