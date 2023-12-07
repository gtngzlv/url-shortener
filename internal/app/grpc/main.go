package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/models"
	pb "github.com/gtngzlv/url-shortener/internal/proto"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

type URLShortenerService struct {
	pb.UnimplementedUrlShortenerServiceServer
	strg storage.MyStorage
	log  zap.SugaredLogger
}

const (
	HeaderUserID = "x-user-id"
)

var (
	ErrMissingMetadata = errors.New("failed to get metadata from context")
)

func (s *URLShortenerService) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
	var resp pb.BatchResponse
	urls := in.GetEntities()
	convertedURLs := protoURLInfoToModel(urls)
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	result, err := s.strg.Batch(userID, convertedURLs)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while batch urls in storage")
	}
	resp.Entities = modelURLInfoToProto(result)
	return &resp, nil
}
func (s *URLShortenerService) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var resp pb.DeleteURLsResponse
	urls := in.GetShortURL()
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	for _, v := range urls {
		err := s.strg.DeleteByUserIDAndShort(userID, v)
		if err != nil {
			return nil, err
		}
	}
	return &resp, nil
}
func (s *URLShortenerService) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.URLInfo, error) {
	var resp pb.URLInfo
	url := in.GetShortURL()
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "no url in request")
	}
	res, err := s.strg.GetByShort(url)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while get short url in storage")
	}
	resp.OriginalUrl = res.OriginalURL
	return &resp, nil
}
func (s *URLShortenerService) GetURLs(ctx context.Context, in *pb.GetURLsRequest) (*pb.BatchResponse, error) {
	var resp pb.BatchResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}
	result, err := s.strg.GetBatchByUserID(userID)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while get urls in storage")
	}
	urls := modelURLInfoToProto(result)
	resp.Entities = urls
	return &resp, nil
}
func (s *URLShortenerService) GetStats(context.Context, *pb.GetStatisticRequest) (*pb.Statistic, error) {
	var resp pb.Statistic
	stat := s.strg.GetStatistic()
	resp.Urls = int32(stat.URLs)
	resp.Users = int32(stat.Users)
	return &resp, nil
}
func (s *URLShortenerService) Ping(context.Context, *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.strg.Ping()
	return &pb.PingResponse{}, err
}
func (s *URLShortenerService) PostAPIShorten(ctx context.Context, in *pb.APIShortenRequest) (*pb.APIShortenResponse, error) {
	var resp pb.APIShortenResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	longURL := in.GetUrl()
	res, err := s.strg.SaveFull(userID, longURL)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while post long url in storage")
	}
	resp.Result = res.ShortURL
	return &resp, nil
}
func (s *URLShortenerService) PostURL(ctx context.Context, in *pb.PostURLRequest) (*pb.PostURLResponse, error) {
	resp := pb.PostURLResponse{}
	url := in.GetLongURL()
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "no url in request")
	}
	short, err := s.strg.SaveFull(userID, url)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while post long url in storage")
	}
	resp.ShortURL = short.ShortURL
	return &resp, nil
}

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()
	st := storage.Init(log, cfg)
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterUrlShortenerServiceServer(
		s, &URLShortenerService{
			strg:                                   st,
			UnimplementedUrlShortenerServiceServer: pb.UnimplementedUrlShortenerServiceServer{},
		})
	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func protoURLInfoToModel(urls []*pb.URLInfo) []models.URLInfo {
	var convertedURLS []models.URLInfo
	for _, v := range urls {
		newURL := models.URLInfo{
			UUID:          uuid.NewString(),
			UserID:        v.UserID,
			CorrelationID: v.CorrelationId,
			OriginalURL:   v.OriginalUrl,
			ShortURL:      v.ShortUrl,
			IsDeleted:     int(v.GetIsDeleted()),
		}
		convertedURLS = append(convertedURLS, newURL)
	}
	return convertedURLS
}

func modelURLInfoToProto(urls []models.URLInfo) []*pb.URLInfo {
	var convertedURLS []*pb.URLInfo
	for _, v := range urls {
		newURL := pb.URLInfo{
			Uuid:          v.UUID,
			UserID:        v.UserID,
			CorrelationId: v.CorrelationID,
			OriginalUrl:   v.OriginalURL,
			ShortUrl:      v.ShortURL,
			IsDeleted:     int32(v.IsDeleted),
		}
		convertedURLS = append(convertedURLS, &newURL)
	}
	return convertedURLS
}

func getUserIDFromMD(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingMetadata
	}
	value, ok := GetMetadataValue(md, HeaderUserID)
	if !ok {
		return "", fmt.Errorf("failed to get %s header from metadata: no values", HeaderUserID)
	}
	return value, nil
}

func GetMetadataValue(md metadata.MD, name string) (string, bool) {
	values := md.Get(name)
	if len(values) == 0 {
		return "", false
	}

	for _, v := range values {
		if v != "" {
			return v, true
		}
	}

	return "", false
}
