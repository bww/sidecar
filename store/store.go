package store

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v3"
)

const storeRoot = ".sidecar"

var (
	ErrNoUserHome = errors.New("No user home directory")
	ErrFileExists = errors.New("File exists")
)

type Resource string

const (
	Config = Resource("profiles.yaml")
)

func Find(r Resource) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	if usr.HomeDir == "" {
		return "", ErrNoUserHome
	}
	return path.Join(usr.HomeDir, storeRoot, string(r)), nil
}

func Load(p string, e interface{}) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	return Read(f, e)
}

func Read(r io.Reader, e interface{}) error {
	err := yaml.NewDecoder(r).Decode(e)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(p string, b []byte, overwrite bool) error {
	_, err := os.Stat(p)
	if err == nil && !overwrite {
		return ErrFileExists
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(path.Dir(p), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, b, 0600)
}

func Store(p string, e interface{}) error {
	err := os.MkdirAll(path.Dir(p), 0755)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, e)
}

func Write(w io.Writer, e interface{}) error {
	err := yaml.NewEncoder(w).Encode(e)
	if err != nil {
		return err
	}
	return nil
}
