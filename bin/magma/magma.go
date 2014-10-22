package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"secret"
)

func walk(path string, info os.FileInfo, current, total int) error {
	fmt.Println(current, path)
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
	defer os.Exit(0)
	s.zip.Cancel()
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
	if args.Password == "" {
		return errors.New("empty passwor")
	}
	if args.OutputName == "" {
		return errors.New("missing output name")
	}
	if len(args.Files) == 0 {
		return errors.New("empty file set")
	}
	s.zip.SetPassword(args.Password)
	if err := s.zip.Create(args.OutputName, args.Files, walk); err != nil {
		return err
	}
	return nil
}
