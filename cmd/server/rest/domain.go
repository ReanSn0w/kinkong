package rest

import (
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/ReanSn0w/gokit/pkg/web/mv/json"
	"github.com/ReanSn0w/kincong/cmd/server/cr"
	"github.com/ReanSn0w/kincong/internal/utils"
)

type DomainRequest struct {
	Value string `json:"value" example:"google.com"`
}

type DomainResponse web.Response[DomainResponseData]

type DomainResponseData cr.DomainInfo

// @Summary		Получение информации о домене
// @Description	Метод возвращает IP, Подсеть и список ASN за которыми закреплен домен
// @Accept		json
// @Produce		json
// @Param       body body DomainRequest true "Параметры запроса"
// @Success		200	    {object} DomainResponse
// @Failure		default	{object} ResponseError
// @Router		/api/domain [post]
func (s *Server) domainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := json.Get[DomainRequest](ctx)

	domainInfo, err := s.resolver.DomainInfo(utils.Value(req.Value))
	if err != nil {
		web.NewResponse(err).Write(http.StatusInternalServerError, w)
		return
	}

	web.NewResponse(domainInfo).Write(http.StatusOK, w)
}
