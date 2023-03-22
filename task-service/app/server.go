package main

import "task-service/app/store"

type Server struct {
	store *store.Store
}

func New(s *Server) *Server {
	return &Server{}
}
