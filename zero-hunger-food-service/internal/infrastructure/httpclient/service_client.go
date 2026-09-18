package httpclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zero-hunger/food-service/internal/domain"
)

type ServiceClient struct {
	httpClient *http.Client
	userURL    string
	productURL string
}

func NewServiceClient(userURL, productURL string) *ServiceClient {
	return &ServiceClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		userURL:    userURL,
		productURL: productURL,
	}
}

func (c *ServiceClient) ValidateUserID(userID string) error {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v1/users/%s", c.userURL, userID))
	if err != nil {
		return domain.ErrFailedToReachService
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return domain.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.ErrServiceError
	}
	return nil
}
