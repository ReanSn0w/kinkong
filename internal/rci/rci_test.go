package rci_test

import (
	"testing"

	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/kincong/internal/rci"
)

var (
	rciClient *rci.Client

	opts = struct {
		RCI rci.Config `group:"rci" namespace:"rci" env-namespace:"RCI"`
	}{}
)

func init() {
	log, err := app.LoadConfiguration("RCI", "test", &opts, "../../run.env")
	if err != nil {
		panic(err)
	}

	rciClient, err = opts.RCI.InitRCI()
	if err != nil {
		log.Logf("[INFO] rci client init failed: %s", err)
		rciClient = nil
	}
}

func TestClient_GetInterfaces(t *testing.T) {
	if rciClient == nil {
		t.Skip("RCI client is not initialized")
	}

	data, err := rciClient.GetInterfaces()
	if err != nil {
		t.Errorf("response failed with err: %v", err)
	}

	t.Logf("data: %v", data)
}

func TestClient_GetWireguardInterfaces(t *testing.T) {
	if rciClient == nil {
		t.Skip("RCI client is not initialized")
	}

	data, err := rciClient.GetWireguardInterfaces()
	if err != nil {
		t.Errorf("response failed with err: %v", err)
	}

	t.Logf("data: %v", data)
}

func TestClient_ListIpRoutes(t *testing.T) {
	if rciClient == nil {
		t.Skip("RCI client is not initialized")
	}

	data, err := rciClient.ListIPRoutes()
	if err != nil {
		t.Errorf("response failed with err: %v", err)
	}

	t.Logf("data: %v", data)
}

func TestClient_IPRoutes(t *testing.T) {
	if rciClient == nil {
		t.Skip("RCI client is not initialized")
	}

	cases := []struct {
		Title string
		Route rci.RouteItem
	}{
		{
			Title: "Network",
			Route: rci.RouteItem{
				Network:   "197.41.128.0",
				Mask:      "255.255.128.0",
				Interface: "Wireguard1",
				Comment:   "Test Fake Route",
			},
		},
		{
			Title: "SingleIP",
			Route: rci.RouteItem{
				Network:   "197.41.128.3",
				Mask:      "255.255.255.255",
				Interface: "Wireguard1",
				Comment:   "Test Fake Single IP",
			},
		},
	}

	rci := rciClient

	for _, c := range cases {
		routes, err := rci.ListIPRoutes()
		if err != nil {
			t.Errorf("listing routes failed. case: %v", c.Title)
			t.FailNow()
		}

		for _, route := range routes {
			if route.Interface == c.Route.Interface && route.GetMask() == c.Route.Mask && route.GetNetwork() == c.Route.Network {
				t.Errorf("finded route. case: %v", c.Title)
				t.FailNow()
			}
		}

		err = rci.AddIPRoute(c.Route)
		if err != nil {
			t.Errorf("add route failed: case: %v err: %v", c.Title, err)
			t.FailNow()
		}

		routes, err = rci.ListIPRoutes()
		if err != nil {
			t.Errorf("listing routes failed.case: %v ", c.Title)
			t.FailNow()
		}

		finded := false
		for _, route := range routes {
			if route.Interface == c.Route.Interface && route.GetMask() == c.Route.Mask && route.GetNetwork() == c.Route.Network {
				finded = true
			}
		}

		if !finded {
			t.Errorf("route not found. case: %v", c.Title)
			t.FailNow()
		}

		err = rci.DeleteIPRoute(c.Route)
		if err != nil {
			t.Errorf("delete route failed. case: %v", c.Title)
			t.FailNow()
		}

		routes, err = rci.ListIPRoutes()
		if err != nil {
			t.Errorf("listing routes failed. case: %v", c.Title)
			t.FailNow()
		}

		finded = false
		for _, route := range routes {
			if route.Interface == c.Route.Interface && route.GetMask() == c.Route.Mask && route.GetNetwork() == c.Route.Network {
				finded = true
			}
		}

		if finded {
			t.Errorf("route was found. case: %v", c.Title)
			t.FailNow()
		}
	}
}
