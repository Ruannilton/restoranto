package abstractions

import "fmt"

type StatefullUseCaseDecorator[TInput IParamValidator, TOutput interface{}] struct {
	useCase    IUseCase[TInput, TOutput]
	unitOfWork IUnitOfWork
}

func NewStatefullUseCase[TInput IParamValidator, TOutput interface{}](useCase IUseCase[TInput, TOutput], unitOfWork IUnitOfWork) StatefullUseCaseDecorator[TInput, TOutput] {
	return StatefullUseCaseDecorator[TInput, TOutput]{
		useCase:    useCase,
		unitOfWork: unitOfWork,
	}
}

func (decorator StatefullUseCaseDecorator[TInput, TOutput]) Execute(input TInput) (result TOutput, err error) {
	var zeroValue TOutput

	if validationErr := input.Validate(); validationErr != nil {
		return zeroValue, validationErr
	}

	defer func() {
		if r := recover(); r != nil {
			_ = decorator.unitOfWork.Rollback()
			err = fmt.Errorf("%v", r)
		}
	}()

	defer decorator.unitOfWork.Rollback()

	result, err = decorator.useCase.Execute(input)

	if err == nil {
		decorator.unitOfWork.Commit()
	}

	return result, err
}
