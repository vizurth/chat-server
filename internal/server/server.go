package server

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	desc "github.com/vizurth/chat-server/pkg/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	desc.UnimplementedChatServer
	db *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	var chatId int64
	query := "INSERT INTO chats DEFAULT VALUES RETURNING id"
	err := s.db.QueryRow(ctx, query).Scan(&chatId)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, username := range req.Usernames {
		builderInsert := sq.Insert("chat_users").
			PlaceholderFormat(sq.Dollar).
			Columns("chat_id", "username").
			Values(chatId, username)

		query, args, err := builderInsert.ToSql()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		_, err = s.db.Exec(ctx, query, args...)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &desc.CreateResponse{
		Id: chatId,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	builderSelect := sq.Select("username").
		From("chat_users").
		Where(sq.Eq{"chat_id": req.GetId()}).PlaceholderFormat(sq.Dollar)
	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var username []string
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		username = append(username, temp)
	}
	return &desc.GetResponse{
		Usernames: username,
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	builderSelect := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "sender", "content").
		Values(req.ChatId, req.Sender, req.Content)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
