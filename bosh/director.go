package bosh

import (
	"fmt"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	"io"
	"io/ioutil"
	"path"
)

type DeployParams struct {
	NoRedact bool
	Cleanup  bool
}

type Director interface {
	Deploy(manifest string, deployParams DeployParams) error
	DownloadManifest() ([]byte, error)
}

type BoshDirector struct {
	source             concourse.Source
	commandRunner      Runner
	resourcesDirectory string
	out                io.Writer
}

func NewBoshDirector(source concourse.Source, commandRunner Runner, resourcesDirectory string, out io.Writer) BoshDirector {
	return BoshDirector{
		source:             source,
		commandRunner:      commandRunner,
		resourcesDirectory: resourcesDirectory,
		out:                out,
	}
}

func (d BoshDirector) Deploy(manifest string, deployParams DeployParams) error {
	manifestBytes, _ := ioutil.ReadFile(path.Join(d.resourcesDirectory, manifest))

	if deployParams.Cleanup {
		d.commandRunner.Execute(&boshcmd.CleanUpOpts{})
	}

	err := d.commandRunner.Execute(&boshcmd.DeployOpts{
		Args:     boshcmd.DeployArgs{Manifest: boshcmd.FileBytesArg{Bytes: manifestBytes}},
		NoRedact: deployParams.NoRedact,
	})
	if err != nil {
		return fmt.Errorf("Could not deploy: %s\n", err)
	}

	return nil
}

func (d BoshDirector) DownloadManifest() ([]byte, error) {
	bytes, err := d.commandRunner.GetResult(&boshcmd.ManifestOpts{})

	if err != nil {
		return nil, fmt.Errorf("Could not get deployment manifest: %s\n", err)
	}

	return bytes, nil
}
