package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"secret"
)

func walk(path string, info os.FileInfo, current, total int) error {
	log.Printf("Adding %s [%d/%d]\n", path, current, total)
	return nil
}

type Magma struct {
	sync.Mutex
	zip *secret.ZipArchive
}

const emptyPassword = ""

func NewMagma() *Magma {
	m := new(Magma)
	m.zip = secret.NewZipArchive(emptyPassword)
	return m
}

type NoArgs struct{}

type ExitReply struct {
	Message string
}

// Exit terminates the process.
func (s *Magma) Exit(r *http.Request, args *NoArgs, reply *ExitReply) error {
	defer os.Exit(0)
	// Stops any processing
	s.zip.Cancel()
	log.Println("Exiting...")
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
		return errors.New("empty password")
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
