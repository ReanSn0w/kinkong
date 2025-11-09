package rest

import (
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/ReanSn0w/gokit/pkg/web/mv/json"
	"github.com/ReanSn0w/kincong/internal/utils"
)

type ResolveRequest struct {
	// Список IP адресов для разрешения
	IPs []string `json:"ip" example:"8.8.8.8"`
	// Список доменных имен для разрешения
	Domains []string `json:"domain" example:"google.com"`
	// Список ASN для разрешения
	ASNs []string `json:"asn" example:"AS12345"`
}

func (r *ResolveRequest) Values() []utils.Value {
	result := make([]utils.Value, 0, len(r.IPs)+len(r.Domains)+len(r.ASNs))
	for _, ip := range r.IPs {
		result = append(result, utils.Value(ip))
	}
	for _, domain := range r.Domains {
		result = append(result, utils.Value(domain))
	}
	for _, asn := range r.ASNs {
		result = append(result, utils.Value(asn))
	}
	return result
}

type ResolveResponse web.Response[ResolveResponseData]

type ResolveResponseData []string

// @Summary		Получение списка подсетей
// @Description	Метод возвращает пользователю список доступных в системе документов
// @Accept		json
// @Produce		json
// @Param       body body ResolveRequest true "Параметры запроса"
// @Success		200	    {object} ResolveResponse
// @Failure		default	{object} ResponseError
// @Router		/api/resolve [post]
func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := json.Get[ResolveRequest](ctx)

	networks, err := s.resolver.Resolve(req.Values())
	if err != nil {
		web.NewResponse(err).Write(http.StatusInternalServerError, w)
		return
	}

	web.NewResponse(networks).Write(http.StatusOK, w)
}
