package textsgrpc

import (
	"context"

	textsv1 "github.com/slavkluev/gophkeeper-contracts/gen/go/texts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"texts/internal/domain/models"
)

type Texts interface {
	GetAll(ctx context.Context) (texts []models.Text, err error)
	SaveText(ctx context.Context, text string, info string) (textID uint64, err error)
	UpdateText(ctx context.Context, id uint64, text string, info string) (err error)
}

type serverAPI struct {
	textsv1.UnimplementedTextsServer
	texts Texts
}

func Register(gRPCServer *grpc.Server, texts Texts) {
	textsv1.RegisterTextsServer(gRPCServer, &serverAPI{texts: texts})
}

func (s *serverAPI) GetAll(
	ctx context.Context,
	in *textsv1.GetAllRequest,
) (*textsv1.GetAllResponse, error) {
	texts, err := s.texts.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all texts")
	}

	var txts []*textsv1.Text
	for _, text := range texts {
		txts = append(txts, &textsv1.Text{
			Id:   text.ID,
			Text: text.Text,
			Info: text.Info,
		})
	}

	return &textsv1.GetAllResponse{Texts: txts}, nil
}

func (s *serverAPI) Save(
	ctx context.Context,
	in *textsv1.SaveRequest,
) (*textsv1.SaveResponse, error) {
	if in.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	textID, err := s.texts.SaveText(ctx, in.GetText(), in.GetInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save text")
	}

	return &textsv1.SaveResponse{Id: textID}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	in *textsv1.UpdateRequest,
) (*textsv1.UpdateResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if in.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	err := s.texts.UpdateText(ctx, in.GetId(), in.GetText(), in.GetInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update text")
	}

	return &textsv1.UpdateResponse{}, nil
}
