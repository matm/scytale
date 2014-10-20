package main

import (
	"net/http"
	"os"
	"sync"

	"secret"
)

func walk(path string, info os.FileInfo, current, total int) error {
	return nil
}

type Magma struct {
	sync.Mutex
	zip *secret.ZipArchive
}

func NewMagma() *Magma {
	m := new(Magma)
	// FIXME: empty password to be set later
	m.zip = secret.NewZipArchive("")
	return m
}

type NoArgs struct{}

type ExitReply struct {
	Message string
}

func (s *Magma) Exit(r *http.Request, args *NoArgs, reply *ExitReply) error {
	// TODO: cleanup things
	defer os.Exit(0)
	return nil
}

type CreateArchiveArgs struct {
	Password   string
	OutputName string
	Files      []string
}

type CreateArchiveReply struct {
	Message string
}

func (s *Magma) CreateArchive(r *http.Request, args *CreateArchiveArgs, reply *ExitReply) error {
	s.zip.SetPassword(args.Password)
	if err := s.zip.Create(args.OutputName, args.Files, walk); err != nil {
		return err
	}
	return nil
}
