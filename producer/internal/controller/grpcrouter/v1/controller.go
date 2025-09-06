package v1

import (
	v1 "github.com/evrone/go-clean-template/docs/proto/v1"
	"github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
	usecase "github.com/rauan06/realtime-map/distance-calculator"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

type V1 struct {
	v1.TranslationServer

	t usecase.IDistanceCalcUseCase
	l logger.Interface
	v *validator.Validate
}
