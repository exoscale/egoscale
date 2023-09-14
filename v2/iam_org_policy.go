package v2

import (
	"context"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
)

// GetIAMOrgPolicy returns the IAM Organization policy.
func (c *Client) GetIAMOrgPolicy(ctx context.Context, zone string) (*IAMPolicy, error) {
	resp, err := c.GetIamOrganizationPolicyWithResponse(apiv2.WithZone(ctx, zone))
	if err != nil {
		return nil, err
	}

	return iamPolicyFromAPI(resp.JSON200), nil
}

// UpdateIAMOrgPolicy updates existing IAM Organization policy.
func (c *Client) UpdateIAMOrgPolicy(ctx context.Context, zone string, policy *IAMPolicy) error {
	services := map[string]oapi.IamServicePolicy{}
	if len(policy.Services) > 0 {
		for name, service := range policy.Services {
			rules := []oapi.IamServicePolicyRule{}
			for _, rule := range service.Rules {
				resources := []string{}
				if rule.Resources != nil {
					resources = rule.Resources
				}

				rules = append(rules, oapi.IamServicePolicyRule{
					Action:     (*oapi.IamServicePolicyRuleAction)(rule.Action),
					Expression: rule.Expression,
					Resources:  &resources,
				})
			}

			services[name] = oapi.IamServicePolicy{
				Type:  (*oapi.IamServicePolicyType)(service.Type),
				Rules: &rules,
			}
		}
	}

	resp, err := c.UpdateIamOrganizationPolicyWithResponse(
		apiv2.WithZone(ctx, zone),
		oapi.UpdateIamOrganizationPolicyJSONRequestBody{
			DefaultServiceStrategy: oapi.IamPolicyDefaultServiceStrategy(policy.DefaultServiceStrategy),
			Services: oapi.IamPolicy_Services{
				AdditionalProperties: services,
			},
		})
	if err != nil {
		return err
	}

	_, err = oapi.NewPoller().
		WithTimeout(c.timeout).
		WithInterval(c.pollInterval).
		Poll(ctx, oapi.OperationPoller(c, zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}
