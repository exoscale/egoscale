package runstatus

import (
	"github.com/exoscale/egoscale/api"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// Page represents a Runstatus page.
type Page struct {
	api.Resource

	ID                   int
	Name                 string
	Title                string
	DefaultStatusMessage string
	CustomDomain         string
	TimeZone             string

	c *Client
}

// TODO: Page.AddService()

// TODO: Page.Services()

// TODO: Page.AddIncident()

// TODO: Page.Incidents()

// TODO: Page.AddMaintenance()

// TODO: Page.Maintenances()

// TODO: Page.Update() // <- Not implemented in internal/egoscale

// TODO: Page.Delete()

// TODO: CreatePage()

// ListPages returns the list of Runstatus pages owned, or an error if the API query failed.
func (c *Client) ListPages() ([]*Page, error) {
	pages := make([]*Page, 0)

	res, err := c.c.ListRunstatusPages(c.ctx)
	if err != nil {
		return nil, err
	}

	for _, page := range res {
		page := page
		pages = append(pages, c.pageFromAPI(&page))
	}

	return pages, nil
}

// TODO: GetPage()

func (c *Client) pageFromAPI(page *egoapi.RunstatusPage) *Page {
	return &Page{
		Resource:             api.MarshalResource(page),
		ID:                   page.ID,
		Name:                 page.Subdomain,
		Title:                page.Title,
		DefaultStatusMessage: page.OkText,
		CustomDomain:         page.Domain,
		TimeZone:             page.TimeZone,
		c:                    c,
	}
}
