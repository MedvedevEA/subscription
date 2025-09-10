package service

import (
	"subscription/internal/pkg/servererrors"
	"subscription/internal/pkg/validator"
	"subscription/internal/repository"
	repoDto "subscription/internal/repository/dto"
	svcDto "subscription/internal/service/dto"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Service interface {
	AddService(ctx *fiber.Ctx) error
	GetService(ctx *fiber.Ctx) error
	GetServices(ctx *fiber.Ctx) error
	UpdateService(ctx *fiber.Ctx) error
	RemoveService(ctx *fiber.Ctx) error

	AddSubscription(ctx *fiber.Ctx) error
	GetSubscription(ctx *fiber.Ctx) error
	GetSubscriptions(ctx *fiber.Ctx) error
	GetSubscriptionTotal(ctx *fiber.Ctx) error
	UpdateSubscription(ctx *fiber.Ctx) error
	RemoveSubscription(ctx *fiber.Ctx) error
}

type service struct {
	repo repository.Repository
	lg   *zap.SugaredLogger
}

func New(repo repository.Repository, lg *zap.SugaredLogger) *service {
	return &service{
		repo: repo,
		lg:   lg,
	}

}

func (s *service) AddService(ctx *fiber.Ctx) error {
	req := new(svcDto.AddService)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	resp, err := s.repo.AddService(*req.Name)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(resp)
}
func (s *service) GetService(ctx *fiber.Ctx) error {
	req := new(svcDto.GetService)
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	resp, err := s.repo.GetService(*req.ServiceId)
	if err == servererrors.ErrorRecordNotFound {
		return ctx.SendStatus(404)
	}
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(resp)
}
func (s *service) GetServices(ctx *fiber.Ctx) error {
	resp, err := s.repo.GetServices()
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(resp)
}
func (s *service) UpdateService(ctx *fiber.Ctx) error {
	req := new(svcDto.UpdateService)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	err := s.repo.UpdateService(&repoDto.UpdateService{
		ServiceId: *req.ServiceId,
		Name:      *req.Name,
	})
	if err == servererrors.ErrorRecordNotFound {
		return ctx.SendStatus(404)
	}
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
func (s *service) RemoveService(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveService)
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	err := s.repo.RemoveService(*req.ServiceId)
	if err == servererrors.ErrorRecordNotFound {
		return ctx.SendStatus(404)
	}
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}

func (s *service) AddSubscription(ctx *fiber.Ctx) error {

	req := new(svcDto.AddSubscription)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	resp, err := s.repo.AddSubscription(&repoDto.AddSubscription{
		ServiceName: *req.ServiceName,
		Price:       *req.Price,
		UserId:      *req.UserId,
		StartDate:   *req.StartDate,
		StopDate:    req.StopDate,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(resp)
}
func (s *service) GetSubscription(ctx *fiber.Ctx) error {
	req := new(svcDto.GetSubscription)
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	resp, err := s.repo.GetSubscription(*req.SubscriptionId)
	if err == servererrors.ErrorRecordNotFound {
		return ctx.SendStatus(404)
	}
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(resp)
}
func (s *service) GetSubscriptions(ctx *fiber.Ctx) error {
	req := new(svcDto.GetSubscriptions)
	if err := ctx.QueryParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	resp, err := s.repo.GetSubscriptions(&repoDto.GetSubscriptions{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(resp)
}
func (s *service) GetSubscriptionTotal(ctx *fiber.Ctx) error {
	req := new(svcDto.GetSubscriptionTotal)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	resp, err := s.repo.GetSubscriptionTotal(&repoDto.GetSubscriptionTotal{
		StartDate:   *req.StartDate,
		StopDate:    *req.StopDate,
		UserId:      req.UserId,
		ServiceName: req.ServiceName,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(fiber.Map{"total": resp})
}
func (s *service) UpdateSubscription(ctx *fiber.Ctx) error {
	req := new(svcDto.UpdateSubscription)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	err := s.repo.UpdateSubscription(&repoDto.UpdateSubscription{
		SubscriptionId: *req.SubscriptionId,
		ServiceName:    *req.ServiceName,
		Price:          *req.Price,
		UserId:         *req.UserId,
		StartDate:      *req.StartDate,
		StopDate:       req.StopDate,
	})
	if err == servererrors.ErrorRecordNotFound {
		return ctx.SendStatus(404)
	}
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
func (s *service) RemoveSubscription(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveSubscription)
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	err := s.repo.RemoveSubscription(*req.SubscriptionId)
	if err == servererrors.ErrorRecordNotFound {
		return ctx.SendStatus(404)
	}
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
