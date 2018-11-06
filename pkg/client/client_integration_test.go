// +build integration

package client_test

import (
	"context"
	"crypto/sha256"
	"hash"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func TestITClient_Register(t *testing.T) {

	bits := 8
	hashFn := sha256.New
	baseURL, err := startIT(t, bits, hashFn)
	if err != nil {
		_ = stopIT(t)
	}
	t.Run("should register a new user ID", func(t *testing.T) {
		clt := client.New(
			&http.Client{},
			client.NewConfiguration(baseURL, bits, hashFn),
			client.NewInMemoryUserRepository())
		user, err := client.NewUser("new-user", 8)
		if err != nil {
			t.Errorf("NewUser() error = %v", err)
		}
		err = clt.Register(user)
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

	t.Run("should not be able to register with an existing user ID", func(t *testing.T) {
		clt := client.New(
			&http.Client{},
			client.NewConfiguration(baseURL, bits, hashFn),
			client.NewInMemoryUserRepository())
		// given
		user, _ := client.NewUser("another-new-user", 8)
		err := clt.Register(user)
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
		// when
		err = clt.Register(user)
		if err == nil {
			t.Errorf("Register() no error but got err = %v", err)
		}
	})

}

func TestITClient_Login(t *testing.T) {

	bits := 8
	hashFn := sha256.New
	baseURL, err := startIT(t, bits, hashFn)
	if err != nil {
		_ = stopIT(t)
	}
	var two = big.NewInt(2)

	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
	}

	max := new(big.Int)
	max.Exp(two, big.NewInt(int64(bits)), nil)

	clt := client.New(
		httpClient,
		client.NewConfiguration(baseURL, bits, hashFn),
		client.NewInMemoryUserRepository())

	user, err := client.NewUser("username", bits)
	if err != nil {
		t.Errorf("client.NewUser() error = %v", err)
	}

	err = clt.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should login with a valid password", func(t *testing.T) {
		err := clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Challenge()
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})

	t.Run("should recv. common error if wrong password", func(t *testing.T) {
		err := clt.Login("username", "wrong-password")
		if err == nil {
			t.Errorf("Login() want no error but got wantErr = %v", err)
		}
		err = clt.Challenge()
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})
	stopIT(t)
}

func TestITClient_GetMetadata(t *testing.T) {

	bits := 8
	hashFn := sha256.New
	baseURL, err := startIT(t, bits, hashFn)
	if err != nil {
		_ = stopIT(t)
	}
	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
	}

	clt := client.New(
		httpClient,
		client.NewConfiguration(baseURL, bits, hashFn),
		client.NewInMemoryUserRepository())

	user, _ := client.NewUser("username", bits)
	err = clt.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should have no domains", func(t *testing.T) {
		err := clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if len(domains) != 0 {
			t.Errorf("domains = %v wantDomains = %v", domains, []string{})
		}
	})

	t.Run("should have google.com domain", func(t *testing.T) {
		// given
		wantDomains := []string{"google.com"}
		err := clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Add(wantDomains[0])
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if !reflect.DeepEqual(domains, wantDomains) {
			t.Errorf("domains = %v wantDomains = %v", domains, wantDomains)
		}
	})
	stopIT(t)
}

func TestITClient_Add(t *testing.T) {

	bits := 8
	hashFn := sha256.New
	baseURL, err := startIT(t, bits, hashFn)
	if err != nil {
		_ = stopIT(t)
	}
	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
	}

	clt := client.New(
		httpClient,
		client.NewConfiguration(baseURL, bits, hashFn),
		client.NewInMemoryUserRepository())

	user, _ := client.NewUser("username", bits)
	err = clt.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should add a new domain with name google.com", func(t *testing.T) {
		// given
		wantDomains := []string{"google.com"}
		err := clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Add(wantDomains[0])
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if !reflect.DeepEqual(domains, wantDomains) {
			t.Errorf("domains = %v wantDomains = %v", domains, wantDomains)
		}
	})
	stopIT(t)
}

func TestITClient_Get(t *testing.T) {

	bits := 8
	hashFn := sha256.New
	baseURL, err := startIT(t, bits, hashFn)
	if err != nil {
		_ = stopIT(t)
	}
	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
	}

	clt := client.New(
		httpClient,
		client.NewConfiguration(baseURL, bits, hashFn),
		client.NewInMemoryUserRepository())

	user, _ := client.NewUser("username", bits)
	err = clt.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should get the same password within one session", func(t *testing.T) {
		// given
		domain := "google.com"
		err := clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Add(domain)
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		pwda, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		pwdb, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		if !reflect.DeepEqual(pwda, pwdb) {
			t.Errorf("pwda = %v pwdb = %v", pwda, pwdb)
		}
	})

	t.Run("should get the same password from two different sessions", func(t *testing.T) {
		// given
		domain := "google.com"
		err := clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Add(domain)
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		pwda, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		clt.Logout()
		err = clt.Login("username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		pwdb, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		if !reflect.DeepEqual(pwda, pwdb) {
			t.Errorf("pwda = %v pwdb = %v", pwda, pwdb)
		}
	})

	stopIT(t)
}

// === docker container utils ===
func startIT(t *testing.T, bits int, hashFn func() hash.Hash) (string, error) {
	stopIT(t)
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		return "", err
	}

	imageName := "latanassov/ossrv:0.1.0"

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	config := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9090",
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Fatal(err)
	}

	return "http://localhost:9090", nil
}

func stopIT(t *testing.T) error {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	for _, container := range containers {
		cli.ContainerStop(ctx, container.ID, nil)
	}

	return nil
}
