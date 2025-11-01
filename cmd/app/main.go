package main

import (
	"fmt"
	"strings"

	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/kincong/internal/configuration"
	"github.com/ReanSn0w/kincong/internal/rci"
	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/resolver/asn"
	"github.com/ReanSn0w/kincong/internal/resolver/dns"
	"github.com/ReanSn0w/kincong/internal/resolver/ip"
	"github.com/ReanSn0w/kincong/internal/utils"
)

const (
	ActionUpload      Action = "upload"
	ActionFileInspect Action = "file-inspect"
	ActionFilePreview Action = "file-preview"
	ActionRCIClear    Action = "rci-clear"
	ActionRCIList     Action = "rci-list"
	ActionRCIInspect  Action = "rci-inspect"
)

type Action string

var (
	revision    = "main"
	application *app.Application
	opts        = struct {
		app.Debug

		Action struct {
			Upload bool `long:"upload" description:"upload file params to rci (default action)"`

			File struct {
				Inspect bool `long:"inspect" description:"validate file"`
				Preview bool `long:"preview" description:"preview file"`
			} `group:"File Actions" namespace:"file"`

			RCI struct {
				List    bool `long:"list" description:"list avaliable rci interfaces"`
				Inspect bool `long:"inspect" description:"list avaliable static routes"`
				Clear   bool `long:"clear" description:"clear all static routes"`
			} `group:"RCI Actions" namespace:"rci"`
		} `group:"Main Actions" namespace:"action"`

		Interface string `short:"i" long:"interface" description:"select interface"`
		File      string `short:"f" long:"file" description:"select file"`

		RCI rci.Config `group:"RCI Params" namespace:"rci" env-namespace:"RCI"`
		ASN asn.Config `group:"ASN Resolver Params" namespace:"asn" env-namespace:"ASN"`
		DNS dns.Config `group:"DNS Resolver Params" namespace:"dns" env-namespace:"DNS"`
	}{}
)

func main() {
	application = app.New("Kincong", revision, &opts)

	var err error

	action := selectAction()
	switch action {
	case ActionUpload:
		err = uploadFn()
	case ActionFileInspect:
		err = fileInspectFn()
	case ActionFilePreview:
		err = filePreviewFn()
	case ActionRCIClear:
		err = rciClearFn()
	case ActionRCIList:
		err = rciListFn()
	case ActionRCIInspect:
		err = rciInspectFn()
	default:
		err = fmt.Errorf("unknown action")
	}

	if err != nil {
		application.Log().Logf("[ERROR] application stopped with error: %s", err)
	}
}

func uploadFn() error {
	conf, err := configuration.Load(opts.File)
	if err != nil {
		return err
	}

	rc := opts.RCI.MustInitRCI()

	rv, err := initResolver()
	if err != nil {
		return err
	}

	interfaces, err := rc.GetInterfaces()
	if err != nil {
		return err
	}

	// Проверка наличия интерфейса на роутере
	{
		err = fmt.Errorf("interface %s not found", opts.Interface)
		for name := range interfaces {
			if name == opts.Interface {
				err = nil
				break
			}
		}
		if err != nil {
			return err
		}
	}

	// Сбор параметров из конфигурации
	var values []utils.Value
	{
		for _, value := range conf {
			application.Log().Logf("[INFO] сбор параметров из %s", value.Title)
			values = append(values, value.Values...)
		}
	}

	subnets, err := rv.Resolve(values...)
	if err != nil {
		return err
	}

	application.Log().Logf("[INFO] сформирован список из %d подсетей", len(subnets))

	for _, subnet := range subnets {
		ip, mask, err := utils.SubnetToIPAndMask(string(subnet))
		if err != nil {
			application.Log().Logf("[ERROR] upload subnet err: %s", err)
			continue
		}

		item := rci.RouteItem{
			Network:   ip,
			Mask:      mask,
			Interface: opts.Interface,
		}

		item.Comment = item.Hash()

		err = rc.AddIPRoute(item)
		if err != nil {
			application.Log().Logf("[ERROR] upload subnet err: %s", err)
		}
	}

	return nil
}

func fileInspectFn() error {
	_, err := configuration.Load(opts.File)
	if err != nil {
		fmt.Printf("Configuration file %s has error:\n%v\n", opts.File, err)
	} else {
		fmt.Printf("Configuration file %s is valid\n", opts.File)
	}

	return nil
}

func filePreviewFn() error {
	conf, err := configuration.Load(opts.File)
	if err != nil {
		return err
	}

	rv, err := initResolver()
	if err != nil {
		return err
	}

	for index, value := range conf {
		fmt.Printf("Connection %d: %s\n", index+1, value.Title)
		subnets := make(map[string]utils.ResolvedSubnet)

		for _, value := range value.Values {
			ss, err := rv.Resolve(value)
			if err != nil {
				fmt.Printf("\tHas error: %s\n", err)
				continue
			}

			for _, subnet := range ss {
				subnets[subnet.Hash()] = subnet
			}
		}

		fmt.Println("Subnets:")

		count := 0
		for _, subnet := range subnets {
			count++
			fmt.Printf("\t%d) %s\n", count, subnet)
		}
	}

	return nil
}

func rciClearFn() error {
	rci := opts.RCI.MustInitRCI()

	resp, err := rci.ListIPRoutes()
	if err != nil {
		return err
	}

	for _, route := range resp {
		if route.IsAppRoute() {
			err := rci.DeleteIPRoute(route)
			if err != nil {
				fmt.Printf("Failed to delete route %s: %s\n", route.Hash(), err)
			}
		}
	}

	return nil
}

func rciListFn() error {
	rci := opts.RCI.MustInitRCI()

	resp, err := rci.GetInterfaces()
	if err != nil {
		return err
	}

	fmt.Printf("rci %s interfaces list:\n", opts.RCI.Host)
	for name, value := range resp {
		fmt.Printf("\t%s: %s\n", name, value)
	}

	return nil
}

func rciInspectFn() error {
	rci := opts.RCI.MustInitRCI()

	resp, err := rci.ListIPRoutes()
	if err != nil {
		return err
	}

	fmt.Printf("rci %s routes list:\n", opts.RCI.Host)

	for _, route := range resp {
		fmt.Printf(
			"\t- Host: %s\n\t  Interface: %s\n\t  Mask: %s\n\t  Network: %s\n",
			route.Host, route.Interface, route.Mask, route.Network)
	}

	return nil
}

func selectAction() Action {
	counter := 0
	actions := map[Action]bool{}

	if opts.Action.Upload {
		actions[ActionUpload] = true
		counter++
	}

	if opts.Action.File.Inspect {
		actions[ActionFileInspect] = true
		counter++
	}

	if opts.Action.File.Preview {
		actions[ActionFilePreview] = true
		counter++
	}

	if opts.Action.RCI.Clear {
		actions[ActionRCIClear] = true
		counter++
	}

	if opts.Action.RCI.Inspect {
		actions[ActionRCIInspect] = true
		counter++
	}

	if opts.Action.RCI.List {
		actions[ActionRCIList] = true
		counter++
	}

	if counter > 2 {
		return "select only one action"
	}

	for action, selected := range actions {
		if selected {
			return action
		}
	}

	return ActionUpload
}

func initResolver() (*resolver.Resolver, error) {
	resolvers := []resolver.ResolverItem{
		ip.New(),
	}

	if opts.DNS.Enabled {
		resolvers = append(resolvers, dns.New(opts.DNS.Host...))
		application.Log().Logf("[INFO] DNS resolver enabled by %s", strings.Join(opts.DNS.Host, ", "))
	}

	if opts.ASN.Enabled {
		resolvers = append(resolvers, asn.New(opts.ASN.Key))
		application.Log().Logf("[INFO] ASN resolver enabled")
	}

	return resolver.NewResolver(resolvers...), nil
}
