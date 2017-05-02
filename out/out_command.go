package out

import (
	"github.com/starkandwayne/bosh2-errand-resource/bosh"
	"github.com/starkandwayne/bosh2-errand-resource/concourse"
	"github.com/starkandwayne/bosh2-errand-resource/storage"
)

type OutResponse struct {
	Version  concourse.Version    `json:"version"`
	Metadata []concourse.Metadata `json:"metadata"`
}

type OutCommand struct {
	director           bosh.Director
	storageClient      storage.StorageClient
	resourcesDirectory string
}

func NewOutCommand(director bosh.Director, storageClient storage.StorageClient, resourcesDirectory string) OutCommand {
	return OutCommand{
		director:           director,
		storageClient:      storageClient,
		resourcesDirectory: resourcesDirectory,
	}
}

func (c OutCommand) Run(outRequest concourse.OutRequest) (OutResponse, error) {

	runErrandParams := bosh.RunErrandParams{
		ErrandName:  outRequest.Params.ErrandName,
		KeepAlive:   outRequest.Params.KeepAlive,
		WhenChanged: outRequest.Params.WhenChanged,
	}

	if err := c.director.RunErrand(runErrandParams); err != nil {
		return OutResponse{}, err
	}

	uploadedManifest, err := c.director.DownloadManifest()
	if err != nil {
		return OutResponse{}, err
	}

	concourseOutput := OutResponse{
		Version:  concourse.NewVersion(uploadedManifest, outRequest.Source.Target),
		Metadata: append([]concourse.Metadata{}),
	}

	return concourseOutput, nil
}
