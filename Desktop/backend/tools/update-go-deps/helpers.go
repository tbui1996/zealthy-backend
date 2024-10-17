package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

// Structure of a Go Module (packages/router, etc)
type GoModule struct {
  Path string
  Module modfile.File
}

// Find all Go Modules in the monorepo
func GetGoModules(path string) []string {
  paths := make([]string, 0)

  // Walk dir looking for go.mod files (avoiding the root file)
  err := filepath.WalkDir(path, func(subpath string, dir fs.DirEntry, err error) error {
    if filepath.Base(subpath) == "go.mod" && subpath != "go.mod" {
      paths = append(paths, subpath)
    }
    return nil
  })

  if err != nil {
    log.Fatal(err)
  }

  return paths
}

// Create GoModule objects for each Go Module found
func LoadGoModules(paths []string) []GoModule {
  modules := make([]GoModule, 0)
  for _, path := range paths {
    modules = append(modules, GoModule{
      Module: LoadGoModule(path),
      Path: path,
    })
  }
  return modules
}

// Load in the go.mod file to read/write
func LoadGoModule(path string) modfile.File {
  fileData, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatal(err)
  }

  file, err := modfile.Parse("go.mod", fileData, nil)
  if err != nil {
    log.Fatal(err)
  }
  return *file
}

// Take all go.mod files, take latest versions of each package, and merge into one
func MergeGoModule(target modfile.File, source modfile.File) modfile.File {
  for _, req := range source.Require {
    // Check for paths inside same repository, cannot be used by Gazelle
    if strings.HasPrefix(req.Mod.Path, target.Module.Mod.Path) {
      continue
    }

    // Determine if module is new to root go.mod and add accordingly
    newModule := true
    for _, oldReq := range target.Require {
      if req.Mod.Path != oldReq.Mod.Path {
        continue
      }

      if semver.Compare(req.Mod.Version, oldReq.Mod.Version) == 1 {
        if err := target.AddRequire(req.Mod.Path, req.Mod.Version); err != nil {
          log.Fatal(err)
        }
      }
      newModule = false
    }

    if newModule {
      target.AddNewRequire(req.Mod.Path, req.Mod.Version, req.Indirect)
    }
  }

  return target
}

// Take merged go.mod and it's latest package versions and apply to all individual go.mod files
func UpdateGoModule(target modfile.File, source modfile.File) (bool, modfile.File) {
  changed := false
  for _, req := range source.Require {
    for _, oldReq := range target.Require {
      if req.Mod.Path == oldReq.Mod.Path && req.Mod.Version != oldReq.Mod.Version {
        changed = true
        if err := target.AddRequire(req.Mod.Path, req.Mod.Version); err != nil {
          log.Fatal(err)
        }
      }
    }
  }

  // Check for Go version as well
  if target.Go == nil || target.Go.Version != source.Go.Version {
    changed = true
    if err := target.AddGoStmt(source.Go.Version); err != nil {
      log.Fatal()
    }
  }

  return changed, target
}

