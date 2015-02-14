/*
Copyright 2014 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"syscall"
)

const (
	HTTP_GET    = "GET"
	HTTP_PUT    = "PUT"
	HTTP_DELETE = "DELETE"
	HTTP_POST   = "POST"
)

type Marathon interface {
	/* watch for changes on a application */
	Watch(name string, channel chan bool)
	/* remove me from watching this service */
	RemoveWatch(name string, channel chan bool)
	/* a list of service being watched */
	WatchList() []string
	/* check it see if a application exists */
	HasApplication(name string) (bool, error)
	/* get a listing of the application ids */
	ListApplications() ([]string, error)
	/* a list of application versions */
	ApplicationVersions(name string) (*ApplicationVersions, error)
	/* check a application version exists */
	HasApplicationVersion(name, version string) (bool, error)
	/* change an application to a different version */
	ChangeApplicationVersion(name string, version *ApplicationVersion) (*DeploymentID, error)
	/* check if an application is ok */
	ApplicationOK(name string) (bool, error)
	/* create an application in marathon */
	CreateApplication(application *Application) (bool, error)
	/* delete an application */
	DeleteApplication(application *Application) (bool, error)
	/* restart an application */
	RestartApplication(application *Application, force bool) (*Deployment, error)
	/* get a list of applications from marathon */
	Applications() (*Applications, error)
	/* get a specific application */
	Application(id string) (*Application, error)
	/* get a list of tasks for a specific application */
	Tasks(application string) (*Tasks, error)
	/* get a list of all tasks */
	AllTasks() (*Tasks, error)
	/* get a list of the deployments */
	Deployments() ([]Deployment, error)
	/* delete a deployment */
	DeleteDeployment(deployment Deployment, force bool) (Deployment, error)
	/* a list of current subscriptions */
	Subscriptions() (*Subscriptions, error)
	/* get the marathon url */
	GetMarathonURL() string
	/* ping the marathon */
	Ping() (bool, error)
	/* grab the marathon server info */
	Info() (*Info, error)
}

var (
	/* the url specified was invalid */
	ErrInvalidEndpoint = errors.New("Invalid Marathon endpoint specified")
	/* invalid or error response from marathon */
	ErrInvalidResponse = errors.New("Invalid response from Marathon")
	/* some resource does not exists */
	ErrDoesNotExist = errors.New("The resource does not exist")
	/* all the marathon endpoints are down */
	ErrMarathonDown = errors.New("All the Marathon hosts are presently down")
	/* unable to decode the response */
	ErrInvalidResult = errors.New("Unable to decode the response from Marathon")
	/* invalid argument */
	ErrInvalidArgument = errors.New("The argument passed is invalid")
)

type Client struct {
	sync.RWMutex
	/* the configuration for the client */
	config Config
	/* the callback url for subscription */
	subscription_url string
	/* the binding for the http service */
	subscription_iface string
	/* protocol */
	protocol string
	/* the http client */
	http *http.Client
	/* the marathon cluster */
	cluster Cluster
	/* a map of service you wish to listen to */
	services map[string]chan bool
}

type Message struct {
	Message string `json:"message"`
}

func NewClient(config Config) (Marathon, error) {
	/* step: we parse the url and build a cluster */
	if cluster, err := NewMarathonCluster(config.URL); err != nil {
		return nil, err
	} else {
		/* step: create the service marathon client */
		service := new(Client)
		service.services = make(map[string]chan bool, 0)
		service.cluster = cluster
		service.http = &http.Client{}
		return service, nil
	}
}

func (client *Client) GetMarathonURL() string {
	return client.cluster.Url()
}

func (client *Client) Ping() (bool, error) {
	if err := client.ApiGet(MARATHON_API_PING, "", nil); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (client *Client) MarshallJSON(data interface{}) (string, error) {
	if response, err := json.Marshal(data); err != nil {
		return "", err
	} else {
		return string(response), err
	}
}

func (client *Client) UnMarshallDataToJson(stream io.Reader, result interface{}) error {
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(result); err != nil {
		return err
	}
	return nil
}

func (client *Client) ApiGet(uri, body string, result interface{}) error {
	_, _, error := client.ApiCall(HTTP_GET, uri, body, result)
	return error
}

func (client *Client) ApiPut(uri string, post interface{}, result interface{}) error {
	var content string
	var err error
	if post == nil {
		content = ""
	} else {
		content, err = client.MarshallJSON(post)
		if err != nil {
			return err
		}
	}
	_, _, error := client.ApiCall(HTTP_PUT, uri, content, result)
	return error
}

func (client *Client) ApiPost(uri string, post interface{}, result interface{}) error {
	/* step: we need to marshall the post data into json */
	var content string
	var err error
	if post == nil {
		content = ""
	} else {
		content, err = client.MarshallJSON(post)
		if err != nil {
			return err
		}
	}
	_, _, error := client.ApiCall(HTTP_PUT, uri, content, result)
	return error
}

func (client *Client) ApiDelete(uri, body string, result interface{}) error {
	_, _, error := client.ApiCall(HTTP_DELETE, uri, body, result)
	return error
}

func (client *Client) ApiCall(method, uri, body string, result interface{}) (int, string, error) {
	if status, content, _, err := client.HttpCall(method, uri, body); err != nil {
		return 0, "", err
	} else {
		client.Debug("ApiCall() status: %s, content: %s\n", status, content)
		if status >= 200 && status <= 299 {
			if result != nil {
				if err := client.UnMarshallDataToJson(strings.NewReader(content), result); err != nil {
					return status, content, err
				}
			}
			return status, content, nil
		}
		switch status {
		case 500:
			return 0, "", ErrInvalidResponse
		case 404:
			return 0, "", ErrDoesNotExist
		}

		/* step: lets decode into a error message */
		var message Message
		if err := client.UnMarshallDataToJson(strings.NewReader(content), &message); err != nil {
			return status, content, ErrInvalidResponse
		} else {
			return status, message.Message, ErrInvalidResult
		}
	}
}

func (client *Client) HttpCall(method, uri, body string) (int, string, *http.Response, error) {
	/* step: get a member from the cluster */
	if marathon, err := client.cluster.GetMember(); err != nil {
		return 0, "", nil, err
	} else {
		url := fmt.Sprintf("%s/%s", marathon, uri)
		client.Debug("HTTPCall() method: %s, uri: %s, url: %s\n", method, uri, url)

		if request, err := http.NewRequest(method, url, strings.NewReader(body)); err != nil {
			return 0, "", nil, err
		} else {
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("X-Client", "go-marathon")
			request.Header.Add("X-Client-Version", VERSION)

			var content string
			/* step: perform the request */
			if response, err := client.http.Do(request); err != nil {
				switch error_type := err.(type) {
				case *net.OpError:
					switch error_type.Op {
					case "dial", "read":
						/* step: we need to mark the host down */
						client.cluster.MarkDown()
						/* step: retry the request */
						return client.HttpCall(method, uri, body)
					default:
					}
				case *syscall.Errno:
					switch *error_type {
					case syscall.ECONNREFUSED:
						/* step: we need to mark the host down */
						client.cluster.MarkDown()
						/* step: retry the request */
						return client.HttpCall(method, uri, body)
					}
				}
				return 0, "", response, err
			} else {
				client.Debug("HTTPCall() call successful, status: %d\n", response.StatusCode)
				/* step: lets read in any content */
				if response.ContentLength > 0 {
					client.Debug("HTTPCall() method: %s, uri: %s, url: %s\n", method, uri, url)
					/* step: read in the content from the request */
					response_content, err := ioutil.ReadAll(response.Body)
					if err != nil {
						return response.StatusCode, "", response, err
					}
					content = string(response_content)
				}
				/* step: return the request */
				return response.StatusCode, content, response, nil
			}
		}
	}
	return 0, "", nil, errors.New("Unable to make call to marathon")
}