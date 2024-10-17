package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// update-go-deps so that this hits on a search when someone is looking for it

//nolint
func main() {
  log.SetFlags(0)

  // Change to WORKSPACE working directory
  workingDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
  err := os.Chdir(workingDir)
  if err != nil {
    log.Fatal(err)
  }

  // Find all packages from source modules (and print paths for debugging)
  goPackages := flag.Args()
  if len(goPackages) == 0 {
    goPackages = GetGoModules("./packages")
    log.Printf("Got %v go packages and go.mod files to merge:", len(goPackages))
    for _, goPackage := range goPackages {
      fmt.Println("\t", goPackage)
    }
  }

  // Use the helper function to load the root go.mod file
  rootGoModule := LoadGoModule("go.mod")

  // Load all of the found Go modules in the package directory (packages/router/go.mod, etc)
  goModules := LoadGoModules(goPackages)

  // Iterate on Go module packages, retrieving highest verions for modules, and any new modules
  rootGoModule.SetRequire(make([]*modfile.Require, 0))
  for _, pkg := range goPackages {
		trimmed := strings.Split(pkg, "/go.mod")
		rootGoModule.AddNewRequire("github.com/circulohealth/sonar-backend/"+trimmed[0], "v0.0.0-00010101000000-000000000000", false)
	}
  for _, source := range goModules {
    rootGoModule = MergeGoModule(rootGoModule, source.Module)
  }

  // Read in the root go.mod file for read/write
  fileData, err := rootGoModule.Format()
  if err != nil {
    log.Fatal(err)
  }
  // Write updated changes to the root go.mod file
  //nolint
  if err = ioutil.WriteFile("go.mod", fileData, 644); err != nil {
    log.Fatal(err)
  }
  log.Printf("Updated root go.mod file")

  // Download the newly edited root go.mod modules
  out, err := exec.Command("go", "get", ".").CombinedOutput()
  if err != nil {
    log.Fatal(string(out))
  }
  log.Println("Installed packages from updated root go.mod and updated go.sum")

  // Run gazelle on the root go.mod file to update the go_repositories.bzl file containing the
  // bazel configurations for every used go module across the monorepo
  out, err = exec.Command("bazel", "run", "//:gazelle", "--", "update-repos",
          "-from_file=go.mod", "-prune", "-to_macro=go_repositories.bzl%go_repositories").CombinedOutput()
  if err != nil {
    log.Fatal(string(out))
  }
  log.Println("Updated root go_repositories.bzl with gazelle update-repos command")

  // Take merged go.mod and it's latest package versions and apply to all individual go.mod files
  for _, source := range goModules {
    // Update specific go.mod file (packages/router/go.mod, etc)
    changed, file := UpdateGoModule(source.Module, rootGoModule)
    if !changed {
      continue
    }

    // Format file contents for write
    fileContent, err := file.Format()
    if err != nil {
      log.Fatal(err)
    }

    // Write edited go.mod file
    //nolint
    if err = ioutil.WriteFile(source.Path, fileContent, 644); err != nil {
      log.Fatal(err)
    }

    // Change to the directory of the edited go.mod file for module downloads
    if err = os.Chdir(filepath.Dir(source.Path)); err != nil {
      log.Fatal(err)
    }

    // Download the newly edited go.mod modules
    if out, err = exec.Command("go", "get", "-v", ".").CombinedOutput(); err != nil {
      log.Fatal(string(out))
    }

    // Change back to root of workspace to continue on for other go.mod files
    log.Printf("Changed package go.mod and go.sum: %v", source.Path)
    if err = os.Chdir(workingDir); err != nil {
      log.Fatal(err)
    }
  }
}
