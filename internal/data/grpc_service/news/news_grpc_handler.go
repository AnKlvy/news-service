package news

import (
	"context"
	"errors"
	"github.com/AnKlvy/news-service/internal/data/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/AnKlvy/news-service/internal/validator"
	"github.com/AnKlvy/news-service/protobuf/gen_news"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	repo database.Models
	news_proto.UnimplementedNewsServiceServer
}

func NewNewsService(grpc *grpc.Server, repo database.Models) {
	newsService := &Service{repo: repo}
	news_proto.RegisterNewsServiceServer(grpc, newsService)
}

func (s *Service) CreateNewsHandler(ctx context.Context, req *news_proto.CreateNewsRequest) (*news_proto.News, error) {
	imageURL := req.GetImageUrl()
	news := &database.News{
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
		Categories: req.GetCategories(),
		Status:     req.GetStatus(),
		ImageURL:   &imageURL,
	}

	v := validator.New()
	if database.ValidateNews(v, news); !v.Valid() {
		return nil, errors.New("invalid input data")
	}

	err := s.repo.News.Insert(news)
	if err != nil {
		return nil, err
	}

	return convertNewsToPB(news), nil
}

func (s *Service) ShowNewsHandler(ctx context.Context, req *news_proto.NewsId) (*news_proto.News, error) {
	news, err := s.repo.News.Get(req.GetId())
	if err != nil {
		return nil, err
	}

	return convertNewsToPB(news), nil
}

func (s *Service) UpdateNewsHandler(ctx context.Context, req *news_proto.UpdateNewsRequest) (*news_proto.News, error) {
	news, err := s.repo.News.Get(req.GetId())
	if err != nil {
		return nil, err
	}

	// Проверка на версию, аналогичная твоему коду
	if req.GetVersion() != news.Version {
		return nil, status.Errorf(codes.Aborted, "Version conflict: The news resource has been modified by another process.")
	}

	if req.Title != nil {
		news.Title = *req.Title
	}
	if req.Content != nil {
		news.Content = *req.Content
	}
	if len(req.GetCategories()) > 0 {
		news.Categories = req.GetCategories()
	}
	if req.Status != nil {
		news.Status = *req.Status
	}
	if req.ImageUrl != nil {
		imageURL := *req.ImageUrl
		news.ImageURL = &imageURL
	}

	v := validator.New()
	if database.ValidateNews(v, news); !v.Valid() {
		return nil, errors.New("invalid input data")
	}

	err = s.repo.News.Update(news)
	if err != nil {
		return nil, err
	}

	return convertNewsToPB(news), nil
}

func (s *Service) DeleteNewsHandler(ctx context.Context, req *news_proto.NewsId) (*emptypb.Empty, error) {
	err := s.repo.News.Delete(req.GetId())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) ListNewsHandler(ctx context.Context, req *news_proto.GetAllRequest) (*news_proto.NewsList, error) {
	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}

	filters := database.Filters{
		Page:         page,
		PageSize:     pageSize,
		Sort:         req.GetSort(),
		SortSafelist: []string{"id", "title", "status", "-id", "-title", "-status"},
	}
	v := validator.New()

	if database.ValidateFilters(v, filters); !v.Valid() {
		return nil, errors.New("invalid filters input data")
	}

	categories := req.GetCategories()
	if categories == nil {
		categories = []string{}
	}

	news, metadata, err := s.repo.News.GetAll(req.GetTitle(), categories, req.GetStatus(), filters)
	if err != nil {
		return nil, err
	}

	pbNews := make([]*news_proto.News, 0, len(news))
	for _, n := range news {
		pbNews = append(pbNews, convertNewsToPB(n))
	}

	metadataProto := &news_proto.Metadata{
		CurrentPage:  int32(metadata.CurrentPage),
		PageSize:     int32(metadata.PageSize),
		FirstPage:    int32(metadata.FirstPage),
		LastPage:     int32(metadata.LastPage),
		TotalRecords: int32(metadata.TotalRecords),
	}

	return &news_proto.NewsList{News: pbNews, Metadata: metadataProto}, nil
}

func convertNewsToPB(n *database.News) *news_proto.News {
	if n == nil {
		return &news_proto.News{}
	}
	imageURL := n.ImageURL
	return &news_proto.News{
		Id:         n.ID,
		Title:      n.Title,
		Content:    n.Content,
		Categories: n.Categories,
		Status:     n.Status,
		ImageUrl:   imageURL,
		CreatedAt:  timestamppb.New(n.CreatedAt),
		UpdatedAt:  timestamppb.New(n.UpdatedAt),
		Version:    n.Version,
	}
}
