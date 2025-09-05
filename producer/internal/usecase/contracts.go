package usecase

type IProducerUseCase interface {
	StartTracking()
	ProcessOBUData()
}
