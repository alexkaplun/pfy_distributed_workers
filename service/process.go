package service

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/alexkaplun/pfy_distributed_workers/model"
)

func (s *Service) processUrl(url *model.URL) error {
	// update url status to PROCESSING
	if err := s.updateUrlStatus(url, model.StatusProcessing); err != nil {
		return err
	}

	// get the URL, don't care for redirects for simplicity purpose
	resp, err := http.Get(url.Url)
	if err != nil {
		// update url status to error
		if errStatusUpdate := s.updateUrlStatus(url, model.StatusError); errStatusUpdate != nil {
			err = errors.Wrap(err, errStatusUpdate.Error())
		}
		return errors.Wrap(err, "failed to GET the url")
	}

	// update the url with status=DONE and received http_code
	urlUpdated, err := s.storage.UpdateURLById(url.Id, string(model.StatusDone), resp.StatusCode)
	if err != nil {
		return errors.Wrap(err, "failed to update the URL record")
	}

	if !urlUpdated {
		return errors.New(fmt.Sprintf("url does not exist, id = %v", url.Id))
	}

	return nil
}

func (s *Service) updateUrlStatus(url *model.URL, status model.UrlStatus) error {
	statusUpdated, err := s.storage.UpdateStatusById(url.Id, string(status))
	if err != nil {
		return err
	}

	if !statusUpdated {
		return errors.New(fmt.Sprintf("url does not exist, id = %v", url.Id))
	}

	return nil
}
