package main

import (
	"errors"
	"fmt"
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
	reply.Message = fmt.Sprintf("Archive %s created successfully.", args.OutputName)
	return nil
}

type ExtractAllArgs struct {
	Archive   string
	Password  string
	OutputDir string
}

func (s *Magma) ExtractAll(r *http.Request, args *ExtractAllArgs, reply *ExitReply) error {
	if args.Password == "" {
		return errors.New("empty password")
	}
	if args.Archive == "" {
		return errors.New("missing archive name")
	}
	if args.OutputDir == "" {
		return errors.New("missing output directory")
	}
	s.zip.SetPassword(args.Password)
	if err := s.zip.ExtractAll(args.Archive, args.OutputDir); err != nil {
		return err
	}
	reply.Message = fmt.Sprintf("All files extracted.")
	return nil
}

type ExtractAtArgs struct {
	ExtractAllArgs
	Pos int
}

func (s *Magma) ExtractAt(r *http.Request, args *ExtractAtArgs, reply *ExitReply) error {
	if args.Password == "" {
		return errors.New("empty password")
	}
	if args.Archive == "" {
		return errors.New("missing archive name")
	}
	if args.OutputDir == "" {
		return errors.New("missing output directory")
	}
	s.zip.SetPassword(args.Password)
	if err := s.zip.ExtractAt(args.Pos, args.Archive, args.OutputDir); err != nil {
		return err
	}
	reply.Message = fmt.Sprintf("File at pos %d extracted.", args.Pos)
	return nil
}
