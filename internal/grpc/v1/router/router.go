package router

import (
	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/grpc/v1/controller"
	"google.golang.org/grpc"
)

func Register(
	srv *grpc.Server,
	ctrl *controller.Controller,
) {
	servicev1.RegisterGoLucaServiceServer(srv, ctrl)
}
